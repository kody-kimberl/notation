package experimental

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	envName = "NOTATION_EXPERIMENTAL"
	enabled = "1"
)

// IsDisabled determines whether an experimental feature is disabled.
func IsDisabled() bool {
	return os.Getenv(envName) != enabled
}

// CheckCommandAndWarn checks whether an experimental command can be run.
func CheckCommandAndWarn(cmd *cobra.Command, _ []string) error {
	return CheckAndWarn(func() (string, bool) {
		return fmt.Sprintf("%q", cmd.CommandPath()), true
	})
}

// CheckFlagsAndWarn checks whether experimental flags can be run.
func CheckFlagsAndWarn(cmd *cobra.Command, flags ...string) error {
	return CheckAndWarn(func() (string, bool) {
		var changedFlags []string
		flagSet := cmd.Flags()
		for _, flag := range flags {
			if flagSet.Changed(flag) {
				changedFlags = append(changedFlags, "--"+flag)
			}
		}
		if len(changedFlags) == 0 {
			// no experimental flag used
			return "", false
		}
		return fmt.Sprintf("flag(s) %s in %q", strings.Join(changedFlags, ","), cmd.CommandPath()), true
	})
}

// CheckAndWarn checks whether a feature can be used.
func CheckAndWarn(doCheck func() (feature string, isExperimental bool)) error {
	feature, isExperimental := doCheck()
	if isExperimental {
		if IsDisabled() {
			// feature is experimental and disabled
			return fmt.Errorf("%s is experimental and not enabled by default. To use, please set %s=%s environment variable", feature, envName, enabled)
		}
		return warn()
	}
	return nil
}

// warn prints a warning message for using the experimental feature.
func warn() error {
	_, err := fmt.Fprintf(os.Stderr, "Warning: This feature is experimental and may not be fully tested or completed and may be deprecated. Report any issues to \"https://github/notaryproject/notation\"\n")
	return err
}

// HideFlags hide experimental flags when NOTATION_EXPERIMENTAL is disabled.
func HideFlags(cmd *cobra.Command, flags ...string) {
	if IsDisabled() {
		flagsSet := cmd.Flags()
		for _, flag := range flags {
			flagsSet.MarkHidden(flag)
		}
	}
}
