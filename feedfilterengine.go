/*
 * package feedfilterengine
 *
 * filter rss feeds by given string and create a new one
 *
 * (c) 2016 by Björn Winkler
 *
 */

package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Rss2 struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XmlnsContent string   `xml:"xmlns content,attr"`
	XmlnsAtom    string   `xml:"xmlns atom,attr"`

	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	// Optional
	PubDate string     `xml:"channel>pubDate,omitempty"`
	Items   []Rss2Item `xml:"channel>item,omitempty"`
}

type Rss2Item struct {
	// Required
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`

	// Optional
	Content  template.HTML `xml:"encoded,omitempty"`
	PubDate  string        `xml:"pubDate,omitempty"`
	Comments string        `xml:"comments,omitempty"`
}

func filterRss(url string) []byte {
	var err error

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	rss := Rss2{}
	err = xml.Unmarshal(body, &rss)

	newItems := []Rss2Item{}

	for _, item := range rss.Items {
		if strings.Contains(item.Title, "Freistetters Formelwelt") {
			newItems = append(newItems, item)
		}
	}

	rss.Items = newItems

	marshalled, err := xml.Marshal(rss)

	if err != nil {
		log.Fatal(err)
	}

	return append([]byte(xml.Header), marshalled...)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func rss(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", filterRss("http://www.spektrum.de/alias/rss/spektrum-de-rss-feed/996406"))
	//fmt.Fprintf(w, "%s", readRss("http://www.zdnet.de/feed/"))
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/rss", rss)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
