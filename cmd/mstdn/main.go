package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/mattn/go-tty"
	"golang.org/x/net/html"
)

var blockTags = []string{"div", "br", "p", "blockquote", "pre", "h1", "h2", "h3", "h4", "h5", "h6"}

func extractText(node *html.Node, w *bytes.Buffer) {
	if node.Type == html.TextNode {
		data := strings.Trim(node.Data, "\r\n")
		if data != "" {
			w.WriteString(data)
		}
	} else if node.Type == html.ElementNode {
		if node.Data == "li" {
			w.WriteString("\n* ")
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, w)
	}
	if node.Type == html.ElementNode {
		for _, bt := range blockTags {
			if strings.ToLower(node.Data) == bt {
				w.WriteString("\n")
				break
			}
		}
	}
}

func prompt() (string, string, error) {
	t, err := tty.Open()
	if err != nil {
		return "", "", err
	}
	defer t.Close()

	fmt.Print("E-Mail: ")
	b, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		return "", "", err
	}
	email := string(b)

	fmt.Print("Password: ")
	password, err := t.ReadPassword()
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

func getConfig() (string, *mastodon.Config, error) {
	dir := os.Getenv("HOME")
	if dir == "" && runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "mstdn")
		}
		dir = filepath.Join(dir, "mstdn")
	} else {
		dir = filepath.Join(dir, ".config", "mstdn")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", nil, err
	}
	file := filepath.Join(dir, "settings.json")
	b, err := ioutil.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return "", nil, err
	}
	config := &mastodon.Config{
		Server:       "https://mstdn.jp",
		ClientID:     "654a15390204e70d74f8d9264526e017e26d323e20e3f983409c157115009862",
		ClientSecret: "17274242a0846ebadcdda77727666c9d475f1989b56ad9bd959021f62f92a84c",
	}
	if err == nil {
		err = json.Unmarshal(b, &config)
		if err != nil {
			return "", nil, fmt.Errorf("could not unmarshal %v: %v", file, err)
		}
	}
	return file, config, nil
}

func main() {
	file, config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	c := mastodon.NewClient(config)

	if config.AccessToken == "" {
		email, password, err := prompt()
		if err != nil {
			log.Fatal(err)
		}
		err = c.Authenticate(email, password)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		err = ioutil.WriteFile(file, b, 0700)
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		return
	}

	timeline, err := c.GetTimeline("/api/v1/timelines/home")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range timeline {
		doc, err := html.Parse(strings.NewReader(t.Content))
		if err != nil {
			log.Fatal(err)
		}
		var buf bytes.Buffer
		extractText(doc, &buf)
		fmt.Println(t.Account.Username)
		fmt.Println(buf.String())
	}
}
