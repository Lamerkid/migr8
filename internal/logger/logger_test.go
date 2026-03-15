package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	logg := NewLogger("WARN")
	logg.output = &buf

	logg.Debug("do not show")
	require.Empty(t, buf)

	logg.Info("do not show")
	require.Empty(t, buf)

	logg.Warn("show this")
	require.Contains(t, buf.String(), "[WARN]: show this\n")

	buf.Reset()
	logg.Error("and this")
	require.Contains(t, buf.String(), "[ERROR]: and this\n")
}
