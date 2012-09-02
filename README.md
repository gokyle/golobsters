# golobsters

### Overview
`golobsters` is an application that posts stories from 
[lobste.rs](https://lobste.rs) to [Twitter](https://twitter.com/lobsternews).

There is an experimental [web interface](http://lobsters.kyleisom.net)
where I am learning how to do that type of work in Go.

### Background
The first version was written in 92 source lines of code in Python, and is
a fairly basic system based on SQLite. I've been learning
[Go](http://www.golang.org) lately, and needed a project to work on,
preferably writing some type of webapp as that is an area I have no experience
in. Ergo, the decision to rewrite lobsterpie in Go.

### Architecture
`golobsters` is comprised of two main components, `bot` (the backend) and
`frontend`. The backend employs a worker pool using goroutines and channels,
while the frontend simply displays the last time the bot updated.

The backend starts up the worker pool, which communicate via a channel for
new stories. This channel is written to by the RSS feed parser, which is in
a separate goroutine. When the RSS feed is updated, it reads the entries,
converts them to the internal data structure that represents a story, and
writes that to the channel. The next available worker picks it up, checks
whether the story has already been posted, and if not, posts the story to
Twitter and marks the story as read in the database. Reading from a channel
blocks, so each worker essentially sleeps while waiting for new stories.

### Deployment
`golobsters` currently runs on [Heroku](https://www.heroku.com), and uses
a [Postgres](http://www.postgres.org) database stored on my postgres VPS. As
this application effectively only has one user, the transaction cost isn't
an issue.

I had originally wanted to use [Redis](https://www.redis.io) as my datastore,
but Heroku's Redis addon costs money. This app generates no revenue, so I
didn't want to end up paying for keeping it running; Redis also requires
an SSL tunnel in order to communicate securely with remote datastores which
would end up being a hassle to setup. I already had a VPS set up with Postgres,
so I ended up just using that. It is overkill, as the table merely stores the
guid of the story, but the infrastructure was already in place.

I am using the Heroku Go buildpack, which so far seems to run pretty smoothly;
setting up the deployment to Heroku was pretty easy and very
[straightforward](https://gist.github.com/299535bbf56bf3016cba).


