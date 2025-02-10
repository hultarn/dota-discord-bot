#!/bin/sh

tmux new-session -d -s dota-league-service "go run main.go discord start-league-bot --env='/path'"

tmux new-session -d -s dota-signup-service "go run main.go discord start-signup-bot --env='/path'"

