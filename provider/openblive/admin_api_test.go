package openblive

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAdminNames(t *testing.T) {
	result := getAdmins(3819533)
	require.Equal(t, 4, len(result))
}

func TestGetAdminNamesFails(t *testing.T) {
	result := getAdmins(3819531231231233)
	require.Equal(t, 1, len(result))
}
