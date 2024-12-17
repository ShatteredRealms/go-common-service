package bus

import "errors"

var (
	ErrDecodingMessage = errors.New("unable to decode message")
	ErrSendingMessage  = errors.New("unable to send message on bus")
	ErrEncodingMessage = errors.New("unable to encode message")
)
