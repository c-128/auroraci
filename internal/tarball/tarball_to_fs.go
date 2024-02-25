package tarball

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

func TarballToFS(
	reader io.ReadCloser,
	root string,
	fs billy.Filesystem,
	prefixToStrip string,
) error {
	tarball := tar.NewReader(reader)
	defer reader.Close()

	for {
		header, err := tarball.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := fs.Join(root, header.Name)
		path, err = filepath.Rel(prefixToStrip, path)
		if err != nil {
			return err
		}

		fileInfo := header.FileInfo()

		if fileInfo.IsDir() {
			err = fs.MkdirAll(path, fileInfo.Mode())
			if err != nil {
				return err
			}

			continue
		}

		file, err := fs.OpenFile(
			path,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			fileInfo.Mode(),
		)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarball)
		if err != nil {
			return err
		}
	}

	return nil
}
