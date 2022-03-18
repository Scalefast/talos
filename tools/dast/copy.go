package dast

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyEmbedFile(src, dst string) (err error) {
	in, err := helpers.ReadFile(src)
	if err != nil {
		return err
	}

	var permissions fs.FileMode = 0755
	err = ioutil.WriteFile(dst, in, permissions)
	if err != nil {
		return err
	}
	return nil
}

// Copy src into dst
// src has to be of type embed.FS
// dst has to be of type string
func CopyEmbedDir(src embed.FS, dst string) error {
	dst = filepath.Clean(dst)
	err := fs.WalkDir(src, ".", func(file string, f fs.DirEntry, err error) error {

		// Error check
		if err != nil {
			return err
		}
		// If the file is a directory
		if f.IsDir() {
			nd := filepath.Clean(dst + "/" + file)
			// Create the dst directories
			var permissions fs.FileMode = 0766
			err = os.MkdirAll(nd, permissions)
			if err != nil {
				return err
			}
			// If the file is not a directory
		} else {
			err = CopyEmbedFile(file, dst+"/"+file)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
