FROM golang:1.19.0-alpine3.15 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd /app/cmd
COPY internal /app/internal
RUN cd cmd/github.com/Ghytro/galleryapp && go build -o app

FROM alpine:3.15 as prod
COPY --from=build /app/cmd/github.com/Ghytro/galleryapp/app ./app
COPY ./web/ /web/
EXPOSE 3001
ENTRYPOINT ["/app"]
