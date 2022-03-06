# BOXYD POC

Proof of concept for Boxyd, a simpler way of moving.

## Run by

`docker build -t boxyd_poc:version .`

`docker run -p=8080:8080 --env BOXYD_USERNAME=username --env BOXYD_PASSWORD=password --rm -d boxyd_poc:version`

