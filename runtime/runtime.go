package runtime

import "runtime"

func GetCallPath() string {
	_, path, _, _ := runtime.Caller(1)
	return path
}