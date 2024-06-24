# Dota Discord Bot

## Introduction
This project is part of a weekly in-house Dota2 tournament initiative.

The project consists of multiple parts:

+ `Website`
+ `Discord Bot`

The website can be found [here](https://kungdota.win/).

This repository includes the Discord bot component. It runs as dockerized containers on AWS EC2 that are automatically deployed when pushed to the repository.

The following need to be added as **Repository secrets:**
+ `AWS_KEY`
+ `EC2_INSTANCE`
+ `ENV`

This project consists of the following parts:

+ `Signup`
+ `Discord Commands`

Both parts use the following services:

+ `Kungsdota API`
+ `Steam API`
+ `Opendota API`
+ `AWS Dynamodb`

### Signup
Runs a weekly cron job that posts a list in a specified Discord channel. Users can interact with the list using embedded message buttons.

### Discord Bot
A Discord bot that allows actions such as:
+ `Add New Game`
+ `Shuffle Teams`
+ `Move Players After Shuffle`
+ `Search And Update New League Games`

Can be invoked using **Slash Commands.**