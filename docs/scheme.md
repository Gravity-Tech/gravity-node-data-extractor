
# Gravity built-in data extractor architecture

## Abstract

This document represents UML architecture of Gravity build-in data extractor. The author of document, architecture and implementation is **Gravity Core Team**. Current implementation is *first* and *initial* so it's considered as *built-in*.

## Foreword

Architecture provided in this document is served only as a ***best practise***  for Gravity protocol node operators and contributors. 

Violating the idea and object relations can lead to unexpected behaviour of Gravity Node system parts.

## Scheme Overview

![Gravity data extractor scheme](https://i.imgur.com/xkkFsrU.jpg)


## Main concepts

The main concepts behind the subject implementation we stick to are:
1. Provide stateless system & Avoid data mutability
2. Conformance to ***any*** kind of data (reusability)
3. Manifest available operations

### Stateless VS stateful

Talking about modern system designing we tend to outlook for certain balance in how operate with objects.

The vast majority of applications combine both stateless and stateful parts of the system.

As regards *Gravity protocol*, it's not an exception. Such parts of the system that are responsible for data mutability and storing are considered stateful. 

By design, extractors are aimed to perform only data aggregation and mapping procedures. That is why we have stateless approach. They are not awared of particular data consumers, and here comes the second advantage - *Extractors are isolated.*


### Conformance to any data & Reusability

Current architecture gives an ability to ***transform*** and ***transport*** data ***the way we want***. There are responsible interfaces for that:
1. Extractor<T, R> - This interface declares methods on how we ***transform*** data .
2. IDataBridge<T, R> - This interface declares methods on how we ***transport*** the data.

Generic types provided in declarations represent:
1. T - raw data type
2. R - transformed response data type. 

Furthermore, such approach gives us an ability to represent any kind of data and deliver it differently.

The scheme represents possible implementations:

![Gravity data extractor response controller examples](https://i.imgur.com/RnPi1Kw.png)

Existing ***reusability of every distinct system part*** concludes to the ***reusability of whole system***.

### Available operations

The design requires 3 accessible access endpoints to be implemented. Those are for:
1. Raw data fetching.
2. Mapped/transformed data fetching.
3. Extractor info in JSON format, containing data feed tag and description.


