package tradier

import "go.uber.org/zap"

var _ Trading = &Data{}

type Data struct {
	Logger *zap.Logger
}

func NewData(logger *zap.Logger) *Data {
	return &Data{
		Logger: logger,
	}
}
