package artifacts

import (
	"io"
	"os"
	"path"
)

func NewOSProvider(artifactsDirectory string) UploaderProvider {
	providerFunc := func(projectID, runID string) (Uploader, error) {
		projectDirectory := path.Join(artifactsDirectory, projectID)

		err := os.MkdirAll(projectDirectory, os.ModePerm)
		if err != nil {
			return nil, err
		}

		runDirectory := path.Join(projectDirectory, runID)
		err = os.MkdirAll(runDirectory, os.ModePerm)
		if err != nil {
			return nil, err
		}

		uploader := &osUploader{
			runDirectory: runDirectory,
		}
		return uploader, nil
	}

	return providerFunc
}

type osUploader struct {
	runDirectory string
}

func (o *osUploader) Upload(name string, reader io.Reader) error {
	filePath := path.Join(o.runDirectory, name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (o *osUploader) Close() error {
	return nil
}
