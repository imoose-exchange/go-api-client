package imoose

import "strings"

type ImooseError struct {
	Errors []string `json:"errors" `
}

func (m *ImooseError) Error() string {
	return strings.Join(m.Errors, ",")
}
