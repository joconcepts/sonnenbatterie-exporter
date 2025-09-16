FROM golang:1.24 as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /app/sonnenbatterie-exporter

FROM gcr.io/distroless/static AS final
COPY --from=build --chown=nonroot:nonroot /app/sonnenbatterie-exporter \
  /sonnenbatterie-exporter
CMD ["/sonnenbatterie-exporter"]
