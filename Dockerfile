FROM golang:1.18 AS builder
WORKDIR /src

COPY . .
RUN make build-binary

FROM ubuntu:focal AS final-image

RUN apt-get update && apt-get install --yes --no-install-recommends ca-certificates

COPY --from=builder /src/build/aspire-loan-app /aspire-loan-app
COPY --from=builder /src/migrations/db_setup.sql /db_setup.sql

CMD ["./aspire-loan-app"]