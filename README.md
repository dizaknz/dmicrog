# Dmicrog Service

Event RPC service using [go-micro](https://github.com/micro/go-micro)

## Description

Demo event RPC service - WIP using default pub/sub configuration.

Skeleton project generated with `micro new`

## Quick start

```
make up
```

or 

```
CGO_ENABLED=0 go build
docker-compose build
docker-compose up -d
```

to shutdown 

```
make down
```

## Configuration

```
./dmicrog -h
```

Standard micro environmental vars.

## Build it locally

Build the binary

```
make proto
make build
```

Run the service
```
./dmicrog
```
