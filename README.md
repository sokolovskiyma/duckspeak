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
- `-l /path/` Add fame from a file or from all folders in the dictionary
- `-g n`      Generate n random sentences
- `-c`        Clear dictionary
- `-h`        Help
