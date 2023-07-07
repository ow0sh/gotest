FROM migrate/migrate
RUN apk add wget
RUN wget https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    mv jq-linux64 /usr/local/bin/jq && \
    chmod +x /usr/local/bin/jq
WORKDIR /go/src/github.com/ow0sh/gotest
COPY migrations ./migrations
COPY migrations/*.sql ./migrations
COPY scripts/migrate.sh ./migrate.sh
# COPY scripts/dbconfig.yml ./dbconfig.yml
RUN chmod +x migrate.sh
ENTRYPOINT [ "./migrate.sh" ]