package areship

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_userInfoService_Detail(t *testing.T) {
	a, err := client.Services.UserInfo.Detail(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "J5820", a.Code)
}
