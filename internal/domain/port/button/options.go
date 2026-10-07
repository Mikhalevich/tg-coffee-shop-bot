package button

import (
	"fmt"
)

type Option func(btn *Button) error

func WithDeleteAfterProcess() Option {
	return func(btn *Button) error {
		btn.IsDeleteAfterProcess = true

		return nil
	}
}

func WithStyle(style Style) Option {
	return func(btn *Button) error {
		btn.Style = style

		return nil
	}
}

func WithURL(url string) Option {
	return func(btn *Button) error {
		btn.URL = url

		return nil
	}
}

func WithPayload[T any](payload T) Option {
	return func(btn *Button) error {
		payloadBytes, err := gobEncodePayload(payload)
		if err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}

		btn.Payload = payloadBytes

		return nil
	}
}
