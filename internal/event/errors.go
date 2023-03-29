package event

import "strings"

type Errors []error

func (e Errors) Error() string {
	var errMsg strings.Builder
	for i := range e {
		errMsg.WriteString(e[i].Error())
		errMsg.WriteString("\n")
	}
	return errMsg.String()
}
