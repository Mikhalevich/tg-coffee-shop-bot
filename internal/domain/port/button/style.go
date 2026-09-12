package button

type Style string

const (
	StyleDefault Style = ""
	StyleDanger  Style = "danger"
	StyleSuccess Style = "success"
	StylePrimary Style = "primary"
)

func (s Style) String() string {
	return string(s)
}

func StyleFromString(raw string) Style {
	switch raw {
	case StyleDefault.String(),
		StyleDanger.String(),
		StyleSuccess.String(),
		StylePrimary.String():
		return Style(raw)
	}

	return StyleDefault
}
