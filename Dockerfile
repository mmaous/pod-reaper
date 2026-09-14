FROM golang:1.27.0-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /pod-reaper .

FROM gcr.io/distroless/static-debian12:latest
COPY --from=build /pod-reaper /pod-reaper
ENTRYPOINT ["/pod-reaper"]
