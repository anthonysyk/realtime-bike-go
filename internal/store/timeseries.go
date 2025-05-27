package store

import (
	"context"
	"realtime-bike-go/cmd/graphqlapi/graph/model"
	"time"
)

func ListStationTimeSeries(ctx context.Context, from, to time.Time, period model.WindowPeriod) {

}

/**
peux tu ajouter une logique supplémentaire mais dans une fonction séparée, cela consite à prendre la liste des opendata.Station et la window (qui peut être 15 minutes par exemple), et il faudra que tu remplisse les données manquantes toutes les 15 minutes entre 2 dates, et si ça n'existe pas du prends la valeur la plus récente et si aucune données tu remplis quand meme avec des 0 pour
*/
