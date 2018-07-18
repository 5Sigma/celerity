package celerity

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/5Sigma/vox"
	"github.com/spf13/viper"
)

func EmptyRouteHandler(c Context) Response {
	return c.R(nil)
}

func TestRoutesCommand(t *testing.T) {

	rootCmd := setupCLI(func() *Server {
		svr := New()

		svr.GET("/test", EmptyRouteHandler)
		return svr
	})
	pl := vox.Test()
	rootCmd.SetArgs([]string{"routes"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Error running routes command: %s", err.Error())
	}
	expected := fmt.Sprintf(`All server routes:

%s[SCOPE]%s /
	GET	/test

`, vox.Yellow, vox.ResetColor) + "\n"

	if v := strings.Join(pl.LogLines, ""); v != expected {
		t.Errorf("incorrect output: \n'%s'\n'%s'", expected, v)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	os.Setenv("FOO", "bar")
	cliConfig()
	v := viper.GetString("foo")
	if v != "bar" {
		t.Errorf("environment reading not setup: %s", v)
	}
}
