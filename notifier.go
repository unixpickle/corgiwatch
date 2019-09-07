package main

import (
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
	session, err := n.getSession()
	if err != nil {
		return err
	}
	_, err = session.SendText(n.Recipient, n.PuppyMessage(puppy))
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
		"See more at " + puppy.ListingURL
}

func (n *Notifier) getSession() (*fbmsgr.Session, error) {
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
