FROM golang:1.17-alpine as builder

WORKDIR /app
RUN apk add --no-cache curl

RUN curl -o /bin/kubectl -SL https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x /bin/kubectl


FROM golang:1.17-alpine as runtime

RUN apk add --no-cache git

COPY --from=builder /bin/kubectl /bin/kubectl

WORKDIR /app

# Better dep caching
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /bin/riser -tags=e2e


