####################################
# STEP 1 build executable binary
####################################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
RUN apk add ca-certificates
WORKDIR $GOPATH/src/accelbyte/metric_summary

#copy all the content to container
COPY . .

# Build the binary
RUN export CGO_ENABLED=0 && go build -o /go/bin/metric_summary

#change the permission on binary
RUN chmod +x /go/bin/metric_summary

##############################################
# STEP 2 build a small image using alpine:3.14
##############################################
FROM alpine:3.14

# Copy our static executable.
COPY --from=builder /go/bin/metric_summary ./metric_summary

# Run the entrypoints.
ENTRYPOINT [ "./metric_summary" ]