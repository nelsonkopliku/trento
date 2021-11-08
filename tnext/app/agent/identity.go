package agent

import (
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

const machineIdPath = "/etc/machine-id"

var fileSystem = afero.NewOsFs()

func loadIdentifier() (uuid.UUID, error) {
	var agentID uuid.UUID

	machineIDBytes, err := afero.ReadFile(fileSystem, machineIdPath)
	if err != nil {
		return agentID, err
	}
	machineID := strings.TrimSpace(string(machineIDBytes))

	agentID, err = uuid.NewRandomFromReader(strings.NewReader(machineID))

	if err != nil {
		return agentID, err
	}

	return agentID, nil
}
