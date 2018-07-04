package celerity

import (
	"fmt"
	"os"
	"testing"

	"github.com/5Sigma/vox"
	"github.com/spf13/viper"
)

func EmptyRouteHandler(c Context) Response {
	return c.R(nil)
}

func TestRoutesCommand(t *testing.T) {
	svr := New()

	svr.GET("/test", EmptyRouteHandler)

	rootCmd := setupCLI(svr)
	vox.Test()
	rootCmd.SetArgs([]string{"routes"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Error running routes command: %s", err.Error())
	}
	vox.AssertOutput(t, "All server routes:\n")
	vox.AssertOutput(t, "\n")
	vox.AssertOutput(t, fmt.Sprint(vox.Yellow, "[SCOPE]", vox.ResetColor, " /\n"))
	vox.AssertOutput(t, "\tGET\t/test\n")
}

func TestEnvironmentVariables(t *testing.T) {
	os.Setenv("FOO", "bar")
	cliConfig()
	v := viper.GetString("foo")
	if v != "bar" {
		t.Errorf("environment reading not setup: %s", v)
	}
}
