// Package assetlibrary exposes live-verified public Asset Library reads.
package assetlibrary

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// Command returns the Asset Library command group.
func Command() *ffcli.Command {
	return group("asset-library", "Inspect App Store Asset Library media and specifications.",
		readCommand("view", "Read an app's Asset Library.", "app", "/v1/apps/%s/assetLibrary", false),
		group("images", "Inspect Asset Library images.",
			readCommand("list", "List images in an Asset Library.", "library-id", "/v1/appAssetLibraries/%s/images", true),
			readCommand("view", "Read an Asset Library image.", "id", "/v1/appAssetLibraryImages/%s", false),
			readCommand("placements", "List where an Asset Library image is used.", "id", "/v1/appAssetLibraryImages/%s/placements", true)),
		group("videos", "Inspect Asset Library videos.",
			readCommand("list", "List videos in an Asset Library.", "library-id", "/v1/appAssetLibraries/%s/videos", true)),
		readCommand("specs", "Read Apple's current asset dimensions and placement policies.", "", "/v1/appAssetLibraryRefData", false))
}

func group(name, help string, children ...*ffcli.Command) *ffcli.Command {
	return &ffcli.Command{
		Name: name, ShortHelp: help,
		LongHelp: help + ` These public GET endpoints are live-verified but absent from Apple's published OpenAPI 4.5. Upload, assignment, and submission are not exposed.

Examples:
  asc asset-library view --app APP_ID
  asc asset-library images list --library-id LIBRARY_ID --paginate
  asc asset-library images view --id IMAGE_ID
  asc asset-library images placements --id IMAGE_ID
  asc asset-library videos list --library-id LIBRARY_ID
  asc asset-library specs`,
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError), UsageFunc: shared.DefaultUsageFunc,
		Subcommands: children, Exec: func(context.Context, []string) error { return flag.ErrHelp },
	}
}

var resourceIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func readCommand(name, help, selector, path string, collection bool) *ffcli.Command {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	id := new(string)
	if selector != "" {
		resourceType := "appAssetLibraries"
		switch selector {
		case "app":
			resourceType = "apps"
		case "id":
			resourceType = "appAssetLibraryImages"
		}
		id = shared.BindResourceIDFlag(fs, selector, resourceType, "Resource ID or API self-link (use asset-library view --app to find the library ID)")
	}
	var limit int
	var next string
	var paginate bool
	if collection {
		fs.IntVar(&limit, "limit", 0, "Maximum results per page (0 uses server default)")
		fs.StringVar(&next, "next", "", "Fetch a links.next URL instead of selecting a resource")
		fs.BoolVar(&paginate, "paginate", false, "Fetch all pages and aggregate the collection")
	}
	output := shared.BindOutputFlags(fs)
	return &ffcli.Command{
		Name: name, ShortHelp: help, FlagSet: fs, UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return shared.UsageErrorf("asset-library %s: unexpected argument %q", name, args[0])
			}
			if err := shared.ValidateBoundOutputFlags(fs); err != nil {
				return shared.UsageErrorf("asset-library %s: %v", name, err)
			}
			if limit < 0 {
				return shared.UsageError("asset-library: --limit must be zero or a positive integer")
			}
			if next != "" {
				if err := shared.ValidateNextURL(next); err != nil {
					return shared.UsageErrorf("asset-library: %v", err)
				}
				if err := shared.RejectNextFlagConflicts(fs, next, "asset-library", selector, "limit"); err != nil {
					return err
				}
			}
			resourceID := strings.TrimSpace(*id)
			if selector == "app" {
				resourceID = shared.ResolveAppID(resourceID)
			}
			target := path
			if selector != "" && next == "" {
				if resourceID == "" {
					return shared.UsageErrorf("asset-library: --%s is required", selector)
				}
				if !resourceIDPattern.MatchString(resourceID) {
					return shared.UsageErrorf("asset-library: --%s must be a resource ID", selector)
				}
				target = fmt.Sprintf(path, url.PathEscape(resourceID))
			}
			if next != "" {
				target = next
			} else if limit > 0 {
				target += "?limit=" + strconv.Itoa(limit)
			}
			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("asset-library: %w", err)
			}
			var response []byte
			if paginate {
				response, err = client.RawPaginatedGET(ctx, target, shared.ContextWithTimeout)
			} else {
				requestCtx, cancel := shared.ContextWithTimeout(ctx)
				defer cancel()
				response, err = client.RawRequest(requestCtx, "GET", target, nil)
			}
			if err != nil {
				return fmt.Errorf("asset-library: %w", err)
			}
			if err := asc.ValidateAssetLibraryResponse(response); err != nil {
				return err
			}
			render := func(markdown bool) error {
				headers, rows, err := asc.AssetLibraryRows(response)
				if err != nil {
					return err
				}
				if markdown {
					return asc.WriteMarkdown(headers, rows)
				}
				return asc.WriteTable(headers, rows)
			}
			return shared.PrintOutputWithRenderers(json.RawMessage(response), *output.Output, *output.Pretty, func() error { return render(false) }, func() error { return render(true) })
		},
	}
}
