package endpoint_inputs

type CommonErrorOutput struct {
	Status bool `json:"status"`
	ErrorMsg string	`json:"msg"`
}

type VerifyUserEmailInput struct {
	Email string     `validate:"required,email"`
}