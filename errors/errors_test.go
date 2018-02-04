package errors

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	Separator = ": "
}

func TestWrap(t *testing.T) {
	var (
		err     = New("base")
		wrapped = Wrap(err, "wrapping")
	)
	require.Equal(t, "base", err.Error())
	require.Equal(t, "wrapping: base", wrapped.Error())
}

func TestWrapf(t *testing.T) {
	var (
		now     = time.Now()
		err     = New("base")
		wrapped = Wrapf(err, "wrapping at %s", now)
	)

	require.Equal(t, fmt.Sprintf("wrapping at %s: base", now), wrapped.Error())
}

func TestErrorf(t *testing.T) {
	var (
		now = time.Now()
		err = Errorf("error at %s", now)
	)

	require.Equal(t, fmt.Sprintf("error at %s", now), err.Error())
}
