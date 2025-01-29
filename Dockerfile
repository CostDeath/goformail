FROM alpine:latest
WORKDIR /app

# Copy the executable binary
# If you do not have this file, please build the application first.
COPY src/goformail/goformail .

# Expose API port
EXPOSE 8000

# Run app
ENTRYPOINT ["goformail"]