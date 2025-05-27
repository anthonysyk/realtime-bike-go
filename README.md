# 🚲 Real-Time Bike Data Processing

This side project is a technical playground for practicing a **Go-based** stack with **MongoDB, GraphQL**.
It processes open data on bike availability in Paris, aggregating the number of available docks and bikes over time.

## 🛠️ Tech Stack
- **Go** – Backend services and data processing
- **MongoDB** – Stores raw and aggregated bike station data
- **GraphQL** – API layer for querying station data

## ⚙️ Process Overview
1. **Fetch real-time bike station data** at regular intervals and store it in **MongoDB**.
2. **Use MongoDB Change Streams** (watcher) to detect and record station update events.
3. **Aggregate updates** to compute rolling sums over **15, 30, and 60 minutes**.
4. **Expose aggregated data** through a **GraphQL API**.

## 🔗 Data Flow
1. The **Collector** service fetches data and writes it to MongoDB.
2. The **Watcher** listens to MongoDB **change streams** (oplog) and logs events to a dedicated collection.
3. The **TimeSeries** service aggregates those events into time-based summaries.
4. The **GraphQL API** exposes both raw and aggregated data to clients.

![schema-data-flow.png](doc%2Fschema-data-flow.png)
