package helper

import (
	"bufio"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"go.uber.org/zap"
)

func ReadConfig(configfile string) kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)

	// open the file
	file, err := os.Open(configfile)
	if err != nil {
		connection.Coldfinancelog().Debug("could not open file", zap.Any("error", err.Error()))
		return nil
	}
	defer file.Close()

	// scan items in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// read the lines
		lines := strings.TrimSpace(scanner.Text())
		// if there's content and its not a comment
		if len(lines) > 0 && !strings.HasPrefix(lines, "#") {
			val := strings.Split(lines, "=")
			parameter := val[0]
			value := val[1]
			m[parameter] = value
		}
	}

	if scanner.Err() != nil {
		connection.Coldfinancelog().Debug("error scanning file", zap.Any("error", scanner.Err()))
		os.Exit(1)
	}

	return m
}
