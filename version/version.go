package version

import "runtime"

var (
	// Package package name
	Package = "github.com/dumbogo/chess"

	// Version holds the complete version number. Filled in at linking time.
	Version = "beta"

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	// GoVersion is Go tree's version.
	GoVersion = runtime.Version()
)
