<p align="center">
  <a href="https://github.com/sokolovskiyma/duckspeak"><img src="https://goreportcard.com/badge/github.com/sokolovskiyma/duckspeak" alt="Go Report Card"></a>
</p>

# duckspeak
CLI generator of random sentences based on [Markov chains](https://en.wikipedia.org/wiki/Markov_chain)

# Installation
## Linux/Mac
Clone this repo
```
go build duckspeak.go
./guckspeak -l texts/ -g 5 -c
```

# Usage
- `-l /path/` Adds words to the dictionary from one file or from the entire directory
- `-g n`      Generate n random sentences
- `-c`        Clear dictionary
- `-h`        Help
