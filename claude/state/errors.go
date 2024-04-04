package state

type ErrMissingFPQN struct {
	Err error
	Msg string
}

func (e *ErrMissingFPQN) Error() string {
	if e.Msg != "" {
		e.Msg = "missing FPQN- use WithFPQN to set the FPQN"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFileNotFound struct {
	Err      error
	Msg      string
	Filepath string
}

func (e *ErrFileNotFound) Error() string {
	if e.Msg != "" {
		e.Msg = "file not found"
	}
	if e.Filepath != "" {
		e.Filepath = "file does not exist: " + e.Filepath
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrOpeningFile struct {
	Err error
	Msg string
}

func (e *ErrOpeningFile) Error() string {
	if e.Msg != "" {
		e.Msg = "error opening file"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrLoadingGOB struct {
	Err error
	Msg string
}

func (e *ErrLoadingGOB) Error() string {
	if e.Msg != "" {
		e.Msg = "error loading GOB"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrSavingGOB struct {
	Err error
	Msg string
}

func (e *ErrSavingGOB) Error() string {
	if e.Msg != "" {
		e.Msg = "error saving GOB"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
