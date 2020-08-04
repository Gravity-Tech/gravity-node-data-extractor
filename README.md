# Gravity Node Data Extractor
Gravity Node Data Extractor

[![Build Status](https://drone.gravityhub.org/api/badges/Gravity-Tech/gravity-node-data-extractor/status.svg)](https://drone.gravityhub.org/Gravity-Tech/gravity-node-data-extractor)

## Common prerequisite

You need to have valid Binance API key data reading. (--api param)

## I. Running public Docker image

### 1. Pull the image

```
docker pull gravityhuborg/gravity-data-extractor:master
```

### 2. Run it

For example, on port 8090

```
docker run -itd -p 8090:8090 gravityhuborg/gravity-data-extractor:master
```

## II. Running the binary

### 1. Install Go

### 2. Build the binary

go build -o data-extractor

### 3. Run it with appropriate params

```
./data-extractor --api 'thekey' --pair 'WAVESBTC' --port 8099
```
