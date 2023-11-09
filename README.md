# Toresa
Toresa is a full-stack project. It pulls contents in the web.

# Target
* [ ] collector-manager: a daemon to control the collection job.
* [ ] toresa-shop: a front application to show the collection target.

# Quick Start
```bash
# run MongoDB in docker
docker run --name mongod -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root -v <your data storge path>:/data/db -p 27017:27017 -d mongo:7.0.0

# then you can connect to it using mongosh
docker run -it --rm mongo:7.0.0 mongosh --host 172.17.0.1 -u root -p root

# run the collector
go run ./cmd/collector/main.go -c=./cmd/collector/config.yaml --v=4
```
