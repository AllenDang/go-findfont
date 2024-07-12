// +build dragonfly freebsd linux nacl netbsd openbsd solaris
//go:build dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build dragonfly freebsd linux nacl netbsd openbsd solaris

// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package findfont

import (
	"os"
	"path/filepath"
)

var unixSystemFontGlobs = []string{
	"/nix/store/*/share/fonts/truetype",
	"/nix/store/*/share/fonts/opentype",
}

var unixSystemFontPaths = []string{
	"/usr/share/fonts/",
	"/usr/share/fonts/truetype/",
	"/usr/share/fonts/opentype/",
}

func getFontDirectories() (paths []string) {
	directories := getUserFontDirs()
	directories = append(directories, getSystemFontDirs()...)
	return directories
}

func getUserFontDirs() (paths []string) {
	if dataPath := os.Getenv("XDG_DATA_HOME"); dataPath != "" {
		return []string{expandUser("~/.fonts/"), filepath.Join(expandUser(dataPath), "fonts")}
	}
	return []string{expandUser("~/.fonts/"), expandUser("~/.local/share/fonts/")}
}

func getSystemFontDirs() (paths []string) {
	if dataPaths := os.Getenv("XDG_DATA_DIRS"); dataPaths != "" {
		for _, dataPath := range filepath.SplitList(dataPaths) {
			dir := filepath.Join(expandUser(dataPath), "fonts")
			if _, err := os.Stat(dir); err == nil {
				paths = append(paths, dir)
			}
			dir = filepath.Join(dir, "truetype")
			if _, err := os.Stat(dir); err == nil {
				paths = append(paths, dir)
			}
			dir = filepath.Join(dir, "opentype")
			if _, err := os.Stat(dir); err == nil {
				paths = append(paths, dir)
			}
		}
	}

	paths = append(paths, "/usr/local/share/fonts/", "/usr/share/fonts/")

	if _, err := os.Stat("/nix/store"); err == nil {
		for _, glob := range unixSystemFontGlobs {
			res, _ := filepath.Glob(glob)
			if len(res) > 0 {
				paths = append(paths, res...)
			}
		}
	}

	return paths
}
