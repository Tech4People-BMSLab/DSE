package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsDirectory checks if a path is a directory.
//
// Parameters:
//   - path: The filesystem path to check.
//
// Returns:
//   - bool: True if the path is a directory, false otherwise.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Exists checks if a directory path exists.
//
// Parameters:
//   - path: The filesystem path to check for existence.
//
// Returns:
//   - bool: True if the directory exists, false otherwise.
func Exists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Create creates a directory at the specified path.
//
// Parameters:
//   - path: The filesystem path where the directory should be created.
//
// Returns:
//   - error: Any error that occurred while creating the directory, or nil if successful.
func Create(path string) error {
	return os.MkdirAll(path, 0755)
}

// Delete removes a directory at the specified path.
//
// Parameters:
//   - path: The filesystem path to the directory to delete.
//
// Returns:
//   - error: Any error that occurred while deleting the directory, or nil if successful.
func Delete(path string) error {
	if IsDirectory(path) == false {
		return fmt.Errorf("Path is not a directory")
	}
	return os.RemoveAll(path)
}

// Rename renames or moves a directory from oldpath to newpath.
//
// Parameters:
//   - oldpath: The current filesystem path of the directory.
//   - newpath: The new filesystem path of the directory.
//
// Returns:
//   - error: Any error that occurred while renaming the directory, or nil if successful.
func Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// ListFiles lists all files in a directory.
//
// Parameters:
//   - path: The filesystem path to the directory.
//
// Returns:
//   - []string: A slice of file paths found in the directory.
//   - error: Any error that occurred while reading the directory, or nil if successful.
func ListFiles(path string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return files, nil
}

// ListDirectories lists all subdirectories in a directory.
//
// Parameters:
//   - path: The filesystem path to the directory.
//
// Returns:
//   - []string: A slice of subdirectory paths found in the directory.
//   - error: Any error that occurred while reading the directory, or nil if successful.
func ListDirectories(path string) ([]string, error) {
	var directories []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(path, entry.Name()))
		}
	}
	return directories, nil
}

// IsEmpty checks if a directory is empty.
//
// Parameters:
//   - path: The filesystem path to the directory to check.
//
// Returns:
//   - bool: True if the directory is empty, false otherwise.
//   - error: Any error encountered while checking the directory.
func IsEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// Size calculates the total size of all files in a directory.
//
// Parameters:
//   - path: The filesystem path to the directory.
//
// Returns:
//   - int64: The total size of all files in the directory in bytes.
//   - error: Any error encountered while calculating the size.
func Size(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// Copy copies the contents of a directory to a new destination.
//
// Parameters:
//   - src: The source directory path.
//   - dst: The destination directory path.
//
// Returns:
//   - error: Any error encountered while copying the directory.
func Copy(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(targetPath, data, info.Mode())
		}
	})
}

// GetCurrentDirectory returns the current working directory.
//
// Returns:
//   - string: The current working directory path.
//   - error: Any error encountered while retrieving the current directory.
func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}
