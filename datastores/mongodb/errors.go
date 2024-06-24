package mongodb

type ErrMissingMongoDbURI struct {
	Err error
	Msg string
}

func (e *ErrMissingMongoDbURI) Error() string {
	if e.Msg != "" {
		e.Msg = "MongoDB URI not set- use WithMongoDbURI() to set the URI"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrMissingDatabaseName struct {
	Err error
	Msg string
}

func (e *ErrMissingDatabaseName) Error() string {
	if e.Msg != "" {
		e.Msg = "database name is required"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToConnectToMongoDB struct {
	Err error
	Msg string
}

func (e *ErrFailedToConnectToMongoDB) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to connect to MongoDB"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

// failed to unmarshal request
type ErrFailedToUnmarshalRequest struct {
	Err error
	Msg string
}

func (e *ErrFailedToUnmarshalRequest) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to unmarshal request"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrActionMissingOrInvalid struct {
	Err error
	Msg string
}

func (e *ErrActionMissingOrInvalid) Error() string {
	if e.Msg != "" {
		e.Msg = "action is required and must be a string"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrActionUnknown struct {
	Err    error
	Msg    string
	Action string
}

func (e *ErrActionUnknown) Error() string {
	if e.Msg != "" {
		e.Msg = "unknown action"
	}
	if e.Action != "" {
		e.Msg = e.Msg + ": " + e.Action
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrRecordTypeMissingOrInvalid struct {
	Err error
	Msg string
}

func (e *ErrRecordTypeMissingOrInvalid) Error() string {
	if e.Msg != "" {
		e.Msg = "record_type is required and must be an object"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrDataMissingOrInvalid struct {
	Err error
	Msg string
}

func (e *ErrDataMissingOrInvalid) Error() string {
	if e.Msg != "" {
		e.Msg = "data is required and must be a string"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToInsertRecord struct {
	Err error
	Msg string
}

func (e *ErrFailedToInsertRecord) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to insert record"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrIdMissingOrInvalid struct {
	Err error
	Msg string
}

func (e *ErrIdMissingOrInvalid) Error() string {
	if e.Msg != "" {
		e.Msg = "id is required and must be a string"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrRecordNotFound struct {
	Err error
	Msg string
}

func (e *ErrRecordNotFound) Error() string {
	if e.Msg != "" {
		e.Msg = "record not found"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToRetrieveRecord struct {
	Err error
	Msg string
}

func (e *ErrFailedToRetrieveRecord) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to retrieve record"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToUpdateRecord struct {
	Err error
	Msg string
}

func (e *ErrFailedToUpdateRecord) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to update record"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToDeleteRecord struct {
	Err error
	Msg string
}

func (e *ErrFailedToDeleteRecord) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to delete record"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrQueryMissingOrInvalid struct {
	Err error
	Msg string
}

func (e *ErrQueryMissingOrInvalid) Error() string {
	if e.Msg != "" {
		e.Msg = "query is required and must be an object"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToSearchRecords struct {
	Err error
	Msg string
}

func (e *ErrFailedToSearchRecords) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to search records"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToDecodeSearchResults struct {
	Err error
	Msg string
}

func (e *ErrFailedToDecodeSearchResults) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to decode search results"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToExecuteInventoryAggregation struct {
	Err error
	Msg string
}

func (e *ErrFailedToExecuteInventoryAggregation) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to execute inventory aggregation"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

type ErrFailedToDecodeInventoryResults struct {
	Err error
	Msg string
}

func (e *ErrFailedToDecodeInventoryResults) Error() string {
	if e.Msg != "" {
		e.Msg = "failed to decode inventory results"
	}
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
