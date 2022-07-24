package gerror

import "fmt"

func NewIllegalParameterError(text string) error {
	return fmt.Errorf("Illegal Parameter %s", text)
}

func StatusChangeError(text string) error {
	return fmt.Errorf("Status change fail with %s", text)
}
