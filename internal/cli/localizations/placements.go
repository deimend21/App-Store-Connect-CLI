package localizations

import (
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// LocalizationsPlacementsCommand exposes localized header and search placement reads.
func LocalizationsPlacementsCommand() *ffcli.Command {
	return shared.CreativePlacementsCommand("appStoreVersionLocalizations", "localizations", "an App Store version localization")
}
