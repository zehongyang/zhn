package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T)  {
	D(zap.Any("1",1))
}
