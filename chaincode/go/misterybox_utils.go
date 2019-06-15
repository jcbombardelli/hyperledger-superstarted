package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ResponseSuccess struct {
	ID        string    `json:"id"`
	Message   string    `json:"localizacao"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

//JSONResponseSuccess - Json Response
func JSONResponseSuccess(id string, message string, datetime time.Time) []byte {

	response := ResponseSuccess{
		ID:        id,
		Message:   message,
		Timestamp: datetime,
	}
	json, _ := json.Marshal(response)
	return []byte(json)
}

//JSONResponseError - JSON Response
func JSONResponseError(id string, message string, code int) string {
	return fmt.Sprintf("{ \"type\":\"error\", \"code\":\"%v\", \"id\":\"%s\", \"message\":\"%s\"}", code, id, message)
}

//ConstructQueryResponseFromIterator - Used By Query String
func ConstructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (string, *bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	if resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		buffer.WriteString(string(queryResponse.Value))
		return queryResponse.GetKey(), &buffer, err
	}

	return "", &buffer, nil
}
