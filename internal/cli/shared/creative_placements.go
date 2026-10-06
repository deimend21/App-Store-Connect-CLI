package shared

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
)

var placementLocalizationID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// CreativePlacementsCommand returns a live-verified placement read group.
func CreativePlacementsCommand(resource, prefix, label string) *ffcli.Command {
	return &ffcli.Command{Name: "placements", ShortHelp: "Inspect localized product page and search asset placements.", FlagSet: flag.NewFlagSet("placements", flag.ExitOnError), UsageFunc: DefaultUsageFunc, Subcommands: []*ffcli.Command{creativePlacementsListCommand(resource, prefix, label)}, Exec: func(context.Context, []string) error { return flag.ErrHelp }}
}

func creativePlacementsListCommand(resource, prefix, label string) *ffcli.Command {
	fs := flag.NewFlagSet("placements list", flag.ExitOnError)
	id := BindResourceIDFlag(fs, "localization-id", resource, label+" ID or API self-link")
	kind := fs.String("placement-type", "", "Filter by placement types, comma-separated: PRODUCT_PAGE_HEADER_ASSET, APP_STORE_SEARCH_RESULTS_ASSET, APP_SCREENSHOT, APP_PREVIEW, IMESSAGE_APP_SCREENSHOT")
	include := fs.String("include", "", "Include related media, comma-separated: image, video")
	sort := fs.String("sort", "", "Sort by placementGroupPosition")
	limit := fs.Int("limit", 0, "Maximum results per page (0 uses server default)")
	next := fs.String("next", "", "Fetch a links.next URL instead of selecting a localization")
	paginate := fs.Bool("paginate", false, "Fetch all pages and aggregate the collection")
	output := BindOutputFlags(fs)
	return &ffcli.Command{
		Name: "list", ShortHelp: "List asset placements for " + label + ".",
		ShortUsage: "asc " + prefix + " placements list [flags]", FlagSet: fs, UsageFunc: DefaultUsageFunc,
		LongHelp: "List asset placements for " + label + ". These public GET endpoints and header/search filters are live-verified but absent from Apple's published OpenAPI 4.5. Empty data means no matching asset is assigned. JSON preserves Apple's full envelope and included media.\n\nExamples:\n  asc " + prefix + " placements list --localization-id ID --placement-type PRODUCT_PAGE_HEADER_ASSET\n  asc " + prefix + " placements list --localization-id ID --placement-type APP_STORE_SEARCH_RESULTS_ASSET --include image,video\n  asc " + prefix + " placements list --localization-id ID --paginate --output table",
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return UsageErrorf(prefix+" placements list: unexpected argument %q", args[0])
			}
			if err := ValidateBoundOutputFlags(fs); err != nil {
				return UsageErrorf(prefix+" placements list: %v", err)
			}
			if *limit < 0 {
				return UsageError(prefix + " placements list: --limit must be zero or a positive integer")
			}
			if err := ValidateNextURL(*next); err != nil {
				return UsageErrorf(prefix+" placements list: %v", err)
			}
			if err := RejectNextFlagConflicts(fs, *next, prefix+" placements list", "localization-id", "placement-type", "include", "sort", "limit"); err != nil {
				return err
			}
			kinds, err := validatedPlacementCSV(*kind, []string{"PRODUCT_PAGE_HEADER_ASSET", "APP_STORE_SEARCH_RESULTS_ASSET", "APP_SCREENSHOT", "APP_PREVIEW", "IMESSAGE_APP_SCREENSHOT"})
			if err != nil {
				return UsageErrorf(prefix+" placements list: --placement-type %v", err)
			}
			includes, err := validatedPlacementCSV(*include, []string{"image", "video"})
			if err != nil {
				return UsageErrorf(prefix+" placements list: --include %v", err)
			}
			if *sort != "" && *sort != "placementGroupPosition" {
				return UsageError(prefix + " placements list: --sort must be placementGroupPosition")
			}
			target := *next
			if target == "" {
				localizationID := strings.TrimSpace(*id)
				if localizationID == "" {
					return UsageError(prefix + " placements list: --localization-id is required")
				}
				if !placementLocalizationID.MatchString(localizationID) {
					return UsageError(prefix + " placements list: --localization-id must be a resource ID")
				}
				target = "/v1/" + resource + "/" + url.PathEscape(localizationID) + "/placements"
				q := url.Values{}
				if kinds != "" {
					q.Set("filter[placementType]", kinds)
				}
				if includes != "" {
					q.Set("include", includes)
				}
				if *sort != "" {
					q.Set("sort", *sort)
				}
				if *limit > 0 {
					q.Set("limit", strconv.Itoa(*limit))
				}
				if len(q) > 0 {
					target += "?" + q.Encode()
				}
			}
			client, err := GetASCClient()
			if err != nil {
				return fmt.Errorf(prefix+" placements list: %w", err)
			}
			var response []byte
			if *paginate {
				response, err = client.RawPaginatedGET(ctx, target, ContextWithTimeout)
			} else {
				requestCtx, cancel := ContextWithTimeout(ctx)
				defer cancel()
				response, err = client.RawRequest(requestCtx, "GET", target, nil)
			}
			if err != nil {
				return fmt.Errorf(prefix+" placements list: %w", err)
			}
			headers, rows, err := asc.CreativePlacementRows(response)
			if err != nil {
				return err
			}
			return PrintOutputRows(json.RawMessage(response), *output.Output, *output.Pretty, headers, rows)
		},
	}
}

func validatedPlacementCSV(raw string, allowed []string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", nil
	}
	values := strings.Split(raw, ",")
	for i, value := range values {
		value = strings.TrimSpace(value)
		valid := false
		for _, candidate := range allowed {
			if value == candidate {
				valid = true
				break
			}
		}
		if !valid {
			return "", fmt.Errorf("contains unsupported value %q", value)
		}
		values[i] = value
	}
	return strings.Join(values, ","), nil
}
