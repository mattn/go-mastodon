module github.com/mattn/go-mastodon

go 1.23

require (
	github.com/gorilla/websocket v1.5.3
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80
)

retract [v0.0.7+incompatible, v0.0.7+incompatible] // Accidental; no major changes or features.
