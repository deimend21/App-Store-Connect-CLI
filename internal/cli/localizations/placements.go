package localizations

import (
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// LocalizationsPlacementsCommand exposes localized header and search placements.
func LocalizationsPlacementsCommand() *ffcli.Command {
	cmd := shared.CreativePlacementsCommand("appStoreVersionLocalizations", "localizations", "an App Store version localization")
	cmd.Subcommands = append(cmd.Subcommands, localizationsPlacementCreateCommand(), localizationsPlacementDeleteCommand())
	return cmd
}
