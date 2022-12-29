package yahoofin

import "go.uber.org/zap"

var _ Trading = &Data{}

type Data struct {
	logger *zap.Logger
}

func NewData(logger *zap.Logger) *Data {
	return &Data{
		logger: logger,
	}
}
