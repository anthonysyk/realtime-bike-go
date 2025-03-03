# 🚲 Real-Time Bike Data Processing

This side project is a technical playground for practicing a **Go-based** stack with **MongoDB, Kafka, GraphQL, Protobuf, and Redis**.
It processes open data on bike availability in Paris, aggregating the number of available docks and bikes over time.

## 🛠️ Tech Stack
- **Go** – Backend processing
- **MongoDB** – Storing bike station states
- **Kafka** – Event-driven messaging
- **GraphQL** – API for querying data
- **Protobuf** – Efficient data serialization
- **Redis** – Caching for fast data access

## ⚙️ Process Overview
1. **Fetch bike station data** at a **given time** and store it in **MongoDB**.
2. **Generate Kafka events** tracking changes (+/- bikes, +/- docks) per station.
3. **Aggregate events** to compute rolling sums over **15, 30, and 60 minutes**.
4. **Expose aggregated data** via a **GraphQL API**.  

## 🔗 **Data Flow**
1. **`Collector`** writes station data to MongoDB.
2. **`Watcher`** captures changes (`oplogs`) and pushes them to Kafka.
3. **`Aggregator`** consumes events from Kafka and generates a global state.
4. **`Exporter`** updates Redis with the aggregated state.
5. **`GraphQL API`** exposes this data to clients.  

![schema-data-flow.png](doc%2Fschema-data-flow.png)
