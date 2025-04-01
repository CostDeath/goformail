FROM alpine:latest
WORKDIR /app

# Copy the executable binary
# If you do not have this file, please build the application first.
COPY src/goformail/goformail .

# Expose relevant ports for API & LMTP
EXPOSE 8000
EXPOSE 24

# Run app
ENTRYPOINT ["./goformail"]