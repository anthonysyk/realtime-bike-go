package opendata

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDocksAvailability(t *testing.T) {
	docks, err := GetStationAvailability(100, 1400)
	assert.NoError(t, err)
	fmt.Println(docks.TotalCount, len(docks.Results), docks.HasNext)
}
