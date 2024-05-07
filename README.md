# r/place

A toy version of [Reddit r/place](https://www.reddit.com/r/place/), a
collaborative real-time canvas.

## Quickstart

- Start a local Redis server on port 6379.
- Start the server by `go run .`
- Open `localhost:8080` in the browser.

## TODO

- [ ] Make the cooldown bar more friendly
- [ ] Implement rate limit in the backend
- [ ] Fix: WebSocket interrupted on loading after long idle

## Reference

[How We Built r/Place](https://www.redditinc.com/blog/how-we-built-rplace)
