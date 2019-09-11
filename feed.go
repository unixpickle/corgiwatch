package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Feed struct {
	PuppyURL string
	SavePath string

	history []*Puppy
}

func NewFeed(puppyURL, savePath string) (*Feed, error) {
	res := &Feed{PuppyURL: puppyURL, SavePath: savePath}
	if data, err := ioutil.ReadFile(savePath); os.IsNotExist(err) {
		res.history, err = FetchPuppies(puppyURL)
		if err != nil {
			return nil, err
		}
		if err := res.saveHistory(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &res.history); err != nil {
		return nil, err
	}
	return res, nil
}

func (f *Feed) Pull() ([]*Puppy, error) {
	puppies, err := FetchPuppies(f.PuppyURL)
	if err != nil {
		return nil, err
	}
	var newPuppies []*Puppy
	for _, p := range puppies {
		found := false
		for _, p1 := range f.history {
			if p.FullName == p1.FullName && p.ListingURL == p1.ListingURL {
				found = true
				break
			}
		}
		if !found {
			if err := p.FetchDetails(); err != nil {
				return nil, err
			}
			newPuppies = append(newPuppies, p)
		}
	}
	f.history = append(f.history, newPuppies...)
	if err := f.saveHistory(); err != nil {
		return nil, err
	}
	return newPuppies, nil
}

func (f *Feed) saveHistory() error {
	data, _ := json.Marshal(f.history)
	return ioutil.WriteFile(f.SavePath, data, 0755)
}
