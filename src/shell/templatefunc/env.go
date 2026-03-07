package templatefunc

import "os"

func Env() Provider {
	return New("env", os.Getenv)
}
