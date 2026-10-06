package capabilities

import (
	"strings"
	"testing"
)

func TestAssetLibraryCapabilitiesAreReadOnly(t *testing.T) {
	for _, c := range capabilityRows() {
		if c.Capability == "Asset Library media and specifications" {
			if c.Status != statusCLISupported || !strings.Contains(strings.Join(c.Notes, " "), "Read-only") {
				t.Fatalf("incorrect scope: %+v", c)
			}
			return
		}
	}
	t.Fatal("missing verified Asset Library capability")
}
