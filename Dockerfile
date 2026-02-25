FROM golang:1.25 AS build

WORKDIR /src

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" \
    -o /out/steamcmd-healthcheck ./main.go

FROM scratch AS export

COPY --from=build /out/steamcmd-healthcheck \
    /steamcmd-healthcheck

ENTRYPOINT [ "/steamcmd-healthcheck" ]
