package main

import (
	"fmt"

	"github.com/designsbysm/timber/v2"
	"github.com/spf13/viper"
)

func pathWidth(statuses []Status) int {
	width := 0

	for _, status := range statuses {
		if len(status.Path) > width {
			width = len(status.Path)
		}
	}

	return width
}

func outputStatuses(statuses []Status) {
	changes := viper.GetBool("changes")

	filtered := []Status{}
	for _, status := range statuses {
		if !changes {
			filtered = append(filtered, status)
		} else if changes && (status.Modified || (status.Remote != NoRemote && status.Remote != InSync)) {
			filtered = append(filtered, status)
		}
	}

	maxPathWidth := pathWidth(filtered)

	if len(filtered) == 0 {
		timber.Info("All repos are clean and up to date.")
		return
	}

	for _, status := range filtered {
		modifyCode := " "
		remoteCode := " "
		color := colorWhite

		if status.Modified {
			modifyCode = "*"
		}

		if status.Remote == InSync {
			color = colorGreen
		} else if status.Remote == LocalAhead {
			remoteCode = "↑"
			color = colorPurple
		} else if status.Remote == RemoteAhead {
			remoteCode = "↓"
			color = colorYellow
		} else if status.Remote == Diverged {
			remoteCode = "*"
			color = colorRed
		} else if status.Remote == Gone {
			remoteCode = "∅"
			color = colorGray
		}

		fmt.Printf("%-*s   %s[%s%s]   %s%s\n", maxPathWidth, status.Path, color, modifyCode, remoteCode, status.Branch, colorReset)
	}

	timber.Info("")
}
