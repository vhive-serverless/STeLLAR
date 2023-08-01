package setup

import (
	"github.com/stretchr/testify/require"
	"stellar/setup"
	"testing"
)

func TestAssignEndpointIDs(t *testing.T) {
	endpointId := "endpointId"
	actual := &setup.SubExperiment{Parallelism: 3}
	actual.AssignEndpointIDs(endpointId)

	expected := &setup.SubExperiment{Parallelism: 3, Endpoints: []setup.EndpointInfo{{ID: "endpointId"}, {ID: "endpointId"}, {ID: "endpointId"}}}

	require.Equal(t, expected, actual)
}

func TestAddRoutes(t *testing.T) {
	actual := &setup.SubExperiment{Routes: []string{"route1", "route2"}}
	route := "route3"
	actual.AddRoute(route)

	expected := &setup.SubExperiment{Routes: []string{"route1", "route2", "route3"}}

	require.Equal(t, expected, actual)

	// nil route
	actual = &setup.SubExperiment{}
	actual.AddRoute(route)
	expected = &setup.SubExperiment{Routes: []string{"route3"}}

	require.Equal(t, expected, actual)
}
