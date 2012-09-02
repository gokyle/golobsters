# golobsters

### Overview
`golobsters` is an application that posts stories from 
[lobste.rs](https://lobste.rs) to [Twitter](https://twitter.com/lobsternews).

### Background
The first version was written in 92 source lines of code in Python, and is
a fairly basic system based on SQLite. I've been learning
[Go](http://www.golang.org) lately, and needed a project to work on, and
decided to rewrite lobsterpie to employ some of the fun parts of Go.

### Architecture
`golobsters` is comprised of two main components, `bot` (the backend) and
`frontend`. The backend employs a worker pool using goroutines and channels,
while the frontend simply displays the last time the bot updated.
