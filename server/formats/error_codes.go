package formats

const (
	ErrSqlFailure            = "E_DB_QUERY_FAILURE"
	ErrDbScanFailure         = "E_DB_SCAN_FAILURE"
	ErrDbRowsFailure         = "E_DB_READ_FAILURE"
	ErrDbModuleFailure       = "E_DB_MODULE_FAILURE"
	ErrDbTranscationFailure  = "E_DB_TRANSACTION_FAILURE"
	ErrResponseWriterFailure = "E_RESPONSE_WRITER_FAILURE"
	ErrJSONMarshalFailure    = "E_JSON_FAILURE"
	Err405                   = "E_METHOD_NOT_ALLOWED"
	Err404                   = "E_NOT_FOUND"
	ErrValidation            = "E_VALIDATION_FAILED"
	ErrInvalidJSON           = "E_INVALID_JSON"
)

const (
	ErrFieldTooShort           = "E_TOO_SHORT"
	ErrFieldTooLong            = "E_TOO_LONG"
	ErrPasswordTooWeak0        = "E_WEAK_PASSWORD_0"
	ErrPasswordTooWeak1        = "E_WEAK_PASSWORD_1"
	ErrPasswordTooWeak2        = "E_WEAK_PASSWORD_2"
	ErrPasswordTooWeak3        = "E_WEAK_PASSWORD_3"
	ErrInvalidEmail            = "E_INVALID_EMAI:"
	ErrEmailDoesNotExist       = "E_EMAIL_DOES_NOT_EXIST"
	ErrEmailDomainDoesNotExist = "E_EMAIL_DOMAIN_DOES_NOT_EXIST"
	ErrInvalidDate             = "E_INVALID_DATE_FORMAT"
	ErrDateOutOfRange          = "E_DATE_OUT_OF_RANGE"
	ErrUniqueViolation         = "E_DUPLICATE_NOT_ALLOWED"
)

var ErrPasswordTooWeak = map[int]string{
	0: ErrPasswordTooWeak0,
	1: ErrPasswordTooWeak1,
	2: ErrPasswordTooWeak2,
	3: ErrPasswordTooWeak3,
}
