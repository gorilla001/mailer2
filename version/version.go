package version

var (
	version   string
	gitCommit string
)

// GetVersion is exported
func GetVersion() string {
	return version
}

// GetGitCommit is exported
func GetGitCommit() string {
	return gitCommit
}
