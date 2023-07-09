package utility

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/gen2brain/beeep"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Notify(notifyGamesName []string) {
	if len(notifyGamesName) == 0 {
		return
	}

	title := "New games cheaper"
	message := "The following games are cheaper than before:\n"
	for _, game := range notifyGamesName {
		message += game + "\n"
	}

	beeep.Notify(title, message, "assets/information.png")
}

func CheckForDataFolder() {
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.Mkdir("./data", 0755)
		Check(err)
	}
}

func OpenBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
