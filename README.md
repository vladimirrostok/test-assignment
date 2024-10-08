# Summary
It's a brief assessment, that's a command-line version of the web game Wordle in Go.

This repository contains the ["Programming Skill Assessment" definition file](wordletest_go_v1.odt) and the [Go project code itself](/app/Go/).

# Running
Go to the [/app/Go/wordle/cmd](/app/Go/wordle/cmd/) directory and execute `go run wordle.go`  

To run the entire Go codebase tests go to [/app/Go/](/app/Go/) directory and execute `go test ./...`

![Game screenshot](assets/Game-Screenshot.png) 

# Running in Docker Compose
Go to the [/app/Go/](/app/Go/) and execute the `docker-compose up --build`

## Build the Wordle service Dockerfile manually
Go to the [/app/Go/](/app/Go/) and execute the `docker build -f wordle/Dockerfile . -t wordle-image`

NB! `go.sum` file is located at the `/Go/` directory while the Dockerfile is located at the `/Go/wordle/` directory, so we execute `build` command from the directory with `go.sum` file directory which is a `/Go`/ directory and point to Dockerfile with `-f` flag, executing from local (go.sum file) context `.` and naming the container with `-t` flag.


## Dockerfile build 
![Dockerfile build screenshot](assets/Dockerfile-Build.png) 

## Docker Compose build 
![Docker Compose build screenshot](assets/Docker-Compose-Build.png) 

## Tech Stack
* Golang 1.22.4 standard library + github.com/fatih/color library                 
    * Docker
    * Docker Compose   

## Requirements
* The only non-standard library allowed is the github.com/fatih/color for the terminal colors.
* All errors must be handled in main(); i.e. no log.Fatalln or os.Exit in other functions.
* Package math/rand can be used instead of crypto/rand.     

# Methods & tweaks
## [CMD pattern](https://github.com/golang-standards/project-layout/blob/master/cmd/README.md)
CMD pattern - a file convention in Go, helps to manage multiple main.go entry-points in the future and reuse the code, this also helps to keep the root directory clean [e.g., Kubernetes uses this pattern a lot](https://github.com/kubernetes/kubernetes/tree/master/cmd). 

## DDD - Domain Driven Design
Domain Driven Design was the key idea behind thisc codebase. Altough the task is tiny, it carries few domain errors related to system data format (only unicode characters allowed, consistant data input length) and a model file with little type-related validation functions.

## [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
Table Driven Tests is a way to write cleaner tests, this reduces the amount of repeated code in the tests and boosts code readability. 

