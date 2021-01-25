FROM golang:1.14-alpine
RUN mkdir /src
WORKDIR /src
COPY . .
RUN go get ./...
RUN go build -o /app .

# vimagick/youtube-dl contains ffmpeg, youtube-dl.
FROM vimagick/youtube-dl
COPY --from=0 /app /app
ENTRYPOINT ["/app"]