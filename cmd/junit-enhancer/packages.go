package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Absolute path to root directory -> PackageInfo
var (
	// ErrPackageNotFound is returned when a package is not found in the go code.
	ErrPackageNotFound = errors.New("package not found")
)

// PackagesInfo contains all the packages found in a Go project.
type PackagesInfo struct {
	Packages map[string]PackageInfo
}

// Get returns the PackageInfo for the given import path.
// Note that the import path is the full import path, not just the package name.
func (p *PackagesInfo) Get(importPath string) (PackageInfo, PackageInfo, error) {
	pkg, ok := p.Packages[importPath]
	testPkg, okTest := p.Packages[importPath+"_test"]

	if ok && okTest {
		return pkg, testPkg, nil
	}
	if ok {
		return pkg, PackageInfo{}, nil
	}

	all := make([]string, 0, len(p.Packages))
	for k := range p.Packages {
		all = append(all, k)
	}

	return PackageInfo{}, PackageInfo{}, fmt.Errorf(
		"%w: %s\nall packages:\n%s",
		ErrPackageNotFound,
		importPath,
		strings.Join(all, "\n"),
	)
}

// FileInfo ties an absolute (full) path to its repository-root-relative path.
type FileInfo struct {
	FullPath string // Absolute path to the file
	RelPath  string // Path relative to the repository root (forward slashes)
}

// PackageInfo contains comprehensive information about a Go package
type PackageInfo struct {
	ImportPath   string           // Package import path (e.g., "github.com/user/repo/pkg")
	Name         string           // Package name
	Dir          string           // Directory containing the package (absolute)
	GoFiles      []FileInfo       // .go source files
	TestGoFiles  []FileInfo       // _test.go files
	XTestGoFiles []FileInfo       // _test.go files with different package names (external tests)
	Module       *packages.Module // Module information for the package
	IsCommand    bool             // True if this is a main package
}

// Packages finds all Go packages in the given Go project directory and subdirectories.
// This includes packages of nested Go projects by discovering all go.mod files and
// loading packages from each module root.
// If modulePath is provided, only the module at that path (relative to rootDir) will be processed.
func Packages(rootDir, modulePath string) (*PackagesInfo, error) {
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	// Find all go.mod files to identify all Go modules in the directory tree
	var goModDirs []string

	if modulePath != "" {
		// If module path is specified, only process that specific module
		targetModuleDir := filepath.Join(absRootDir, modulePath)
		if _, statErr := os.Stat(filepath.Join(targetModuleDir, "go.mod")); statErr != nil {
			if os.IsNotExist(statErr) {
				return nil, fmt.Errorf("no go.mod file found at specified module path: %s", modulePath)
			}
			return nil, fmt.Errorf("failed to check go.mod at module path %s: %w", modulePath, statErr)
		}
		goModDirs = []string{targetModuleDir}
	} else {
		// Find all go.mod files in the directory tree
		goModDirs, err = findGoModDirectories(absRootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to find go.mod files: %w", err)
		}
	}

	if len(goModDirs) == 0 {
		return &PackagesInfo{Packages: make(map[string]PackageInfo)}, nil
	}

	result := &PackagesInfo{
		Packages: make(map[string]PackageInfo),
	}

	// Load packages from each Go module
	for _, modDir := range goModDirs {
		modulePackages, err := loadPackagesFromModule(modDir, absRootDir) // pass repo root
		if err != nil {
			return nil, fmt.Errorf("failed to load packages from module %s: %w", modDir, err)
		}
		// Merge packages from this module into the result
		maps.Copy(result.Packages, modulePackages)
	}

	return result, nil
}

// findGoModDirectories recursively finds all directories containing go.mod files
func findGoModDirectories(rootDir string) ([]string, error) {
	var goModDirs []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and vendor directories
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor") {
			return filepath.SkipDir
		}

		if info.Name() == "go.mod" {
			goModDirs = append(goModDirs, filepath.Dir(path))
		}

		return nil
	})

	return goModDirs, err
}

// loadPackagesFromModule loads all packages from a single Go module
func loadPackagesFromModule(moduleDir, repoRoot string) (map[string]PackageInfo, error) {
	config := &packages.Config{
		Mode:  packages.NeedName | packages.NeedModule | packages.NeedFiles,
		Dir:   moduleDir,
		Tests: true,
	}

	pkgs, err := packages.Load(config, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages from module %s: %w", moduleDir, err)
	}

	result := make(map[string]PackageInfo)

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, pkgErr := range pkg.Errors {
				fmt.Fprintf(os.Stderr, "Warning: package load error in %s: %v\n", pkg.PkgPath, pkgErr)
			}
		}

		info := PackageInfo{
			ImportPath: pkg.PkgPath,
			Name:       pkg.Name,
			IsCommand:  pkg.Name == "main",
			Module:     pkg.Module,
		}

		// Separate regular files from test files based on suffix
		goFiles, testFiles := splitGoAndTest(repoRoot, pkg.GoFiles)
		info.GoFiles = goFiles
		info.TestGoFiles = testFiles

		if len(pkg.GoFiles) > 0 {
			info.Dir = filepath.Dir(pkg.GoFiles[0])
		}

		result[info.ImportPath] = info
	}

	return result, nil
}

// splitGoAndTest converts a list of files to FileInfo and separates regular and _test files.
func splitGoAndTest(repoRoot string, files []string) (goFiles, testFiles []FileInfo) {
	for _, f := range files {
		fi := makeFileInfo(repoRoot, f)
		if strings.HasSuffix(f, "_test.go") {
			testFiles = append(testFiles, fi)
		} else {
			goFiles = append(goFiles, fi)
		}
	}
	return
}

// makeFileInfo builds a FileInfo with absolute and repository-root-relative paths.
// Falls back to absolute path for RelPath if it cannot be made relative.
// Normalizes RelPath to forward slashes for stability across OSes.
func makeFileInfo(repoRoot, abs string) FileInfo {
	rel, err := filepath.Rel(repoRoot, abs)
	if err != nil || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return FileInfo{FullPath: abs, RelPath: abs}
	}
	return FileInfo{FullPath: abs, RelPath: filepath.ToSlash(rel)}
}
