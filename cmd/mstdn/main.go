package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/mattn/go-tty"
	"golang.org/x/net/html"
)

var (
	toot     = flag.String("t", "", "toot text")
	stream   = flag.Bool("S", false, "streaming public")
	fromfile = flag.String("ff", "", "post utf-8 string from a file(\"-\" means STDIN)")
)

func readFile(filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

func textContent(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer

	var extractText func(node *html.Node, w *bytes.Buffer)
	extractText = func(node *html.Node, w *bytes.Buffer) {
		if node.Type == html.TextNode {
			data := strings.Trim(node.Data, "\r\n")
			if data != "" {
				w.WriteString(data)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extractText(c, w)
		}
		if node.Type == html.ElementNode {
			name := strings.ToLower(node.Data)
			if name == "br" {
				w.WriteString("\n")
			}
		}
	}
	extractText(doc, &buf)
	return buf.String()
}

var (
	readUsername func() (string, error) = func() (string, error) {
		b, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	readPassword func() (string, error)
)

func prompt() (string, string, error) {
	t, err := tty.Open()
	if err != nil {
		return "", "", err
	}
	defer t.Close()

	fmt.Print("E-Mail: ")
	email, err := readUsername()
	if err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	var password string
	if readPassword == nil {
		password, err = t.ReadPassword()
	} else {
		password, err = readPassword()
	}
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

func getConfig() (string, *mastodon.Config, error) {
	dir := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
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
		ClientID:     "171d45f22068a5dddbd927b9d966f5b97971ed1d3256b03d489f5b3a83cdba59",
		ClientSecret: "574a2cf4b3f28a5fa0cfd285fc80cfe9daa419945163ef18f5f3d0022f4add28",
	}
	if err == nil {
		err = json.Unmarshal(b, &config)
		if err != nil {
			return "", nil, fmt.Errorf("could not unmarshal %v: %v", file, err)
		}
	}
	return file, config, nil
}

func authenticate(client *mastodon.Client, config *mastodon.Config, file string) {
	email, password, err := prompt()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Authenticate(email, password)
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
}

func streaming(client *mastodon.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	q, err := client.StreamingPublic(ctx)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-sc
		cancel()
		close(q)
	}()
	for e := range q {
		switch t := e.(type) {
		case *mastodon.UpdateEvent:
			color.Set(color.FgHiRed)
			fmt.Println(t.Status.Account.Username)
			color.Set(color.Reset)
			fmt.Println(textContent(t.Status.Content))
		case *mastodon.ErrorEvent:
			color.Set(color.FgYellow)
			fmt.Println(t.Error())
			color.Set(color.Reset)
		}
	}
}

func init() {
	flag.Parse()
	if *fromfile != "" {
		text, err := readFile(*fromfile)
		if err != nil {
			log.Fatal(err)
		}
		*toot = string(text)
	}
}

func post(client *mastodon.Client, text string) {
	_, err := client.PostStatus(&mastodon.Toot{
		Status: text,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func timeline(client *mastodon.Client) {
	timeline, err := client.GetTimelineHome()
	if err != nil {
		log.Fatal(err)
	}
	for i := len(timeline) - 1; i >= 0; i-- {
		t := timeline[i]
		color.Set(color.FgHiRed)
		fmt.Println(t.Account.Username)
		color.Set(color.Reset)
		fmt.Println(textContent(t.Content))
	}
}

func main() {
	file, config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := mastodon.NewClient(config)

	if config.AccessToken == "" {
		authenticate(client, config, file)
		return
	}

	if *toot != "" {
		post(client, *toot)
		return
	}

	if *stream {
		streaming(client)
		return
	}

	timeline(client)
}
