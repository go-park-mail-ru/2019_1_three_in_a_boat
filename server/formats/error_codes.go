package formats

// These constants represent error messages - not error objects. They represent
// strings that are sent to client when a certain error occurs, since sending
// the whole error.Error() is probably unwise.
const (
	Err404                   = "E_NOT_FOUND"
	Err405                   = "E_METHOD_NOT_ALLOWED"
	Err403                   = "E_FORBIDDEN"
	ErrCSRF                  = "E_CSRF_VALIDATION_FAILED"
	ErrDbModuleFailure       = "E_DB_MODULE_FAILURE"
	ErrDbRowsFailure         = "E_DB_READ_FAILURE"
	ErrDbScanFailure         = "E_DB_SCAN_FAILURE"
	ErrDbTransactionFailure  = "E_DB_TRANSACTION_FAILURE"
	ErrJSONMarshalFailure    = "E_JSON_FAILURE"
	ErrJWTDecryptionEmpty    = "E_JWT_EMPTY_FAILURE"
	ErrJWTDecryptionFailure  = "E_JWT_DECRYPTION_FAILURE"
	ErrJWTEncryptionFailure  = "E_JWT_ENCRYPTION_FAILURE"
	ErrJWTOutdated           = "E_JWT_OUTDATED"
	ErrPasswordHashing       = "E_BCRYPT_FAILURE"
	ErrPanic                 = "E_PANIC"
	ErrResponseWriterFailure = "E_RESPONSE_WRITER_FAILURE"
	ErrSavingImg             = "E_IMG_SAVE_FAILURE"
	ErrSignupAuthFailure     = "E_SIGNUP_OK_AUTH_FAILED"
	ErrSqlFailure            = "E_DB_QUERY_FAILURE"
)

const (
	ErrDateOutOfRange          = "E_DATE_OUT_OF_RANGE"
	ErrBase64Decoding          = "E_INVALID_BASE64"
	ErrEmailDoesNotExist       = "E_EMAIL_DOES_NOT_EXIST"
	ErrEmailDomainDoesNotExist = "E_EMAIL_DOMAIN_DOES_NOT_EXIST"
	ErrFieldTooLong            = "E_TOO_LONG"
	ErrFieldTooShort           = "E_TOO_SHORT"
	ErrInvalidCredentials      = "E_INVALID_CREDENTIALS"
	ErrInvalidDate             = "E_INVALID_DATE_FORMAT"
	ErrInvalidEmail            = "E_INVALID_EMAIL"
	ErrInvalidJSON             = "E_INVALID_JSON"
	ErrPasswordTooWeak0        = "E_WEAK_PASSWORD_0"
	ErrPasswordTooWeak1        = "E_WEAK_PASSWORD_1"
	ErrPasswordTooWeak2        = "E_WEAK_PASSWORD_2"
	ErrPasswordTooWeak3        = "E_WEAK_PASSWORD_3"
	ErrUniqueViolation         = "E_DUPLICATE_NOT_ALLOWED"
	ErrValidation              = "E_VALIDATION_FAILED"
)

var ErrPasswordTooWeak = map[int]string{
	0: ErrPasswordTooWeak0,
	1: ErrPasswordTooWeak1,
	2: ErrPasswordTooWeak2,
	3: ErrPasswordTooWeak3,
}
