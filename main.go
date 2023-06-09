package main

import (
	"discordBot/setu"
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
    commands = []*discordgo.ApplicationCommand{&setu.SetuCmd}
    dg *discordgo.Session
    err error
	commandHandlers = setu.CommandHandler
)
func init() {
	token := "MTA2ODEwODMzMzU1MjI0MjgwMA.GvuAGW.SrlAAHj_3BlEO46CVyISouFk_cXgRT0iKUWz-Y"

	// creates a new Discord session
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println("创建机器人错误,", err)
		return
	}
}

func init() {
	dg.AddHandler(func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(dg, i)
		}
	})
}

func main() {
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// 只監聽訊息
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 开放通讯
	err = dg.Open()
	if err != nil {
		log.Println("开放通讯错误,", err)
		return
	}

	//注册命令
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if pong, _ := regexp.MatchString("[pong|PONG]", m.Content); pong {
		embed, err := setu.GetSetu("1","")
		if err != nil {
			s.ChannelMessageCrosspost(m.ChannelID, "获取色图失败："+err.Error())
			return
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

