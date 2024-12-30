package student

var EmailAlreadyRegister = &UserErrors{msg: "Email ya registrado"}
var NumberAlreadyRegister = &UserErrors{msg: "Numero ya registrado"}
var UserNotFound = &UserErrors{msg: "Usuario no encontrado"}
var EmailOrNumberAlreadyRegister = &UserErrors{msg: "Email o numero ya registrado"}

var CourseNotValid = &UserErrors{msg: "Courses Not valid"}

type UserErrors struct {
	msg string
}

// Error implements error.
func (e *UserErrors) Error() string {
	return e.msg
}
