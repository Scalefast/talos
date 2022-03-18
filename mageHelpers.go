//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func flags() (string, error) {
	// Make version increase automatically
	return fmt.Sprintf(`-X "main.version=%s" -X "main.commit=%s" -X "main.date=%s" -X "main.builtBy=%s"`, version, commit, date, builtBy), nil
}

// Implement automatic semver management
func getVersion() string {
	cmd, ok := sh.Output("git", "describe", "--tags")
	if ok != nil {
		return "-tag"
	}
	return cmd
}
func getCommit() string {
	cmd, ok := sh.Output("git", "describe", "--always", "--long", "--dirty")
	if ok != nil {
		return ""
	}
	return cmd
}
func getDate() string {
	cmd, ok := sh.Output("date")
	if ok != nil {
		return ""
	}
	return cmd
}
func getBuiltBy() string {
	cmd, ok := sh.Output("git", "--no-pager", "show", "-s", "--format='%ae'")
	if ok != nil {
		return ok.Error()
	}
	// change str from 'xxx' to xxx
	if cmd != "" {
		sl := len(cmd)
		cmd = cmd[1 : sl-1]
	}
	return cmd
}
