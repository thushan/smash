package smash

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// EnvironmentInfo contains information about the execution environment
type EnvironmentInfo struct {
	TermType      string
	IsTerminal    bool
	IsCI          bool
	IsTesting     bool
	IsContainer   bool
	ForceColor    bool
	NoColor       bool
	ShouldAnimate bool
}

// DetectEnvironment analyzes the current execution environment
func DetectEnvironment() EnvironmentInfo {
	env := EnvironmentInfo{
		IsTerminal: term.IsTerminal(int(os.Stdout.Fd())) &&
			term.IsTerminal(int(os.Stderr.Fd())),
		TermType: os.Getenv("TERM"),
	}

	// Check for CI environments
	ciVars := []string{
		"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS",
		"GITLAB_CI", "JENKINS_URL", "BUILDKITE", "DRONE",
		"TRAVIS", "CIRCLECI", "APPVEYOR", "CODEBUILD_BUILD_ID",
		"TEAMCITY_VERSION", "TF_BUILD", "BUDDY", "WERCKER",
	}
	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			env.IsCI = true
			break
		}
	}

	// Check for test environment
	testVars := []string{"GO_TEST", "TEST", "TESTING"}
	for _, v := range testVars {
		if os.Getenv(v) != "" {
			env.IsTesting = true
			break
		}
	}

	// Check if running in container
	env.IsContainer = isRunningInContainer()

	// Check color preferences
	env.NoColor = os.Getenv("NO_COLOR") != "" || env.TermType == "dumb"
	env.ForceColor = os.Getenv("FORCE_COLOR") != "" ||
		os.Getenv("CLICOLOR_FORCE") == "1"

	// Determine if we should show animations (spinners)
	// Animations are disabled in CI, testing, non-TTY environments, or when NO_COLOR is set
	env.ShouldAnimate = env.IsTerminal &&
		!env.IsCI &&
		!env.IsTesting &&
		!env.NoColor &&
		env.TermType != "dumb"

	return env
}

// isRunningInContainer detects if the application is running inside a container
func isRunningInContainer() bool {
	// Check for Docker
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check for other container runtimes via environment
	if os.Getenv("container") != "" {
		return true
	}

	// Check cgroup for container signatures
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		cgroupContent := string(data)
		containerSignatures := []string{
			"docker", "kubepods", "containerd",
			"lxc", "machine.slice", "podman",
		}
		for _, sig := range containerSignatures {
			if strings.Contains(cgroupContent, sig) {
				return true
			}
		}
	}

	return false
}

// CICapabilities describes what a CI environment supports
type CICapabilities struct {
	SupportsColor     bool
	SupportsBasicANSI bool
	SupportsGrouping  bool
}

// GetCICapabilities returns the capabilities of the current CI environment
func GetCICapabilities() CICapabilities {
	caps := CICapabilities{}

	// GitHub Actions supports some ANSI codes and grouping
	if os.Getenv("GITHUB_ACTIONS") != "" {
		caps.SupportsColor = true
		caps.SupportsBasicANSI = true
		caps.SupportsGrouping = true // ::group:: syntax
	}

	// GitLab CI has good terminal emulation
	if os.Getenv("GITLAB_CI") != "" {
		caps.SupportsColor = true
		caps.SupportsBasicANSI = true
	}

	// Travis CI supports color
	if os.Getenv("TRAVIS") == "true" {
		caps.SupportsColor = true
		caps.SupportsBasicANSI = true
	}

	// Jenkins depends on plugins, be conservative
	if os.Getenv("JENKINS_URL") != "" {
		caps.SupportsColor = false
		caps.SupportsBasicANSI = false
	}

	return caps
}
