package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-mastodon"
	"github.com/mattn/go-tty"
	"github.com/urfave/cli"
	"golang.org/x/net/html"
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
		return s
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
	readUsername = func() (string, error) {
		b, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	readPassword func() (string, error)
)

func prompt() (string, string, error) {
	fmt.Print("E-Mail: ")
	email, err := readUsername()
	if err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	var password string
	if readPassword == nil {
		t, err := tty.Open()
		if err != nil {
			return "", "", err
		}
		defer t.Close()
		password, err = t.ReadPassword()
	} else {
		password, err = readPassword()
	}
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

func getConfig(c *cli.Context) (string, *mastodon.Config, error) {
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
	var file string
	profile := c.String("profile")
	if profile != "" {
		file = filepath.Join(dir, "settings"+profile+".json")
	} else {
		file = filepath.Join(dir, "settings.json")
	}
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

func authenticate(client *mastodon.Client, config *mastodon.Config, file string) error {
	email, password, err := prompt()
	if err != nil {
		return err
	}
	err = client.Authenticate(context.Background(), email, password)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to store file: %v", err)
	}
	err = ioutil.WriteFile(file, b, 0700)
	if err != nil {
		return fmt.Errorf("failed to store file: %v", err)
	}
	return nil
}

func argstr(c *cli.Context) string {
	a := []string{}
	for i := 0; i < c.NArg(); i++ {
		a = append(a, c.Args().Get(i))
	}
	return strings.Join(a, " ")
}

func fatalIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	os.Exit(1)
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "mstdn"
	app.Usage = "mastodon client"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile",
			Usage: "profile name",
			Value: "",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "toot",
			Usage: "post toot",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ff",
					Usage: "post utf-8 string from a file(\"-\" means STDIN)",
					Value: "",
				},
				cli.IntFlag{
					Name:  "i",
					Usage: "in-reply-to",
					Value: 0,
				},
			},
			Action: cmdToot,
		},
		{
			Name:  "stream",
			Usage: "stream statuses",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "output JSON",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "straem type (public,public/local,user:NAME,hashtag:TAG)",
				},
				cli.BoolFlag{
					Name:  "simplejson",
					Usage: "output simple JSON",
				},
			},
			Action: cmdStream,
		},
		{
			Name:   "timeline",
			Usage:  "show timeline",
			Action: cmdTimeline,
		},
		{
			Name:   "notification",
			Usage:  "show notification",
			Action: cmdNotification,
		},
		{
			Name:   "instance",
			Usage:  "show instance information",
			Action: cmdInstance,
		},
		{
			Name:   "account",
			Usage:  "show account information",
			Action: cmdAccount,
		},
		{
			Name:   "search",
			Usage:  "search content",
			Action: cmdSearch,
		},
		{
			Name:   "follow",
			Usage:  "follow account",
			Action: cmdFollow,
		},
		{
			Name:   "followers",
			Usage:  "show followers",
			Action: cmdFollowers,
		},
		{
			Name:   "upload",
			Usage:  "upload file",
			Action: cmdUpload,
		},
	}
	app.Setup()
	return app
}

func acct(acct, host string) string {
	if !strings.Contains(acct, "@") {
		acct += "@" + host
	}
	return acct
}

type screen struct {
	host string
}

func newScreen(config *mastodon.Config) *screen {
	var host string
	u, err := url.Parse(config.Server)
	if err == nil {
		host = u.Host
	}
	return &screen{host}
}

func (s *screen) displayError(w io.Writer, e error) {
	color.Set(color.FgYellow)
	fmt.Fprintln(w, e.Error())
	color.Set(color.Reset)
}

func (s *screen) displayStatus(w io.Writer, t *mastodon.Status) {
	if t.Reblog != nil {
		color.Set(color.FgHiRed)
		fmt.Fprint(w, acct(t.Account.Acct, s.host))
		color.Set(color.Reset)
		fmt.Fprint(w, " reblogged ")
		color.Set(color.FgHiBlue)
		fmt.Fprintln(w, acct(t.Reblog.Account.Acct, s.host))
		fmt.Fprintln(w, textContent(t.Reblog.Content))
		color.Set(color.Reset)
	} else {
		color.Set(color.FgHiRed)
		fmt.Fprintln(w, acct(t.Account.Acct, s.host))
		color.Set(color.Reset)
		fmt.Fprintln(w, textContent(t.Content))
	}
}

func run() int {
	app := makeApp()

	app.Before = func(c *cli.Context) error {
		file, config, err := getConfig(c)
		if err != nil {
			return err
		}

		client := mastodon.NewClient(config)
		if config.AccessToken == "" {
			return authenticate(client, config, file)
		}
		app.Metadata = map[string]interface{}{
			"client": client,
			"config": config,
		}
		return nil
	}

	app.Run(os.Args)
	return 0
}

func main() {
	os.Exit(run())
}
