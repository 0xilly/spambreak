package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Token     string   `yaml:"token"`
	Playing   []string `yaml:"playing"`
	SpamMatch []string `yaml:"spam"`
}

func main() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	initBot(getConf())
}

func getConf() Config {

	var conf Config
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		conf := Config{Token: "1234",
			Playing: []string{"A Game", " With Media's heart", "The Game... darn I lost", "Jeff 2 Electric Boogaloo"},
			SpamMatch: []string{"With our IRC ad service you can reach a global audience of entrepreneurs and fentanyl addicts with extraordinary engagement rates! https://williampitcock.com/",
				"I thought you guys might be interested in this blog by freenode staff member Bryan 'kloeri' Ostergaard",
				"Read what IRC investigative journalists have uncovered on the freenode pedophilia scandal"},
		}

		marshel, _ := yaml.Marshal(conf)

		err = ioutil.WriteFile("config.yaml", marshel, 0644)
	}

	file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	yaml.Unmarshal(file, &conf)

	return conf
}

func initBot(c Config) {
	dg, err := discordgo.New("Bot " + c.Token)

	if err != nil {
		log.Fatal("Unable to start discord session")
		os.Exit(2)
	}

	dg.AddHandler(msgDelete)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func msgDelete(s *discordgo.Session, m *discordgo.MessageCreate) {

	s.UpdateStatus(0, "Jeff2 Electric Boogaloo")
	if hasSpam(m.Content) {
		log.Printf("Deleted message from user {%s} in channel {%s} for spamming {%s}", m.Author.Username, m.ChannelID, m.Content)
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)

		println(string(m.ChannelID))
		if err != nil {
			fmt.Println(err)
		}

	}

}

func hasSpam(m string) bool {

	var spamlist Config = getConf()
	for _, list := range spamlist.SpamMatch {
		if strings.Contains(m, list) {
			return true
		}

	}
	return false
}
