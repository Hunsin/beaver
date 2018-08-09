package logger_test

import (
	"log"

	"github.com/Hunsin/beaver/logger"
)

func Example() {

	// create a logger by New()
	l := logger.New()

	// set flags via standard package flags
	l.Flags(log.Lshortfile)

	// set the levels to log
	l.Level(logger.Linfo | logger.Lerror | logger.Lfatal)

	// no output if the level not set
	l.Debug("Hello World!")
	// output:

	l.Info("Hello!")
	// output: example_test.go:24: INFO : Hello!
}
