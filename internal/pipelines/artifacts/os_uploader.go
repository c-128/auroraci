package artifacts

import (
	"io"
	"os"
	"path"
)

func NewOSProvider(artifactsDirectory string) UploaderProvider {
	providerFunc := func(pipelineID, runID string) (Uploader, error) {
		projectDirectory := path.Join(artifactsDirectory, pipelineID)

		err := os.MkdirAll(projectDirectory, os.ModePerm)
		if err != nil {
			return nil, err
		}

		runDirectory := path.Join(projectDirectory, runID)
		err = os.MkdirAll(runDirectory, os.ModePerm)
		if err != nil {
			return nil, err
		}

		return osUploader(runDirectory), nil
	}

	return providerFunc
}

type osUploader string

func (uploader osUploader) Upload(name string, reader io.Reader) error {
	filePath := path.Join(string(uploader), name)
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

func (uploader osUploader) Close() error {
	return nil
}
