package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFrappeAppStructure(t *testing.T) {
	appName := "testapp"

	// Helper to create dummy files
	createFile := func(path string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		return f.Close()
	}

	t.Run("valid app structure", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-valid-app-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		innerAppDir := filepath.Join(tmpDir, appName)
		if err := os.Mkdir(innerAppDir, 0755); err != nil {
			t.Fatalf("Failed to create inner app dir: %v", err)
		}

		if err := createFile(filepath.Join(innerAppDir, "__init__.py")); err != nil {
			t.Fatalf("Failed to create __init__.py: %v", err)
		}
		if err := createFile(filepath.Join(innerAppDir, "hooks.py")); err != nil {
			t.Fatalf("Failed to create hooks.py: %v", err)
		}
		if err := createFile(filepath.Join(innerAppDir, "modules.txt")); err != nil {
			t.Fatalf("Failed to create modules.txt: %v", err)
		}

		err = validateFrappeAppStructure(tmpDir, appName)
		if err != nil {
			t.Errorf("Expected no error for valid structure, got %v", err)
		}
	})

	t.Run("missing inner app directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-missing-inner-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Do not create innerAppDir

		err = validateFrappeAppStructure(tmpDir, appName)
		if err == nil {
			t.Errorf("Expected error for missing inner app directory, got nil")
		} else {
			expectedErrorPart := fmt.Sprintf("app directory '%s' not found", filepath.Join(tmpDir, appName))
			if !strings.Contains(err.Error(), expectedErrorPart) {
				t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
			}
		}
	})

	t.Run("inner app path is a file not a directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-inner-is-file-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create a file where the directory should be
		if err := createFile(filepath.Join(tmpDir, appName)); err != nil {
			t.Fatalf("Failed to create dummy file for inner app path: %v", err)
		}

		err = validateFrappeAppStructure(tmpDir, appName)
		if err == nil {
			t.Errorf("Expected error when inner app path is a file, got nil")
		} else {
			expectedErrorPart := fmt.Sprintf("'%s' is not a directory", filepath.Join(tmpDir, appName))
			if !strings.Contains(err.Error(), expectedErrorPart) {
				t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
			}
		}
	})

	t.Run("missing __init__.py", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-missing-init-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		innerAppDir := filepath.Join(tmpDir, appName)
		if err := os.Mkdir(innerAppDir, 0755); err != nil {
			t.Fatalf("Failed to create inner app dir: %v", err)
		}

		// Missing __init__.py
		if err := createFile(filepath.Join(innerAppDir, "hooks.py")); err != nil {
			t.Fatalf("Failed to create hooks.py: %v", err)
		}
		if err := createFile(filepath.Join(innerAppDir, "modules.txt")); err != nil {
			t.Fatalf("Failed to create modules.txt: %v", err)
		}

		err = validateFrappeAppStructure(tmpDir, appName)
		if err == nil {
			t.Errorf("Expected error for missing __init__.py, got nil")
		} else {
			expectedErrorPart := fmt.Sprintf("file '%s' not found", filepath.Join(innerAppDir, "__init__.py"))
			if !strings.Contains(err.Error(), expectedErrorPart) {
				t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
			}
		}
	})

	t.Run("missing hooks.py", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-missing-hooks-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		innerAppDir := filepath.Join(tmpDir, appName)
		if err := os.Mkdir(innerAppDir, 0755); err != nil {
			t.Fatalf("Failed to create inner app dir: %v", err)
		}

		if err := createFile(filepath.Join(innerAppDir, "__init__.py")); err != nil {
			t.Fatalf("Failed to create __init__.py: %v", err)
		}
		// Missing hooks.py
		if err := createFile(filepath.Join(innerAppDir, "modules.txt")); err != nil {
			t.Fatalf("Failed to create modules.txt: %v", err)
		}

		err = validateFrappeAppStructure(tmpDir, appName)
		if err == nil {
			t.Errorf("Expected error for missing hooks.py, got nil")
		} else {
			expectedErrorPart := fmt.Sprintf("file '%s' not found", filepath.Join(innerAppDir, "hooks.py"))
			if !strings.Contains(err.Error(), expectedErrorPart) {
				t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
			}
		}
	})

	t.Run("missing modules.txt", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-missing-modules-")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		innerAppDir := filepath.Join(tmpDir, appName)
		if err := os.Mkdir(innerAppDir, 0755); err != nil {
			t.Fatalf("Failed to create inner app dir: %v", err)
		}

		if err := createFile(filepath.Join(innerAppDir, "__init__.py")); err != nil {
			t.Fatalf("Failed to create __init__.py: %v", err)
		}
		if err := createFile(filepath.Join(innerAppDir, "hooks.py")); err != nil {
			t.Fatalf("Failed to create hooks.py: %v", err)
		}
		// Missing modules.txt

		err = validateFrappeAppStructure(tmpDir, appName)
		if err == nil {
			t.Errorf("Expected error for missing modules.txt, got nil")
		} else {
			expectedErrorPart := fmt.Sprintf("file '%s' not found", filepath.Join(innerAppDir, "modules.txt"))
			if !strings.Contains(err.Error(), expectedErrorPart) {
				t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
			}
		}
	})

	// Test cases for when a required component is a directory instead of a file
	testCasesIsDirectory := []struct {
		name         string
		fileToMakeDir string // e.g., "__init__.py"
	}{
		{"__init__.py is a directory", "__init__.py"},
		{"hooks.py is a directory", "hooks.py"},
		{"modules.txt is a directory", "modules.txt"},
	}

	for _, tc := range testCasesIsDirectory {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "test-isdir-")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			innerAppDir := filepath.Join(tmpDir, appName)
			if err := os.MkdirAll(innerAppDir, 0755); err != nil { // MkdirAll in case appName has slashes (not typical here)
				t.Fatalf("Failed to create inner app dir: %v", err)
			}

			// Create other required files
			if tc.fileToMakeDir != "__init__.py" {
				if err := createFile(filepath.Join(innerAppDir, "__init__.py")); err != nil {
					t.Fatalf("Failed to create __init__.py: %v", err)
				}
			}
			if tc.fileToMakeDir != "hooks.py" {
				if err := createFile(filepath.Join(innerAppDir, "hooks.py")); err != nil {
					t.Fatalf("Failed to create hooks.py: %v", err)
				}
			}
			if tc.fileToMakeDir != "modules.txt" {
				if err := createFile(filepath.Join(innerAppDir, "modules.txt")); err != nil {
					t.Fatalf("Failed to create modules.txt: %v", err)
				}
			}

			// Create the component that should be a file as a directory
			pathAsDir := filepath.Join(innerAppDir, tc.fileToMakeDir)
			if err := os.Mkdir(pathAsDir, 0755); err != nil {
				t.Fatalf("Failed to create %s as directory: %v", tc.fileToMakeDir, err)
			}

			err = validateFrappeAppStructure(tmpDir, appName)
			if err == nil {
				t.Errorf("Expected error for %s being a directory, got nil", tc.fileToMakeDir)
			} else {
				expectedErrorPart := fmt.Sprintf("'%s' is a directory, not a file", pathAsDir)
				if !strings.Contains(err.Error(), expectedErrorPart) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", expectedErrorPart, err.Error())
				}
			}
		})
	}

	t.Run("valid complex mock app structure", func(t *testing.T) {
		tmpComplexAppDir, err := os.MkdirTemp("", "test-complex-app-")
		if err != nil {
			t.Fatalf("Failed to create temp dir for complex app: %v", err)
		}
		defer os.RemoveAll(tmpComplexAppDir)

		// The appName for validateFrappeAppStructure is the name of the inner directory,
		// which is conventionally the repository name or a derivative.
		// For this test, let's use a fixed appName for simplicity inside the temp source dir.
		complexAppName := "mycomplexapp"

		// This is the directory that contains the appName directory
		// So, sourceDir for validateFrappeAppStructure will be tmpComplexAppDir
		// and appName will be complexAppName.
		// Structure: tmpComplexAppDir / complexAppName / <-- (hooks.py etc here)

		innerAppActualDir := filepath.Join(tmpComplexAppDir, complexAppName)
		if err := os.Mkdir(innerAppActualDir, 0755); err != nil {
			t.Fatalf("Failed to create inner app dir '%s': %v", innerAppActualDir, err)
		}

		// 2.d. Create required files
		requiredFiles := []string{"__init__.py", "hooks.py", "modules.txt"}
		for _, fName := range requiredFiles {
			if err := createFile(filepath.Join(innerAppActualDir, fName)); err != nil {
				t.Fatalf("Failed to create required file %s: %v", fName, err)
			}
		}

		// 2.e. Create additional realistic structure
		// Public JS directory
		publicJsDir := filepath.Join(innerAppActualDir, "public", "js")
		if err := os.MkdirAll(publicJsDir, 0755); err != nil {
			t.Fatalf("Failed to create public/js dir: %v", err)
		}

		// Templates
		templatesPagesDir := filepath.Join(innerAppActualDir, "templates", "pages")
		if err := os.MkdirAll(templatesPagesDir, 0755); err != nil {
			t.Fatalf("Failed to create templates/pages dir: %v", err)
		}
		if err := createFile(filepath.Join(templatesPagesDir, "test_page.html")); err != nil {
			t.Fatalf("Failed to create test_page.html: %v", err)
		}

		// Doctype
		doctypeDirName := complexAppName + "_doctype" // e.g., mycomplexapp_doctype
		doctypeActualDir := filepath.Join(innerAppActualDir, doctypeDirName)
		if err := os.Mkdir(doctypeActualDir, 0755); err != nil {
			t.Fatalf("Failed to create doctype dir %s: %v", doctypeDirName, err)
		}
		doctypeJsonFileName := complexAppName + "_doctype.json" // e.g., mycomplexapp_doctype.json
		if err := createFile(filepath.Join(doctypeActualDir, doctypeJsonFileName)); err != nil {
			t.Fatalf("Failed to create %s: %v", doctypeJsonFileName, err)
		}

		// Root setup.py (this is outside the inner appName dir, in the sourceDir)
		if err := createFile(filepath.Join(tmpComplexAppDir, "setup.py")); err != nil {
			t.Fatalf("Failed to create setup.py: %v", err)
		}
		// A MANIFEST.in file
		if err := createFile(filepath.Join(tmpComplexAppDir, "MANIFEST.in")); err != nil {
			t.Fatalf("Failed to create MANIFEST.in: %v", err)
		}


		// 2.f. Call validateFrappeAppStructure
		// sourceDir is tmpComplexAppDir, appName is complexAppName
		err = validateFrappeAppStructure(tmpComplexAppDir, complexAppName)

		// 2.g. Assert that the error returned is nil
		if err != nil {
			t.Errorf("Expected no error for valid complex app structure, but got: %v", err)
		}
	})
}

// Note: This test file assumes that `validateFrappeAppStructure` is in the `cmd` package.
// If it's in a different package, the import path for `validateFrappeAppStructure` would need adjustment,
// but since they are in the same package `cmd`, direct calls are fine.
// The function `validateFrappeAppStructure` itself is not exported (lowercase 'v'),
// so this test file `package_test.go` must be in the same `cmd` package, which it is.
// Build constraints or tags are not needed as long as they are in the same package.

// To run these tests:
// cd <path_to_fpm_directory>/fpm/cmd
// go test

// To run with coverage:
// go test -coverprofile=coverage.out
// go tool cover -html=coverage.out
