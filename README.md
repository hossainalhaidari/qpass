# qpass

**qpass** is a minimal secret-sharing tool [inspired by Yopass](https://github.com/jhaals/yopass) that tries to keep the bare-minimum and stay as lightweight as possible

## Usage

- Start a memcached server: `docker run --name qpass_cache -p 11211:11211 -d memcached:alpine memcached`
- Run `go run .` to start the server at `http://localhost:3000`
- Or check the sample `docker-compose.yml` file for running it with Docker
