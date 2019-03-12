package settings

// The file defines constants used by various handlers and forms, that affect
// the behavior of the API in a non-critical way (e.g. validation settings,
// entries per page etc.)

// Specifies the directory images are saved to
const UploadsPath = "media/images"

// Specifies the string returned to the client when the entry in the DB is NULL or empty
const DefaultImgName = "default.png"

// Defines how many users are returned per page by default.
const UsersDefaultPageSize = 10

// Defines the upperbound for the number of users returned per page.
const UsersMaxPageSize = 10

// Passwords weaker than MinPasswordStrength will raise a validation error.
// has to be in range of 0..4, or can cause a panic. Strength is estimated by
// zxcvbn.PasswordStrength(pwd).Score.
const MinPasswordStrength = 0

// If true, also checks that the email is real when signing up. This can be
// rather long (up to 1 second) since it requires an answer from the SMTP server.
const EmailExistsCheck = false

// Lifespan of a JWE Auth token - in days
const JWETokenLifespan = 30

// Regulates the length of a CSRF Token (in bytes). 20 is probably ok.
const CSRFTokenLength = 20

// Lifespan of a CSRFToken - in days
const CSRFTokenLifespan = 7
