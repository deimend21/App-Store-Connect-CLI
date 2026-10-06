package asc

import "strconv"

// CreativePlacementCreateResult records a successful public placement creation.
type CreativePlacementCreateResult struct {
	ID             string `json:"id"`
	Created        bool   `json:"created"`
	LocalizationID string `json:"localizationId"`
	ImageID        string `json:"imageId"`
	PlacementType  string `json:"placementType"`
	PlacementGroup string `json:"placementGroup"`
	State          string `json:"state,omitempty"`
}

// CreativePlacementDeleteResult records a successful public placement removal.
type CreativePlacementDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func creativePlacementCreateResultRows(r *CreativePlacementCreateResult) ([]string, [][]string) {
	return []string{"ID", "Created", "Localization ID", "Image ID", "Placement Type", "Group", "State"}, [][]string{{r.ID, strconv.FormatBool(r.Created), r.LocalizationID, r.ImageID, r.PlacementType, r.PlacementGroup, r.State}}
}

func creativePlacementDeleteResultRows(r *CreativePlacementDeleteResult) ([]string, [][]string) {
	return []string{"ID", "Deleted"}, [][]string{{r.ID, strconv.FormatBool(r.Deleted)}}
}
