package button

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

const (
	shareURLTemplate = "https://t.me/share/url?url=%s"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(s string) ID {
	return ID(s)
}

type Button struct {
	ID                   ID
	Caption              string
	Operation            Operation
	IsDeleteAfterProcess bool
	Style                Style
	URL                  string
	Payload              []byte
}

type ButtonRow []Button

func Row(buttons ...Button) ButtonRow {
	return buttons
}

func GetPayload[P any](b Button) (P, error) {
	payload, err := gobDecodePayload[P](b.Payload)
	if err != nil {
		return payload, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

func gobEncodePayload(p any) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}

	return buf.Bytes(), nil
}

func gobDecodePayload[Payload any](b []byte) (Payload, error) {
	var payload Payload
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return payload, fmt.Errorf("gob decode: %w", err)
	}

	return payload, nil
}

func CreateButton(
	caption string,
	operation Operation,
	opts ...Option,
) (Button, error) {
	btn := Button{
		ID:        IDFromString(generateID()),
		Caption:   caption,
		Operation: operation,
		Style:     StyleDefault,
	}

	for _, opt := range opts {
		if err := opt(&btn); err != nil {
			return Button{}, fmt.Errorf("button option: %w", err)
		}
	}

	return btn, nil
}

func MustCreateButton(
	caption string,
	operation Operation,
	opts ...Option,
) Button {
	btn := Button{
		ID:        IDFromString(generateID()),
		Caption:   caption,
		Operation: operation,
		Style:     StyleDefault,
	}

	for _, opt := range opts {
		if err := opt(&btn); err != nil {
			panic(fmt.Errorf("button option: %w", err))
		}
	}

	return btn
}

func MustShareButton(
	caption string,
	targetURL string,
) Button {
	return MustCreateButton(
		caption,
		OperationOpenURL,
		WithURL(
			fmt.Sprintf(
				shareURLTemplate,
				url.PathEscape(targetURL),
			),
		),
	)
}

func generateID() string {
	return uuid.NewString()
}
