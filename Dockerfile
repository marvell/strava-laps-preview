FROM golang:1.18 as builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build .

FROM alpine:3
COPY --from=builder /app/strava-laps-preview /bin/
CMD ["/bin/strava-laps-preview"]