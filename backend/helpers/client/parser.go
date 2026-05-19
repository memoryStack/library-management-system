package client

import (
	"sync"

	"github.com/ua-parser/uap-go/uaparser"
)

var (
	parserOnce sync.Once
	parserInst *uaparser.Parser
	parserErr  error
)

// Parser returns a shared UA parser (regex DB embedded in uap-go).
func Parser() (*uaparser.Parser, error) {
	parserOnce.Do(func() {
		parserInst, parserErr = uaparser.NewFromSaved()
	})
	return parserInst, parserErr
}
