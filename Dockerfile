FROM golang:1.15 as base
RUN apt-get update && apt-get install -y libpcap-dev && rm -rf /var/lib/apt/lists/*

FROM base as build
WORKDIR /build
COPY . .
RUN go build

FROM base as runtime
WORKDIR /app
COPY --from=build /build/ground-control .
CMD ["./ground-control"]
