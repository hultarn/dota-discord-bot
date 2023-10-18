#!/bin/sh

# Exit on any error
set -e

mkdir "$HOME/.ssh"
echo "$AWS_KEY" > "$HOME/.ssh/aws_key"
chmod 400 "$HOME/.ssh/aws_key"
echo "ec2-51-20-63-178.eu-north-1.compute.amazonaws.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBPizi++gnfLgkIW3qHSViUeaUXTgtk9RVzwpzPQUXMU/7CF3+i/08WHdKhWfXlSyrAt1u5oGFGxsv2Wgoyvc3+Y=" > "$HOME/.ssh/known_hosts"
echo "$ENV" > .env

# Show commands as they're executed
set -x
scp -i "$HOME/.ssh/aws_key" .env "$EC2_INSTANCE"
ssh -i "$HOME/.ssh/aws_key" "$EC2_INSTANCE" sudo systemctl stop dbot-signup
ssh -i "$HOME/.ssh/aws_key" "$EC2_INSTANCE" sudo systemctl stop dbot-league
sleep 2 # some downtime
scp -i "$HOME/.ssh/aws_key" dbot "$EC2_INSTANCE":/home/ubuntu/dbot
sleep 2
ssh -i "$HOME/.ssh/aws_key" "$EC2_INSTANCE" sudo systemctl start dbot-signup
ssh -i "$HOME/.ssh/aws_key" "$EC2_INSTANCE" sudo systemctl start dbot-league
