package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/jonlinkens/go-markdown/lexer"
)

type TokenSlice []lexer.Token

func (tokens TokenSlice) ToJson() (string, error) {

	jsonData, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON array:", err)
		return "", err
	}

	return string(jsonData), nil
}
