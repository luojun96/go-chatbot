package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// create a test for the function New()
const conf = `
openai:
  token: "test"
`

// write the configuration to the specified file
func writeConfigFile(file, content string) error {
	return os.WriteFile(file, []byte(content), 0644)
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	confFile := tmpDir + "/conf.yaml"
	err := writeConfigFile(confFile, conf)
	assert.Nil(t, err)

	c, err := New(&confFile)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "test", c.OpenAI.Token)
}

func TestGetYamlFile(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// test empty dir
	_, err := getYamlFile(tmpDir)
	assert.NotNil(t, err)
}
