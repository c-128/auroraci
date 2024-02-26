package pipelines

type Pipeline struct {
	Repository *Repository
	Build      *Build
}

type Repository struct {
	Origin string
	Branch string
}

type Build struct {
	Triggers  []BuildTrigger
	Stages    []BuildStage
	Artifacts []string
}

type BuildTrigger struct {
	Cron string
}

type BuildStage struct {
	Image    string
	Workdir  string
	Commands []BuildCommand
}

type BuildCommand struct {
	RunBash  string
	Run      string
	ExitCode int
}
