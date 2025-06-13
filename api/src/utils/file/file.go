package file

import (
	"dse/src/utils/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Exists checks if a file path exists.
//
// Parameters:
//   - path: The filesystem path to check for existence.
//
// Returns:
//   - bool: True if the file exists, false otherwise.
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Create creates a file at the specified path.
//
// Parameters:
//   - path: The filesystem path where the file should be created.
//
// Returns:
//   - error: Any error that occurred while creating the file, or nil if successful.
func Create(path string) (*os.File, error) {
	return os.Create(path)
}

// Open opens a file for reading only.
//
// Parameters:
//   - path: The filesystem path to the file to open.
//
// Returns:
//   - *os.File: A pointer to the opened file.
//   - error: Any error that occurred while opening the file, or nil if successful.
func Open(path string) (*os.File, error) {
	return os.Open(path)
}

// ReadJSON reads a JSON file and decodes it into the provided interface.
//
// Parameters:
//   - path: The filesystem path to the JSON file to read.
//   - v: A pointer to the interface where the JSON data will be decoded.
//
// Returns:
//   - error: Any error that occurred while reading or decoding the file, or nil if successful.
func ReadJSON(path string, v interface{}) error {
	data, err := Read(path)
	if err != nil {
		return err
	}
	return json.FromBytes(data, v)
}

// Read reads the contents of a file and returns it as a byte slice.
//
// Parameters:
//   - path: The filesystem path to the file to read.
//
// Returns:
//   - []byte: The contents of the file.
//   - error: Any error that occurred while reading the file, or nil if successful.
func Read(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Write writes data to a file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file to write to.
//   - data: The data to write to the file.
//
// Returns:
//   - error: Any error that occurred while writing the file, or nil if successful.
func Write(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// WriteJSON writes a JSON object to a file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file to write to.
//   - data: The interface{} object to be marshaled into JSON and written to the file.
//
// Returns:
//   - error: Any error that occurred while marshaling the data or writing to the file, or nil if successful.
func WriteJSON(path string, data interface{}) error {
	b, err := json.ToBytes(data)
	if err != nil {
		return err
	}
	return Write(path, b)
}

// Delete deletes a file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file to delete.
//
// Returns:
//   - error: Any error that occurred while deleting the file, or nil if successful.
func Delete(path string) error {
	if IsFile(path) == false {
		return fmt.Errorf("Path is not a file")
	}
	return os.Remove(path)
}

// Rename renames or moves a file from oldpath to newpath.
//
// Parameters:
//   - oldpath: The current filesystem path of the file.
//   - newpath: The new filesystem path of the file.
//
// Returns:
//   - error: Any error that occurred while renaming the file, or nil if successful.
func Rename(oldpath string, newpath string) error {
	err := os.Rename(oldpath, newpath)
	if err != nil {
		return err
	}
	return nil
}

// Stream opens a file for reading and writing, creating it if it doesn't exist.
//
// Parameters:
//   - path: The filesystem path to the file to open/create.
//
// Returns:
//   - *os.File: A pointer to the opened file.
//   - error: Any error that occurred while opening the file, or nil if successful.
func Stream(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Size retrieves the size of the file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file whose size is to be determined.
//
// Returns:
//   - int64: The size of the file in bytes.
//   - error: Any error that occurred while retrieving the file size, or nil if successful.
func Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Info retrieves information about the file such as its modification date.
//
// Parameters:
//   - path: The filesystem path to the file.
//
// Returns:
//   - os.FileInfo: Information about the file.
//   - error: Any error that occurred while retrieving the file information, or nil if successful.
func Info(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetModified retrieves the modification time of the file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file.
//
// Returns:
//   - time.Time: The modification time of the file.
//   - error: Any error that occurred while retrieving the modification time, or nil if successful.
func GetModified(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

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

// IsFile checks if a path is a regular file.
//
// Parameters:
//   - path: The filesystem path to check.
//
// Returns:
//   - bool: True if the path is a regular file, false otherwise.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Append adds data to the end of a file at the specified path.
//
// Parameters:
//   - path: The filesystem path to the file to append to.
//   - data: The data to append to the file.
//
// Returns:
//   - error: Any error that occurred while appending to the file, or nil if successful.
func Append(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// AppendJSON reads a JSON file, appends new data to it, and writes it back.
//
// Parameters:
//   - path: The filesystem path to the JSON file.
//   - data: The data to append to the JSON array.
//
// Returns:
//   - error: Any error that occurred during the operation, or nil if successful.
func AppendJSON(path string, data interface{}) error {
	var err error
	var left  []byte
	var right []byte

	right, err = json.ToBytes(data)
	if err != nil {
		return err
	}

	if Exists(path) {
		left, err = Read(path)
		if err != nil {
			return err
		}

		keys, err := json.KeysBytes(left)
		if err != nil {
			return err
		}

		for _, k := range keys {
			v := json.GetBytes(left, k).Value()
			right, err = json.SetBytes(right, k, v)
			if err != nil {
				return err
			}
		}
	}

	right, err = json.FormatBytes(right)
	if err != nil {
		return err
	}

	return Write(path, right)
}


// WaitFor waits until a file becomes available using filesystem notifications.
//
// Parameters:
//   - path: The filesystem path to wait for.
//   - timeout: Maximum duration to wait in seconds (0 for infinite).
//
// Returns:
//   - error: Error if timeout is reached or other errors occur.
func WaitFor(path string, timeout int) error {
	// Create a watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Watch the directory containing the file
	dir := filepath.Dir(path)
	if err := watcher.Add(dir); err != nil {
		return err
	}

	// Set up timeout channel
	timeoutChan := make(chan bool, 1)
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			timeoutChan <- true
		}()
	}

	// Check if file already exists and is accessible
	if Exists(path) {
		if file, err := os.OpenFile(path, os.O_RDWR, 0644); err == nil {
			file.Close()
			return nil
		}
	}

	// Wait for events
	for {
		select {
		case event := <-watcher.Events:
			if event.Name == path {
				if Exists(path) {
					if file, err := os.OpenFile(path, os.O_RDWR, 0644); err == nil {
						file.Close()
						return nil
					}
				}
			}
		case err := <-watcher.Errors:
			return err
		case <-timeoutChan:
			return fmt.Errorf("timeout waiting for file: %s", path)
		}
	}
}
