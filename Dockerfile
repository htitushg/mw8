# syntax=docker/dockerfile:1
FROM golang:1.22.0 

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
#RUN useradd henry 
RUN useradd -m -s /bin/bash henry
RUN mkdir /logs
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY assets/* ./assets/
COPY assets/css/* ./assets/css/
COPY cmd/*.go ./cmd/
COPY config/*.json ./config/
COPY controllers/*.go ./controllers/
COPY data/* ./data/
COPY database/* ./database/
COPY internal/middlewares/*.go ./internal/middlewares/
COPY internal/models/*.go ./internal/models/
COPY internal/utils/*.go ./internal/utils/
COPY logs/* ./logs/
RUN chown henry:henry -R /app/logs
RUN chmod 666 /app/logs
COPY router/* ./router/
COPY server/* ./server/
COPY templates/* ./templates/
#COPY templates/layouts/* ./templates/layouts/

VOLUME /app/logs
VOLUME /app/data
VOLUME /app/templates
VOLUME /app/assets

#RUN chmod 777 /usr/local/go
#RUN chmod 777 /usr/local/go/bin/go
#USER henry
# /var/lib/docker/volumes/logs/ à l'extérieur du container
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o mw7 cmd/main.go


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080
USER henry

ENV MYSQL_HOST db
ENV MYSQL_PORT 3306
ENV MYSQL_DATABASE basemw7
ENV MYSQL_USER henry
ENV MYSQL_PASSWORD 1nhri96p
# Run
CMD ["./mw7"]
# Pour construire le container avec le nom mw7
# docker build -t mw7 .
# Pour lancer l'application 
# docker run --name mw7 -it --rm -p 8080:8080 -v /home/henry/go/src/mw7/logs:/app/logs -v /home/henry/go/src/mw7/data:/app/data -v /home/henry/go/src/mw7/templates:/app/templates -v /home/henry/go/src/mw7/assets:/app/assets mw7
