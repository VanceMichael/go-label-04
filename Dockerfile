FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/airbridge ./cmd/server
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/airbridge /airbridge
COPY --from=build /src/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/airbridge"]
