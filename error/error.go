package error

type Error interface {
	error
	Cause() error
}
