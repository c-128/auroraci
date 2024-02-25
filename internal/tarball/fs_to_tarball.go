package tarball

import (
	"archive/tar"
	"errors"
	"io"
	"io/fs"

	"github.com/go-git/go-billy/v5"
)

func FSToTarball(
	fs billy.Filesystem,
	root string,
	writer io.WriteCloser,
) error {
	tarball := tar.NewWriter(writer)
	defer writer.Close()

	err := handleDir(
		fs,
		root,
		tarball,
	)
	if err != nil {
		return err
	}

	return nil
}

func handleDir(
	fs billy.Filesystem,
	root string,
	tarball *tar.Writer,
) error {
	files, err := fs.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := fs.Join(root, file.Name())

		header, isDir, err := fileInfoToTarHeader(
			file,
			fullPath,
		)
		if err != nil {
			return err
		}

		err = tarball.WriteHeader(header)
		if err != nil {
			return err
		}

		if isDir {
			err := handleDir(
				fs,
				fullPath,
				tarball,
			)
			if err != nil {
				return err
			}

			continue
		}

		file, err := fs.Open(fullPath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarball, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func fileInfoToTarHeader(file fs.FileInfo, name string) (*tar.Header, bool, error) {
	fileMode := file.Mode()

	header := &tar.Header{
		Name:    name,
		ModTime: file.ModTime(),
		Mode:    int64(fileMode.Perm()),
	}
	isDir := false

	switch {
	case fileMode.IsRegular():
		header.Typeflag = tar.TypeReg
		header.Size = file.Size()
	case fileMode.IsDir():
		header.Typeflag = tar.TypeDir
		isDir = true
	default:
		return nil, false, errors.New("can not handle file")
	}

	return header, isDir, nil
}
