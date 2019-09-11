package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Puppy struct {
	Name         string
	FullName     string
	Gender       string
	Price        string
	Registration string
	ListingURL   string
	PhotoURL     string

	// Fields on the puppy details page.
	AvailableDate string
	Age           string
}

func FetchPuppies(puppyURL string) ([]*Puppy, error) {
	parsedURL, err := url.Parse(puppyURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(puppyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	searchResults, ok := scrape.Find(parsed, scrape.ByClass("view-puppy-search"))
	if !ok {
		return nil, errors.New("no puppy search results")
	}
	contents, ok := scrape.Find(searchResults, scrape.ByClass("view-content"))
	if !ok {
		return nil, errors.New("no search result view content")
	}

	var puppies []*Puppy
	for _, row := range scrape.FindAll(contents, scrape.ByClass("views-row-inner")) {
		puppy := ParsePuppy(row, parsedURL)
		if puppy != nil {
			puppies = append(puppies, puppy)
		}
	}
	return puppies, nil
}

func ParsePuppy(row *html.Node, u *url.URL) *Puppy {
	res := &Puppy{}
	if img, ok := scrape.Find(row, scrape.ByTag(atom.Img)); ok {
		res.PhotoURL = scrape.Attr(img, "src")
	}
	if title, ok := scrape.Find(row, scrape.ByClass("views-field-title")); ok {
		if link, ok := scrape.Find(title, scrape.ByTag(atom.A)); ok {
			uCopy := *u
			uCopy.Path = scrape.Attr(link, "href")
			uCopy.RawQuery = ""
			res.ListingURL = uCopy.String()
		}
		res.FullName = strings.TrimSpace(scrape.Text(title))
		res.Name = strings.TrimSpace(strings.Split(res.FullName, "-")[0])
	}
	getField := func(name string) string {
		if elem, ok := scrape.Find(row, scrape.ByClass("views-field-field-"+name)); ok {
			if content, ok := scrape.Find(elem, scrape.ByClass("field-content")); ok {
				return strings.TrimSpace(scrape.Text(content))
			}
		}
		return ""
	}
	res.Gender = getField("sex")
	res.Price = getField("asking-price")
	res.Registration = getField("registration")
	return res
}

func (p *Puppy) FetchDetails() error {
	resp, err := http.Get(p.ListingURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	availableDate, ok := scrape.Find(parsed, scrape.ByClass("field-name-field-date-available"))
	if !ok {
		return errors.New("could not find availability date")
	}
	p.AvailableDate = FieldItemValue(availableDate)

	age, ok := scrape.Find(parsed, scrape.ByClass("age-in-weeks"))
	if !ok {
		return errors.New("could not find age in weeks")
	}
	p.Age = strings.TrimSpace(scrape.Text(age))

	return nil
}

func FieldItemValue(item *html.Node) string {
	items, ok := scrape.Find(item, scrape.ByClass("field-items"))
	if !ok {
		return ""
	}
	return strings.TrimSpace(scrape.Text(items))
}
