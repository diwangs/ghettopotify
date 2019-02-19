# Ghettopotify
A stupid simple music streaming program written in Go

Uses UDP. Library used: golang mp3 player & decoder by hajimehoshi (github.com/hajimehoshi)

Compile:
- `export GOPATH=$PWD`
- `go install client server `
- run the programs at `./bin`

Running:
- By default, server will listen at port 6969
- By default, 127.0.0.1:6969 will be added to the client subscription list
- run server, lalu run client

Commands:
- play <filename> (pause and resume with "pause" and "resume" respectively)
- exit
- ls: list songs
- lschan: list channels
- chchan: change channels
- sub: add a server to subscription list

