package src

import (
	"bytes"
	"errors"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/setup"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCompleteFlow_With_Test_Dependencies(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	envConf := config.Load(os.Getenv(""))

	testDependencies := setup.InitializeTestDependencies(envConf)

	err := testDependencies.ProcessFile.Execute("./../path/txns1.csv")

	// Assert no errors occurred
	assert.NoError(t, err)

	// Check if specific log messages were shown
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Reading file")
	assert.Contains(t, logOutput, "Storing 4 transactions")
	assert.Contains(t, logOutput, "Sending mail with summary account information")
}

func Test_InvalidPathFile(t *testing.T) {
	envConf := config.Load(os.Getenv(""))

	testDependencies := setup.InitializeTestDependencies(envConf)

	err := testDependencies.ProcessFile.Execute("./path/non-exist.csv")

	// Assert no errors occurred
	assert.Error(t, err, errors.New("failed to open file ./path/non-exist.csv: open ./path/non-exist.csv: no such file or directory"))
}
