package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract define the Chaincode structure
type SmartContract struct {
}

// Misterybox define the asset structure, with 10 properties. Structure tags are used by encoding/json library
type Misterybox struct {
	DocType  string    `json:"docType"`
	Serial   string    `json:"serial"`
	Size     string    `json:"type"`
	Model    string    `json:"model"`
	Owner    string    `json:"owner"`
	Register time.Time `json:"registerAt"`
}

/*
 * Init method is called when the Smart Contract "misterybox" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "misterybox"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryMisterybox" {
		return s.queryMisterybox(APIstub, args)

	} else if function == "createMisterybox" {
		return s.createMisterybox(APIstub, args)

	} else if function == "queryAllMisteryboxes" {
		return s.queryAllMisteryboxes(APIstub)

	} else if function == "transferMisterybox" {
		return s.transferMisterybox(APIstub, args)

	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryMisterybox(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	cmtbxAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(cmtbxAsBytes)
}

func (s *SmartContract) createMisterybox(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Error: Incorrect number of arguments.")
	}

	misteryboxIsExists := false
	var mtbx = Misterybox{}

	queryString := fmt.Sprintf(`{"selector": {"serial": "%s" }}`, args[0])
	key, queryResults, _ := getQueryResultForQueryString(APIstub, queryString)

	if key != "" {
		misteryboxIsExists = true
		err := json.Unmarshal(queryResults, &mtbx)
		if err != nil {
			return shim.Error(JSONResponseError(args[0], "Error on parse json", 5))
		}
	}

	if !misteryboxIsExists {

		if args[0] != "" {
			mtbx.Serial = args[0]
		}
		if args[1] != "" {
			mtbx.Size = args[1]
		}
		if args[2] != "" {
			mtbx.Model = args[2]
		}
		if args[3] != "" {
			mtbx.Owner = strings.ToLower(args[3])
		}

		mtbx.DocType = "MisteryBox"
		mtbx.Register = time.Now()
		mtbxAsBytes, _ := json.Marshal(mtbx)

		doctypeTxIndexKey, err := APIstub.CreateCompositeKey("doctype~tx", []string{mtbx.DocType, APIstub.GetTxID()})
		if err != nil {
			return shim.Error(JSONResponseError(stub.GetTxID(), PUTSTATE_ERROR_CODE, err.Error()))
		}

		//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the ccee.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		err = stub.PutState(hashTxIndexKey, value)
		if err != nil {
			return shim.Error(JSONResponseError("", "Error on PutState : "+err.Error(), 99))
		}

		err = APIstub.PutState(APIstub.GetTxID(), mtbxAsBytes)
		if err != nil {
			return shim.Error(JSONResponseError("", "Error on PutState : "+err.Error(), 99))
		}

		return shim.Success(JSONResponseSuccess(APIstub.GetTxID(), mtbx.DocType, mtbx.Register))

	} else {
		return shim.Error(JSONResponseError(APIstub.GetTxID(), "Misterybox has already exists!", 409))
	}

}

func (s *SmartContract) transferMisterybox(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error(JSONResponseError("", "Incorrect number of arguments. Expecting 3", 2))
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "" {
			return shim.Error(JSONResponseError("", "All fields are required", 3))
		}
	}
	queryString := fmt.Sprintf(`{"selector": {"_id": "%s" }}`, args[0])
	key, queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(JSONResponseError(args[0], "ID referenced does not exist", 4))
	}

	if key != "" {
		misterybox := Misterybox{}
		err = json.Unmarshal(queryResults, &misterybox)
		if err != nil {
			return shim.Error(JSONResponseError(args[0], "Error on parse json", 5))
		}

		if misterybox.Owner != args[1] {
			return shim.Error(JSONResponseError(args[0], "Owner Incorrect", 20))
		}

		if misterybox.Owner == args[2] {
			return shim.Error(JSONResponseError(args[0], "New Owner is same actual owner", 21))
		}

		mtbxAsBytes, _ := json.Marshal(misterybox)
		APIstub.PutState(key, mtbxAsBytes)

		return shim.Success(JSONResponseSuccess(APIstub.GetTxID(), fmt.Sprintf("The new Owner is %s", args[2]), time.Now()))

	}
	return shim.Error(JSONResponseError(args[0], "Key referenced does not exist", 4))
}

func (s *SmartContract) queryAllMisteryboxes(APIstub shim.ChaincodeStubInterface) sc.Response {

	queryString := `{"selector": { "docType": "MisteryBox" } } }`
	resultsIterator, err := APIstub.GetQueryResult(queryString)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	//fmt.Printf("- queryAllMisteryboxes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) (string, []byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return "", nil, err
	}
	defer resultsIterator.Close()

	key, buffer, err := ConstructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", nil, err
	}

	return key, buffer.Bytes(), nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
