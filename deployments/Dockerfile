FROM golang:1.19.0-alpine3.15 as build
WORKDIR /app
COPY go.mod go.sum ./
COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./vendor /app/vendor
RUN cd cmd/galleryapp && go build -o app

FROM alpine:3.15 as prod
COPY --from=build /app/cmd/github.com/Ghytro/galleryapp/app ./app
COPY ./web/ /web/
EXPOSE 3001
ENTRYPOINT ["/app"]
