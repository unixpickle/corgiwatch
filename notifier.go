package main

import (
	"net/http"
	"time"

	"github.com/unixpickle/fbmsgr"
)

const ExpireTime = time.Hour * 24 * 7

type Notifier struct {
	User      string
	Pass      string
	Recipient string

	session    *fbmsgr.Session
	expiration time.Time
}

func (n *Notifier) Notify(puppy *Puppy) error {
	session, err := n.GetSession()
	if err != nil {
		return err
	}
	if _, err := session.SendText(n.Recipient, n.PuppyMessage(puppy)); err != nil {
		return err
	}

	if puppy.PhotoURL == "" {
		return nil
	}

	resp, err := http.Get(puppy.PhotoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	attachment, err := session.Upload("puppy.jpg", resp.Body)
	if err != nil {
		return err
	}
	_, err = session.SendAttachment(n.Recipient, attachment)
	return err
}

func (n *Notifier) PuppyMessage(puppy *Puppy) string {
	var possessive string
	var pronoun string
	if puppy.Gender == "Male" {
		possessive = "His"
		pronoun = "He"
	} else {
		possessive = "Her"
		pronoun = "She"
	}
	if puppy.Registration == "" {
		puppy.Registration = "not"
	}
	return "New Puppy! " + possessive + " name is " + puppy.Name + ". " +
		pronoun + " is " + puppy.Registration + " registered. " +
		pronoun + " costs " + puppy.Price + ". " +
		pronoun + " is " + puppy.Age + " and is available on " + puppy.AvailableDate + ". " +
		"See more at " + puppy.ListingURL
}

func (n *Notifier) GetSession() (*fbmsgr.Session, error) {
	if n.session != nil && time.Now().Before(n.expiration) {
		return n.session, nil
	}
	var err error
	n.session, err = fbmsgr.Auth(n.User, n.Pass)
	if err != nil {
		return nil, err
	}
	n.expiration = time.Now().Add(ExpireTime)
	return n.session, nil
}
