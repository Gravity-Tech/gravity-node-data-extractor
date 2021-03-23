# Gravity Node Data Extractor
Gravity Node Data Extractor

[![Build Status](https://drone.gravityhub.org/api/badges/Gravity-Tech/gravity-node-data-extractor/status.svg)](https://drone.gravityhub.org/Gravity-Tech/gravity-node-data-extractor)

## Abstract 

This package represents working example of Gravity Data Extractor.

## I. Running public Docker image

### 1. Pull the image

```
docker pull gravityhuborg/gravity-data-extractor:master
```

### 2. Run

You can run SuSy dApp extractor

```
docker run -itd -p 8090:8090 gravityhuborg/gravity-data-extractor:master
```

### 2.1 Run SuSy dApp extractor for EVM-based or WAVES-based chains

Currently, extractor supports only EVM-chain<->WAVES direction.
So we in order to add new custom(non-EVM) chain we need to provide 4 additional params
and 2 bridge implementations as well - for direct and reversal swaps.

So, for EVM-based tokens we run `--waves-based-to-eth-direct` and '--waves-based-to-eth-reverse'.
For WAVES-based tokens we run `--eth-based-to-waves-direct` and `--eth-based-to-waves-reverse'


## II. Running the binary

### 1. Install Go

### 2. Build the binary

```
go build -o data-extractor
```

### 3. Run it with appropriate params

```
./data-extractor --api 'thekey' --pair 'WAVESBTC' --port 8099
```

## Miscellaneous

### 1. Run tests

```
go test tests/*.go -v
```
