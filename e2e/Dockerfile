FROM golang:1.14

WORKDIR /app

RUN curl -o /bin/kubectl -SL https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x /bin/kubectl

# Better dep caching
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# TODO: create way to download specific versions
RUN go build -o /bin/riser
