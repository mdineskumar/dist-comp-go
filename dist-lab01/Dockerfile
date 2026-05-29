FROM golang:1.26.3

# Install networking tools for Experiments D and E
RUN apt-get update && apt-get install -y \
    iptables \
    iproute2 \
    iputils-ping \
    net-tools \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /lab01
# Keep the container running
CMD ["tail", "-f", "/dev/null"]