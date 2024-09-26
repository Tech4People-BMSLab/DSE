package file

import (
	"io"
	"os"
)

// Exists checks if a file exists
func Exists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

// Rename renames a file
func Rename(oldName, newName string) error {
    return os.Rename(oldName, newName)
}

// Move moves a file
func Move(oldLocation, newLocation string) error {
    return os.Rename(oldLocation, newLocation)
}

// Copy copies a file
func Copy(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    return err
}
