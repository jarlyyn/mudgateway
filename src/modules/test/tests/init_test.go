package tests

import (
	"mudgateway/modules/app"
	"mudgateway/modules/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	test.Init()
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
func TestInit(t *testing.T) {
	if app.Development.Testing == false {
		t.Fatal(app.Development.Testing)
	}
}
