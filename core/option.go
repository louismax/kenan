package core

type KenanOption interface {
	Join(settings *Settings) error
}

type ErrorOption struct {
	Error error
}

// Join 返回初始化错误
func (w ErrorOption) Join(_ *Settings) error {
	return w.Error
}
