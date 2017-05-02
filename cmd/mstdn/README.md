# mstdn

command line tool for mstdn.jp

## Usage

```
NAME:
   mstdn - mastodon client

USAGE:
   mstdn [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     toot          post toot
     stream        stream statuses
     timeline      show timeline
     notification  show notification
     instance      show instance information
     account       show account information
     search        search content
     follow        follow account
     followers     show followers
     upload        upload file
     delete        delete status
     init          initialize profile
     mikami        search mikami
     xsearch       cross search
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile value  profile name
   --help, -h       show help
   --version, -v    print the version
```

## Installation

```
$ go get github.com/mattn/go-mastodon/cmd/mstdn
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
