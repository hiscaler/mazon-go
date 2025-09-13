package mazon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shippingLabelService_Detail(t *testing.T) {
	label, err := client.Services.ShippingLabel.Detail(ctx, ShippingLabelDetailRequest{OrderCode: "EPB00120250912114236000021"})

	assert.Nil(t, err)
	assert.NotEmpty(t, label)
}

func Test_shippingLabelService_QueryWithBadTrackingNumbers(t *testing.T) {
	_, err := client.Services.ShippingLabel.Query(ctx, "", " ")

	assert.NotNil(t, err)
}

func Test_shippingLabelService_Query(t *testing.T) {
	labels, err := client.Services.ShippingLabel.Query(ctx, "9234690397703300025653")
	assert.Nil(t, err)
	if err == nil {
		assert.NotEmpty(t, labels)
		if len(labels) != 0 {
			label := labels[0]
			assert.Equal(t, "9234690397703300025653", label.Labels[0].TrackingNumber)
		}
	}
}
