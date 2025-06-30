package main

import (
	"os/exec"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/designsbysm/timber/v2"
	"github.com/spf13/viper"
)

func pull(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("git", "pull", "--prune", "--tags")
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() > 1 {
				timber.Error(path, err)
			}
		} else {
			timber.Error(path, err)
		}
	}

	install := viper.GetBool("install")

	if install {
		cmd := exec.Command("npm", "isntall")
		cmd.Dir = path

		if err := cmd.Run(); err != nil {
			timber.Error(path, err)
		}
	}
}

func remotesPull(statuses []Status) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	var wg sync.WaitGroup

	for _, status := range statuses {
		if !status.Modified && status.Remote == RemoteAhead {
			wg.Add(1)
			go pull(status.Path, &wg)
		}
	}

	wg.Wait()
	s.Stop()

	return nil
}
