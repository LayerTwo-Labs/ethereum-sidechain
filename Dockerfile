# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""


# Avoid alpine, as that doesn't play nice with C stuff
# Build sidegeth in a stock Go builder container
FROM golang:1.18 as builder

RUN apt install git

# Get Rust
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-ethereum/
COPY go.sum /go-ethereum/
RUN cd /go-ethereum && go mod download

ADD . /go-ethereum
WORKDIR /go-ethereum
RUN cargo build --manifest-path ./drivechain/Cargo.toml
RUN go run build/ci.go install ./cmd/sidegeth

# Pull sidegeth into a second stage deploy container
# Avoid alpine, as that doesn't play nice with C stuff
FROM debian:bookworm-slim

COPY --from=builder /go-ethereum/build/bin/sidegeth /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["sidegeth"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
