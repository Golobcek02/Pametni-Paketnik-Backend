FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN apt-get update && \
    apt-get install -y wget build-essential libssl-dev zlib1g-dev libbz2-dev \
    libreadline-dev libsqlite3-dev llvm libncurses5-dev libncursesw5-dev \
    xz-utils tk-dev libffi-dev liblzma-dev python3-openssl git

RUN wget https://www.python.org/ftp/python/3.11.3/Python-3.11.3.tgz && \
    tar xvf Python-3.11.3.tgz && \
    cd Python-3.11.3 && \
    ./configure --enable-optimizations && \
    make altinstall && \
    cd .. && \
    rm -rf Python-3.11.3.tgz Python-3.11.3

ENV PATH="/usr/local/bin:${PATH}"

RUN ln -s /usr/local/bin/python3.11 /usr/local/bin/python

RUN python -m ensurepip && \
    python -m pip install --upgrade pip && \
    python -m pip install tensorflow opencv-python scikit-learn

RUN CGO_ENABLED=0 GOOS=linux go build -o /ppbackend ./main/main.go
EXPOSE 5551
CMD ["/ppbackend"]

