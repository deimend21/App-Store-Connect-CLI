package capabilities

import (
	"slices"
	"strings"
	"testing"
)

func TestAssetLibraryCapabilitiesIncludeImageUpload(t *testing.T) {
	for _, c := range capabilityRows() {
		if c.Capability == "Asset Library media and specifications" {
			if c.Status != statusCLISupported || !strings.Contains(strings.Join(c.Notes, " "), "image upload") {
				t.Fatalf("incorrect scope: %+v", c)
			}
			return
		}
	}
	t.Fatal("missing verified Asset Library capability")
}

func TestAssetLibraryCapabilitiesExposeStandaloneReview(t *testing.T) {
	for _, c := range capabilityRows() {
		if c.Capability != "Asset Library media and specifications" {
			continue
		}
		for _, command := range []string{"asc review submissions-create", "asc review items add", "asc review submissions-submit --confirm"} {
			if !slices.Contains(c.Commands, command) {
				t.Errorf("missing review command %q: %+v", command, c)
			}
		}
		for _, resource := range []string{"appAssetLibraryImages", "appAssetLibraryVideos", "reviewSubmissions", "reviewSubmissionItems"} {
			if !slices.Contains(c.APIResources, resource) {
				t.Errorf("missing review resource %q: %+v", resource, c)
			}
		}
		notes := strings.Join(c.Notes, " ")
		if strings.Contains(notes, "Review submission is not exposed") || !strings.Contains(notes, "public review-item POST") || !strings.Contains(notes, "App Review acceptance remain unverified") {
			t.Errorf("misleading verification scope: %s", notes)
		}
		return
	}
	t.Fatal("missing Asset Library capability")
}
