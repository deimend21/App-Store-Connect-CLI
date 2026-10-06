package capabilities

import (
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
