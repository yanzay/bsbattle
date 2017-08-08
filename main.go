package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

type Buildings struct {
	Barracks  int
	Wall      int
	Trebuchet int
	Storage   int
	Houses    int
}

var buildStore *BuildStore

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN should not be empty")
	}
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatalf("can't create server: %q", err)
	}
	buildStore = NewBuildStore("builds.db")
	bot.Handle("/start", "Welcome to BS Battle Advice!\nForward your Buildings and Workshop to get advice.")
	bot.HandleDefault(parserHanlder)
	bot.ListenAndServe()
}

func parserHanlder(m *tbot.Message) {
	log.Infof("%s - %s", m.From, m.Text())
	log.Info([]byte(m.Text()))
	oldBuilds := buildStore.GetBuildings(m.From)
	builds, err := parseBuildings(m.Text())
	if err != nil {
		log.Errorf("can't parse buildings for %s: %q", m.Text(), err)
	}
	builds = mergeBuildings(oldBuilds, builds)
	buildStore.SaveBuildings(m.From, builds)
	m.Replyf("Barracks: %d\nWall: %d\nTrebuchet: %d\nStorage: %d\nHouses: %d", builds.Barracks, builds.Wall, builds.Trebuchet, builds.Storage, builds.Houses)
	if builds.Barracks == 0 || builds.Wall == 0 || builds.Trebuchet == 0 || builds.Storage == 0 || builds.Houses == 0 {
		m.Reply("Not enough information to advice.")
	} else {
		m.Reply(recommend(builds))
	}
}

func parseBuildings(text string) (*Buildings, error) {
	builds := &Buildings{}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		log.Info(line)
		log.Info([]byte(line))
		if builds.Barracks == 0 {
			fmt.Sscanf(line, "🛡   %d", &builds.Barracks)
		}
		if builds.Wall == 0 {
			fmt.Sscanf(line, "🏰   %d", &builds.Wall)
		}
		if builds.Storage == 0 {
			fmt.Sscanf(line, "🏚   %d", &builds.Storage)
		}
		if builds.Houses == 0 {
			fmt.Sscanf(line, "🏘   %d", &builds.Houses)
		}
		if builds.Trebuchet == 0 {
			str := string([]byte{226, 154, 148, 84, 114, 101, 98, 117, 99, 104, 101, 116})
			fmt.Sscanf(line, str+"%d", &builds.Trebuchet)
		}
	}

	return builds, nil
}

func mergeBuildings(oldBuilds, newBuilds *Buildings) *Buildings {
	if oldBuilds.Barracks > newBuilds.Barracks {
		newBuilds.Barracks = oldBuilds.Barracks
	}
	if oldBuilds.Wall > newBuilds.Wall {
		newBuilds.Wall = oldBuilds.Wall
	}
	if oldBuilds.Trebuchet > newBuilds.Trebuchet {
		newBuilds.Trebuchet = oldBuilds.Trebuchet
	}
	if oldBuilds.Storage > newBuilds.Storage {
		newBuilds.Storage = oldBuilds.Storage
	}
	if oldBuilds.Houses > newBuilds.Houses {
		newBuilds.Houses = oldBuilds.Houses
	}
	return newBuilds
}
