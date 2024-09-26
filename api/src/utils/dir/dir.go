package dir

import (
	"io"
	"os"
)

// Exists checks if a directory exists
func Exists(dirPath string) bool {
    _, err := os.Stat(dirPath)
    return !os.IsNotExist(err)
}

// Rename renames a directory
func Rename(oldName, newName string) error {
    return os.Rename(oldName, newName)
}

// Move moves a directory
func Move(oldLocation, newLocation string) error {
    return os.Rename(oldLocation, newLocation)
}

// Copy copies a directory
func Copy(srcDir, dstDir string) error {
    entries, err := os.ReadDir(srcDir)
    if err != nil {
        return err
    }

    os.MkdirAll(dstDir, 0755)

    for _, entry := range entries {
        srcPath := srcDir + "/" + entry.Name()
        dstPath := dstDir + "/" + entry.Name()

        if entry.IsDir() {
            err = Copy(srcPath, dstPath)
            if err != nil {
                return err
            }
        } else {
            in, err := os.Open(srcPath)
            if err != nil {
                return err
            }
            defer in.Close()

			out, err := os.Create(dstPath)
            if err != nil {
                return err
            }
            defer out.Close()

            _, err = io.Copy(out, in)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func Create(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}
