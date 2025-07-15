package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func gatherFolders(src string, maxDepth int) (files []string, err error) {
	noRoot := viper.GetBool("noRoot")
	rootDepth := len(strings.Split(src, string(os.PathSeparator)))
	if strings.HasPrefix(src, "./") {
		rootDepth -= 1
	}

	err = filepath.WalkDir(src, func(path string, dir fs.DirEntry, err error) error {
		currentDepth := strings.Count(path, string(os.PathSeparator))

		if err != nil {
			return err
		} else if dir.IsDir() && currentDepth > rootDepth+maxDepth || strings.Contains(path, "node_modules") {
			return fs.SkipDir
		} else if dir.IsDir() && strings.Contains(path, ".git") {
			if noRoot && filepath.Dir(path) == src {
				return fs.SkipDir
			}

			files = append(files, filepath.Dir(path))

			return fs.SkipDir
		} else if dir.IsDir() {
			gitDir := filepath.Join(path, ".git")

			if info, err := os.Stat(gitDir); err == nil && !info.IsDir() {
				if noRoot && path == src {
					return fs.SkipDir
				}

				files = append(files, path)

				return fs.SkipDir
			}
		}

		return nil
	})

	return
}
