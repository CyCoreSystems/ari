package ari

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEndpoints(t *testing.T) {
	list, err := DefaultClient.ListEndpoints()
	assert.Nil(t, err, "Error getting list of endpoints")
	if len(list) > 0 {
		assert.NotNil(t, list[0].Resource)
	} else {
		fmt.Println("No endpoints received, but no errors.")
	}
}

func TestGetEndpointsByTech(t *testing.T) {
	list, err := DefaultClient.GetEndpointsByTech("PJSIP")
	assert.Nil(t, err, "Error getting list of endpoints with pjsip tech")
	for _, element := range list {
		fmt.Println("Endpoint: ", element)
	}
}

func TestGetEndpoint(t *testing.T) {
	endpoint, err := DefaultClient.GetEndpoint("PJSIP", "101")
	assert.Nil(t, err, "Error getting specific endpoint PJSIP/101")
	assert.NotNil(t, endpoint.Resource)
}
