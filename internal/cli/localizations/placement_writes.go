package localizations

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

var creativePlacementID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func localizationsPlacementCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("placements create", flag.ExitOnError)
	loc := shared.BindResourceIDFlag(fs, "localization-id", "appStoreVersionLocalizations", "App Store version localization ID or API self-link")
	image := shared.BindResourceIDFlag(fs, "image-id", "appAssetLibraryImages", "Uploaded asset library image ID or API self-link")
	kind := fs.String("placement-type", "", "PRODUCT_PAGE_HEADER_ASSET or APP_STORE_SEARCH_RESULTS_ASSET")
	output := shared.BindOutputFlags(fs)
	return &ffcli.Command{
		Name: "create", ShortHelp: "Assign an asset library image as a product page header or search asset.", ShortUsage: "asc localizations placements create --localization-id ID --image-id ID --placement-type TYPE", FlagSet: fs, UsageFunc: shared.DefaultUsageFunc,
		LongHelp: `Assign an uploaded asset library image to an App Store version localization.
Only PRODUCT_PAGE_HEADER_ASSET and APP_STORE_SEARCH_RESULTS_ASSET in
DEFAULT_PROFILE have been verified with the public API. This assigns the image;
it does not submit it for review or publish it. Apple validates image suitability
and slot availability. No existing placement is removed automatically.

Examples:
  asc localizations placements create --localization-id ID --image-id IMAGE_ID --placement-type PRODUCT_PAGE_HEADER_ASSET
  asc localizations placements create --localization-id ID --image-id IMAGE_ID --placement-type APP_STORE_SEARCH_RESULTS_ASSET`,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return shared.UsageErrorf("placements create: unexpected argument %q", args[0])
			}
			if err := shared.ValidateBoundOutputFlags(fs); err != nil {
				return shared.UsageErrorf("placements create: %v", err)
			}
			localizationID, imageID := strings.TrimSpace(*loc), strings.TrimSpace(*image)
			for _, input := range []struct{ name, value string }{{"localization-id", localizationID}, {"image-id", imageID}} {
				if !creativePlacementID.MatchString(input.value) {
					return shared.UsageErrorf("placements create: --%s requires a resource ID", input.name)
				}
			}
			if *kind != "PRODUCT_PAGE_HEADER_ASSET" && *kind != "APP_STORE_SEARCH_RESULTS_ASSET" {
				return shared.UsageError("placements create: --placement-type must be PRODUCT_PAGE_HEADER_ASSET or APP_STORE_SEARCH_RESULTS_ASSET")
			}
			relationship := func(resource, id string) map[string]any {
				return map[string]any{"data": map[string]string{"type": resource, "id": id}}
			}
			payload, err := json.Marshal(map[string]any{"data": map[string]any{"type": "appAssetLibraryPlacements", "attributes": map[string]string{"placementType": *kind, "placementGroup": "DEFAULT_PROFILE"}, "relationships": map[string]any{"image": relationship("appAssetLibraryImages", imageID), "appStoreVersionLocalization": relationship("appStoreVersionLocalizations", localizationID)}}})
			if err != nil {
				return err
			}
			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("placements create: %w", err)
			}
			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()
			response, err := client.RawRequest(requestCtx, "POST", "/v1/appAssetLibraryPlacements", payload)
			if err != nil {
				return fmt.Errorf("placements create: %w", err)
			}
			var envelope struct {
				Data struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Attributes struct {
						State string `json:"state"`
					} `json:"attributes"`
				} `json:"data"`
			}
			if err := json.Unmarshal(response, &envelope); err != nil {
				return fmt.Errorf("placements create: invalid response: %w", err)
			}
			if envelope.Data.ID == "" || envelope.Data.Type != "appAssetLibraryPlacements" {
				return fmt.Errorf("placements create: invalid placement response")
			}
			return shared.PrintOutput(&asc.CreativePlacementCreateResult{ID: envelope.Data.ID, Created: true, LocalizationID: localizationID, ImageID: imageID, PlacementType: *kind, PlacementGroup: "DEFAULT_PROFILE", State: envelope.Data.Attributes.State}, *output.Output, *output.Pretty)
		},
	}
}

func localizationsPlacementDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("placements delete", flag.ExitOnError)
	id := shared.BindResourceIDFlag(fs, "id", "appAssetLibraryPlacements", "Asset library placement ID or API self-link")
	confirm := fs.Bool("confirm", false, "Confirm removal of this placement")
	output := shared.BindOutputFlags(fs)
	return &ffcli.Command{
		Name: "delete", ShortHelp: "Remove an asset placement without deleting its library media.", ShortUsage: "asc localizations placements delete --id ID --confirm", FlagSet: fs, UsageFunc: shared.DefaultUsageFunc,
		LongHelp: `Remove an asset library placement by its placement ID.
Only the placement is removed; its library media remains available for reuse.
This command does not validate the placement's localization or type.
Find placement IDs with
asc localizations placements list. Removal requires --confirm.

Example:
  asc localizations placements delete --id PLACEMENT_ID --confirm`,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return shared.UsageErrorf("placements delete: unexpected argument %q", args[0])
			}
			if err := shared.ValidateBoundOutputFlags(fs); err != nil {
				return shared.UsageErrorf("placements delete: %v", err)
			}
			placementID := strings.TrimSpace(*id)
			if !creativePlacementID.MatchString(placementID) {
				return shared.UsageError("placements delete: --id requires a resource ID")
			}
			if !*confirm {
				return shared.UsageError("placements delete: --confirm is required")
			}
			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("placements delete: %w", err)
			}
			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()
			if _, err := client.RawRequest(requestCtx, "DELETE", "/v1/appAssetLibraryPlacements/"+url.PathEscape(placementID), nil); err != nil {
				return fmt.Errorf("placements delete: %w", err)
			}
			return shared.PrintOutput(&asc.CreativePlacementDeleteResult{ID: placementID, Deleted: true}, *output.Output, *output.Pretty)
		},
	}
}
