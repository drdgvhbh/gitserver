package git_test

import (
	"testing"

	"github.com/drdgvhbh/gitserver/internal/git"

	"github.com/stretchr/testify/assert"
)

func TestActionString(t *testing.T) {
	assert := assert.New(t)

	insertStr, err := git.Insert.String()
	assert.NoError(err)
	assert.EqualValues("INSERT", insertStr)

	modifyStr, err := git.Modify.String()
	assert.NoError(err)
	assert.EqualValues("MODIFY", modifyStr)

	deleteStr, err := git.Delete.String()
	assert.NoError(err)
	assert.EqualValues("DELETE", deleteStr)
}
 