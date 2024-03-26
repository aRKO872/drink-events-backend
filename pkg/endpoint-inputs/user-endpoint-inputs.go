package endpoint_inputs

type CommonErrorOutput struct {
	Status bool `json:"status"`
	ErrorMsg string	`json:"error-msg"`
}

type VerifyUserEmailInput struct {
	Email string     `json:"email" validate:"required,email"`
}

type VerifyOTP struct {
	Event string `json:"event" validate:"required,oneof=verify_email verify_phone login_email login_phone"`
	Otp		int		 `json:"otp" validate:"required,gte=100000,lte=999999"`
	Email string `json:"email" validate:"email,omitempty"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type LogInInput struct {
	LoggedInFrom string `json:"logged_in_from" validate:"required,oneof=email phone"`
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type ResendOTP struct {
	Event string `json:"event" validate:"required,oneof=verify_email verify_phone login_email login_phone"`
	Email string `json:"email" validate:"email,omitempty"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10,regex=^[6-9][0-9]*$"`
}

type SignUpInput struct {
	Name string `json:"name" validate:"min=3"`
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone" validate:"min=10,max=10,regex=^[6-9][0-9]*$"`
	Bio string `json:"bio" validate:"omitempty"`
}

type SignUpLoginOutput struct {
	Status bool `json:"status"`
	ErrorMsg string `json:"error-msg"`
	AccessToken string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}