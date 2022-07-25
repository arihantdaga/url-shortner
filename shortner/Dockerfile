from golang:1.17.11-alpine3.16 AS builder
#create app user
RUN adduser -D -g '' appuser
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify
# copy the source code as the last step
COPY . .
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main .


#build a small image
FROM alpine:3.14.1
LABEL language="golang"
# import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
# copy the executable
COPY --from=builder --chown=appuser:1000 /app/main /main
# use a non-root user
USER appuser

ENTRYPOINT ["./main"]
