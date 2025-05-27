package graph

import "realtime-bike-go/internal/store"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store *store.Store
}
