package messagesvc

func (s *Service) EscapeMarkdown(text string) string {
	return s.escaper.EscapeMarkdown(text)
}
