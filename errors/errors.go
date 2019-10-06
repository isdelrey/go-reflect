package errors

type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

func New(s string) *Error {
	return &Error{
		s,
	}
}
