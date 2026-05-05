package monitoring

type ErrorType string

const (
	ErrBadRequest    ErrorType = "bad_request_error"
	ErrExtractUserId ErrorType = "extract_user_id_error"
	ErrParseUUID     ErrorType = "parse_uuid_error"
	ErrForbidden     ErrorType = "forbidden_error"
	ErrInternal      ErrorType = "internal_error"
	ErrInvalidInput  ErrorType = "invalid_input"
	ErrBusinessLogic ErrorType = "business_logic_error"
)

func (et ErrorType) String() string {
	return string(et)
}
