package artifacts

import "io"

type UploaderProvider func(projectID string, runID string) (Uploader, error)

type Uploader interface {
	Upload(name string, reader io.Reader) error
	Close() error
}
