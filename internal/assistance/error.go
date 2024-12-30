package assistance

var NotFountRegister = &AssistanceError{msg: "Registro no encontrado"}

type AssistanceError struct {
	msg string
}

// Error implements error.
func (e *AssistanceError) Error() string {
	return e.msg
}
