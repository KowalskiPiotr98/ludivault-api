FROM golang:1.23 AS build
WORKDIR /build
COPY . .
RUN go build

FROM golang:1.23 AS run
EXPOSE 5500
ENV LUDIVAULT_LISTEN=0.0.0.0:5500

ARG user=ludivault
RUN adduser --shell /bin/false --no-create-home --disabled-password --disabled-login "$user"

WORKDIR /ludivault
COPY --from=build --chown=$user:$user --chmod=100 /build/ludivault .

USER $user

ENTRYPOINT ["/ludivault/ludivault"]