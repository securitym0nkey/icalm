FROM golang:1.22.3 as build
ARG tags=docker-build-versionless
WORKDIR /go/src/app
COPY . .
ENV TAGS=${tags}
RUN CGO_ENABLED=0 go build -tags ${TAGS} -o bin/ ./cmd/icalm-server/

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /go/src/app/bin/* /usr/bin/
ENTRYPOINT ["/usr/bin/icalm-server"]
CMD ["-networks", "/home/nonroot/icalm/networks.csv", "-line-listen", "4226"]

