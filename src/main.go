package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"time"
)

var (
	messages []string
)

func main() {
	fmt.Printf("Upload time: %s\n", os.Getenv("UPLOAD_TIME"))

	session, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.AddHandler(readyHandler)
	session.AddHandler(interactionHandler)

	if err := session.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go scheduleMessages(session)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	select {}
}

func readyHandler(s *discordgo.Session, _ *discordgo.Ready) {
	// Register slash commands
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "add",
		Description: "Add a message to the anonymous queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "message",
				Description: "The message to add",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
}

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "add":
		message := i.ApplicationCommandData().Options[0].StringValue()
		messages = append(messages, message)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Ihre (hoffentlich) erregende Bewegtbildproduktion wurde in die Schleife übernommen!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			return
		}
	}
}

func scheduleMessages(s *discordgo.Session) {
	for {
		now := time.Now()
		if now.Format("15:04") == os.Getenv("UPLOAD_TIME") {
			// Send accumulated messages
			if len(messages) > 0 {
				channelID := os.Getenv("CHANNEL_ID") // Replace with the ID of the channel to send messages
				for _, msg := range messages {
					_, err := s.ChannelMessageSend(channelID, msg)
					if err != nil {
						return
					}
				}
				messages = nil // Clear the message queue
			} else {
				log.Println("No messages to send.")
			}
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}
