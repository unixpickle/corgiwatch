package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/unixpickle/essentials"
)

const PuppyURL = "https://www.lancasterpuppies.com/puppy-search/breed/welsh%20corgi%20%28pembroke%29?sort=created&order=desc"

func main() {
	var puppyURL string
	var savePath string
	var messengerUser string
	var messengerPass string
	var messengerRecipient string
	var registrationFilter string

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Required flags: -user, -pass, -recipient")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}

	flag.StringVar(&puppyURL, "url", PuppyURL, "URL to track")
	flag.StringVar(&savePath, "save", "puppies.json", "path to save previously seen puppies")
	flag.StringVar(&messengerUser, "user", "", "Facebook messenger username")
	flag.StringVar(&messengerPass, "pass", "", "Facebook messenger password")
	flag.StringVar(&messengerRecipient, "recipient", "", "Facebook messenger recipient ID")
	flag.StringVar(&registrationFilter, "registration", "", "acceptable value for registration")

	flag.Parse()

	if messengerUser == "" || messengerPass == "" || messengerRecipient == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Creating feed...")
	feed, err := NewFeed(puppyURL, savePath)
	essentials.Must(err)

	log.Println("Creating notifier...")
	notifier := &Notifier{
		User:      messengerUser,
		Pass:      messengerPass,
		Recipient: messengerRecipient,
	}
	if _, err := notifier.GetSession(); err != nil {
		essentials.Die(err)
	}

	log.Println("Starting feed loop...")
	for {
		puppies, err := feed.Pull()
		if err != nil {
			log.Println("Pull error:", err)
		} else {
			for _, p := range puppies {
				if registrationFilter != "" && p.Registration != registrationFilter {
					log.Println("Skipping", p.Registration, "registered puppy.")
					continue
				}
				log.Println(notifier.PuppyMessage(p))
				if err := notifier.Notify(p); err != nil {
					log.Println("Send error:", err)
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}
