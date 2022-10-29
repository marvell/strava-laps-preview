FROM golang:1.19 as builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3
COPY --from=builder /app/strava-laps-preview /bin/
CMD ["/bin/strava-laps-preview"]