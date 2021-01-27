FROM golang:1.14-alpine
#RUN apk add --no-cache youtube-dl && apk add --no-cache ffmpeg
WORKDIR /usr/src
COPY . .
RUN go get ./...
RUN go build -o /app .

FROM jauderho/youtube-dl
COPY --from=0 /app /app
ENTRYPOINT ["/app"]