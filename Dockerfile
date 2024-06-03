FROM golang:1.22-bookworm

RUN alias "ll=ls -al" \
  && apt update \
  && apt upgrade -y \
  && apt install -y sqlite3
