# docker-slacks

## Run

```
docker run \
  -h %H \
  -e WEBHOOK \
  -v /var/run/docker.sock:/var/run/docker.sock
  rlister/docker-slacks
```
