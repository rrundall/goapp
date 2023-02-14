package config

const (
	// DB error messagess
	DBConnectErrMsg   = "failed to connect to database"
	DBOperationErrMsg = "request failed. Verify data meets any requirements " +
		"(i.e. uniqueness, null, etc...) and try again"
	DBCloseErrMsg = "fail to close database"

	// Operation error messages
	InvalidDataErrMsg         = "invalid data passing."
	UnsupportedContentType    = "unsupported Content-Type."
	DataCouldNotBeEmptyErrMsg = "field is empty or not define.  Please fill out all required fields"
	FailToSaveLogErrMsg       = "fail to save log:"
	BadRequestErrMsg          = "bad Request. Please check your relative path"

	// Operation warning messages
	FieldsBeEmptyWarningMsg     = "following fields were not included in the update:"
	NoDataUpdateWarningMsg      = "no data update"
	NoQueryDataPassedWarningMsg = "no data to pass in query. All string fields were empty and/or book_id is 0"

	// Operation success messages
	HomepageMsg      = "Welcome To Book Library!"
	AddSuccessMsg    = "Data successfully added."
	UpdateSuccessMsg = "Data successfully updated."
	DeleteSuccessMsg = "Data successfully deleted."
)
