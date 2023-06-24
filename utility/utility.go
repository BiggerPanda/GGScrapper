package utility

import "github.com/gen2brain/beeep"

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
