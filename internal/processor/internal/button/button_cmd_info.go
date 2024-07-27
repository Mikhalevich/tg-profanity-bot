package button

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type ButtonCMDInfo struct {
	CMD     string `json:"cmd"`
	Payload []byte `json:"payload"`
}

func FromBase64(data string) (ButtonCMDInfo, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ButtonCMDInfo{}, fmt.Errorf("base64 decode: %w", err)
	}

	var buttonInfo ButtonCMDInfo
	if err := json.Unmarshal(b, &buttonInfo); err != nil {
		return ButtonCMDInfo{}, fmt.Errorf("json unmarshal: %w", err)
	}

	return buttonInfo, nil
}

func (bci ButtonCMDInfo) ToBase64() (string, error) {
	b, err := json.Marshal(bci)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
