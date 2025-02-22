package cmdutil

import (
	"os"
	"testing"

	"github.com/vercel/turbo/cli/internal/config"

	"github.com/spf13/pflag"
	"github.com/vercel/turbo/cli/internal/fs"
	"gotest.tools/v3/assert"
)

func TestTokenEnvVar(t *testing.T) {

	// Set up an empty config so we're just testing environment variables
	userConfigPath := fs.AbsoluteSystemPathFromUpstream(t.TempDir()).UntypedJoin("turborepo", "config.json")
	expectedPrefix := "my-token"
	vars := []string{"TURBO_TOKEN", "VERCEL_ARTIFACTS_TOKEN"}
	for _, v := range vars {
		t.Run(v, func(t *testing.T) {
			t.Cleanup(func() {
				_ = os.Unsetenv(v)
			})
			flags := pflag.NewFlagSet("test-flags", pflag.ContinueOnError)
			h := NewHelper("test-version")
			h.AddFlags(flags)
			h.UserConfigPath = userConfigPath

			expectedToken := expectedPrefix + v
			err := os.Setenv(v, expectedToken)
			if err != nil {
				t.Fatalf("setenv %v", err)
			}

			base, err := h.GetCmdBase(config.FlagSet{FlagSet: flags})
			if err != nil {
				t.Fatalf("failed to get command base %v", err)
			}
			assert.Equal(t, base.RemoteConfig.Token, expectedToken)
		})
	}
}
