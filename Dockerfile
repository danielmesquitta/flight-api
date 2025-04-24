FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN make build

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/tmp/server .
EXPOSE 8080
CMD ["./server"]
