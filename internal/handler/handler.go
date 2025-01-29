package handler

var (
	ERRUniqueConstraint    = "SQLSTATE 23505"
	ERRInternalServerError = "Internal Server Error"
	ERRRecordNotFound      = "record not found"
)

// Response to formalize http response data
type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
