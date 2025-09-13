package mazon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_userService_Info(t *testing.T) {
	a, err := client.Services.User.Information(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "EPB001", a.Code)
}
