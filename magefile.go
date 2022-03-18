//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// To change talos build version, we have to change this
	version string = getVersion()
	commit  string = getCommit()
	date    string = getDate()
	builtBy string = getBuiltBy()
)

// Runs go mod download
func DownloadDeps() error {
	return sh.Run("go", "mod", "download")
}

// Runs go mod download and then builds the binary.
func Build() error {
	mg.Deps(DownloadDeps)

	ldf, err := flags()
	if err != nil {
		return err
	}
	return sh.Run("go", "build", "--ldflags="+ldf)
}

func Install() error {
	mg.Deps(DownloadDeps)

	ldf, err := flags()
	if err != nil {
		return err
	}
	return sh.Run("go", "install", "--ldflags="+ldf)
}
