package logging

import (
	"testing"
)

func Test_zapLogger_SetLevel(t *testing.T) {
	// p := GetProductionLogger()
	d := GetDevelopmentLogger()

	d.Debugf("q123")
	d.SetLevel(DebugLevel)
	d.Debugf("q123")
}
