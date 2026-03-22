package territory

import (
	"a-star-is-born/internal/model"
	"a-star-is-born/internal/storage"
	"fmt"
)

func GenerateTerritoryJSON(inputPath, outputPath string) error {
	var snapshot model.Snapshot

	if err := storage.LoadJSON(inputPath, &snapshot); err != nil {
		return fmt.Errorf("load snapshot failed: %w", err)
	}

	tree, err := BuildContributorTerritory(snapshot)
	if err != nil {
		return fmt.Errorf("build territory failed: %w", err)
	}

	if err := storage.SaveJSON(outputPath, tree); err != nil {
		return fmt.Errorf("save territory failed: %w", err)
	}

	return nil
}
