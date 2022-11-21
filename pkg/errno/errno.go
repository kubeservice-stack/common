package errno

import (
	"fmt"
)

type Errno struct {
	status  int    `errno:"errno"`
	message string `json:"errmsg"`
}

func (e *Errno) Status() int {
	return e.status
}

func (e *Errno) Message() string {
	return e.message
}

func (e *Errno) Error() string {
	return fmt.Sprintf("Error - errno: %d, errmsg: %s", e.status, e.message)
}

func New(status int, message string) *Errno {
	return &Errno{status: status, message: message}
}

func NewCode(status int) *Errno {
	return &Errno{status: status}
}
