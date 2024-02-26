package pipelines

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var pipelines = make(map[string]*Pipeline, 0)

func LoadPipelines(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		id := fileNameWithoutExt(entry.Name())
		_, found := pipelines[id]
		if found {
			log.Printf("Project with ID \"%s\" is already loaded", id)
			continue
		}

		filePath := path.Join(directory, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		var pipeline *Pipeline
		err = yaml.Unmarshal(bytes, &pipeline)
		if err != nil {
			return err
		}

		log.Printf("Loaded pipeline \"%s\"", id)
		pipelines[id] = pipeline
	}

	return nil
}

func GetPipelines() map[string]*Pipeline {
	return pipelines
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
