package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"fpm/internal/archive"
	"fpm/internal/metadata"

	"github.com/spf13/cobra"
)

// validateFrappeAppStructure checks if the source directory has a valid Frappe app structure.
func validateFrappeAppStructure(sourceDir string, appName string) error {
	// Check 1: Existence of directory sourceDir + "/" + appName
	innerAppPath := filepath.Join(sourceDir, appName)
	info, err := os.Stat(innerAppPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Frappe app validation failed: app directory '%s' not found", innerAppPath)
	}
	if err != nil {
		return fmt.Errorf("Frappe app validation failed: error checking app directory '%s': %w", innerAppPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Frappe app validation failed: '%s' is not a directory", innerAppPath)
	}

	// Check 2: Existence of file sourceDir + "/" + appName + "/__init__.py"
	initPyPath := filepath.Join(innerAppPath, "__init__.py")
	info, err = os.Stat(initPyPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Frappe app validation failed: file '%s' not found", initPyPath)
	}
	if err != nil {
		return fmt.Errorf("Frappe app validation failed: error checking file '%s': %w", initPyPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("Frappe app validation failed: '%s' is a directory, not a file", initPyPath)
	}

	// Check 3: Existence of file sourceDir + "/" + appName + "/hooks.py"
	hooksPyPath := filepath.Join(innerAppPath, "hooks.py")
	info, err = os.Stat(hooksPyPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Frappe app validation failed: file '%s' not found", hooksPyPath)
	}
	if err != nil {
		return fmt.Errorf("Frappe app validation failed: error checking file '%s': %w", hooksPyPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("Frappe app validation failed: '%s' is a directory, not a file", hooksPyPath)
	}

	// Check 4: Existence of file sourceDir + "/" + appName + "/modules.txt"
	modulesTxtPath := filepath.Join(innerAppPath, "modules.txt")
	info, err = os.Stat(modulesTxtPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Frappe app validation failed: file '%s' not found", modulesTxtPath)
	}
	if err != nil {
		return fmt.Errorf("Frappe app validation failed: error checking file '%s': %w", modulesTxtPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("Frappe app validation failed: '%s' is a directory, not a file", modulesTxtPath)
	}

	return nil // All checks passed
}

var (
	packageSourcePath string
	packageOutputPath string
	packageVersion    string
	packageOverwrite  bool
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package a Frappe application into an .fpm file",
	Long: `Packages a Frappe application from a local development directory into an .fpm file.
It reads app metadata, collects source files, and bundles them into a versioned archive.`,
	RunE: func(cmd *cobra.Command, args []string) error { // Using RunE for error handling
		if packageVersion == "" {
			return fmt.Errorf("--version flag is required")
		}

		absSourcePath, err := filepath.Abs(packageSourcePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute source path: %w", err)
		}

		if _, err := os.Stat(absSourcePath); os.IsNotExist(err) {
			return fmt.Errorf("source path '%s' does not exist", absSourcePath)
		}

		// Load existing metadata or generate a new one
		meta, err := metadata.LoadAppMetadata(absSourcePath)
		if err != nil {
			// If LoadAppMetadata returns an error for reasons other than file not found,
			// or if we decide it should error if file not found, handle here.
			// For now, LoadAppMetadata returns empty struct if not found.
			// Let's assume if meta.PackageName is empty after Load, we should generate.
		}

		// If package name is still empty (either file didn't exist or was empty), generate.
		if meta.PackageName == "" {
		    inferredMeta, genErr := metadata.GenerateAppMetadata(absSourcePath, packageVersion)
		    if genErr != nil {
		        return fmt.Errorf("failed to generate default app metadata: %w", genErr)
		    }
		    meta = inferredMeta // Use the generated one
		} else {
            // If loaded, still ensure the CLI version overrides
		    meta.PackageVersion = packageVersion
        }
        // If GenerateAppMetadata was called, it already set the version.
        // If LoadAppMetadata was called and it was successful, PackageVersion in meta
        // will be updated by the GenerateAppMetadata or the line above.

		// Validate Frappe app structure
		if meta.PackageName == "" {
			// This should ideally be caught by GenerateAppMetadata if it's responsible for determining name
			return fmt.Errorf("app package name could not be determined, cannot validate structure")
		}
		if err := validateFrappeAppStructure(absSourcePath, meta.PackageName); err != nil {
			return err // The error from validateFrappeAppStructure is already descriptive
		}

		outputFileName := fmt.Sprintf("%s-%s.fpm", meta.PackageName, packageVersion)
		absOutputPath, err := filepath.Abs(packageOutputPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute output path: %w", err)
		}

		finalFpmFilePath := filepath.Join(absOutputPath, outputFileName)

		if _, err := os.Stat(finalFpmFilePath); err == nil && !packageOverwrite {
			return fmt.Errorf("output file '%s' already exists. Use --overwrite to replace it", finalFpmFilePath)
		}

		fmt.Printf("Packaging '%s' version '%s' from '%s'...\n", meta.PackageName, packageVersion, absSourcePath)

		err = archive.CreateFPMArchive(absSourcePath, absOutputPath, meta, packageVersion)
		if err != nil {
			return fmt.Errorf("failed to create package: %w", err)
		}

		fmt.Printf("Successfully packaged: %s\n", finalFpmFilePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.Flags().StringVarP(&packageSourcePath, "source", "s", ".", "Path to the Frappe app source directory")
	packageCmd.Flags().StringVarP(&packageOutputPath, "output-path", "o", ".", "Directory to save the .fpm file")
	packageCmd.Flags().StringVarP(&packageVersion, "version", "v", "", "Package version (e.g., 1.0.0) (required)")
	packageCmd.Flags().BoolVar(&packageOverwrite, "overwrite", false, "Overwrite if .fpm file already exists")

	// Mark version as required if using cobra's built-in way, though manual check is also fine.
	// packageCmd.MarkFlagRequired("version") // This causes help text to show if not provided.
}
