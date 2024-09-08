# Stage 1: Build the Go binary
FROM golang:1.23.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the workspace
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Run the Go app in a minimal image
FROM surnet/alpine-wkhtmltopdf:3.20.2-0.12.6-small AS wkhtmltopdf
FROM alpine:3.20.2

# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    ca-certificates \
    fontconfig \
    # freetype \
    # ttf-dejavu \
    # ttf-droid \
    # ttf-freefont \
    # ttf-liberation \
    # more fonts
  # && apk add --no-cache --virtual .build-deps \
  #   msttcorefonts-installer \
  # Install microsoft fonts
  # && update-ms-fonts \
  # && fc-cache -f \
  # Clean up when done
  && rm -rf /tmp/* 
  # && apk del .build-deps

# Baixando e instalando a fonte Inter
RUN wget https://github.com/rsms/inter/releases/download/v3.19/Inter-3.19.zip -O /tmp/inter.zip && \
    mkdir -p /usr/share/fonts/inter && \
    unzip /tmp/inter.zip -d /usr/share/fonts/inter && \
    rm /tmp/inter.zip
# # Baixando e instalando a fonte Courier Prime
RUN mkdir -p /usr/share/fonts/courier-prime && \
    wget "https://github.com/google/fonts/raw/main/ofl/courierprime/CourierPrime-Regular.ttf" -O /usr/share/fonts/courier-prime/CourierPrime-Regular.ttf && \
    wget "https://github.com/google/fonts/raw/main/ofl/courierprime/CourierPrime-Bold.ttf" -O /usr/share/fonts/courier-prime/CourierPrime-Bold.ttf

# Atualizar o cache de fontes
RUN fc-cache -fv

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
# COPY --from=wkhtmltopdf /bin/wkhtmltoimage /bin/wkhtmltoimage
# COPY --from=wkhtmltopdf /lib/libwkhtmltox* /lib/

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
