package gerror

import "fmt"

func NewIllegalParameterError(text string) error {
	return fmt.Errorf("Illegal Parameter %s", text)
}
