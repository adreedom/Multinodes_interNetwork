package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type SwiftMessage struct {
	MsgID   string `json:"msgid"`
	MsgContent  string `json:"msgcontent"`
	MsgStatus string `json:"ssgstatus"`
	MsgType  string `json:"msgtype"`
	TimeStamp string `json:"timestamp"`
}

type Assets struct {
	AssetName   string  `json:"assetname"`
	AssetID     string  `json:"assetid"`
	AssetAmount float64 `json:"assetamount"`
}

type Clients struct {
	Name            string   `json:"name"`
	ClientID        string   `json:"clientid"`
	ClientType      string   `json:"clienttype"`
	SafeKeepingAccount string   `jsoasset_ton:"safekeepingaccount"`
	Currency        string   `json:"currency"`
	Asset           []Assets `json:"asset"`
	Status          string   `json:"status"`
}

type Transactions struct {
	TsID      string `json:"tsid"`
	Seller        string `json:"seller"`
	Buyer        string `json:"buyer"`
	Asset     Assets `json:"asset"`
	TimeStamp string `json:"timestamp"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "doInnerTransaction" {
		return s.doInnerTransaction(APIstub, args)
	} else if function == "queryAllAsset" {
		return s.queryAllAsset(APIstub, args)
	} else if function == "queryTransactions" {
		return s.queryTransactions(APIstub, args)
	} else if function == "queryClientInfo" {
		return s.queryClientInfo(APIstub, args)
	} else if function == "setSwiftMessage" {
	    return s.setSwiftMessage(APIstub, args)
	} else if function == "getSwiftMessage" {
	    return s.getSwiftMessage(APIstub, args)
	} else if function == "updateSwiftMessage" {
	    return s.updateSwiftMessage(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryClientInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	clientinfo_bytes, _ := APIstub.GetState(args[0])
	return shim.Success(clientinfo_bytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	asset1 := []Assets{
		Assets{AssetName: "ATVI", AssetID: "A1", AssetAmount: 2500.00},
		Assets{AssetName: "BABA", AssetID: "A2", AssetAmount: 2500.00},
		Assets{AssetName: "JD", AssetID: "A3", AssetAmount: 2500.00},
		Assets{AssetName: "AMD", AssetID: "A4", AssetAmount: 2500.00},
		Assets{AssetName: "INTEL", AssetID: "A5", AssetAmount: 2500.00},
	}

	asset2 := []Assets{
		Assets{AssetName: "ATVI", AssetID: "A1", AssetAmount: 1500.00},
		Assets{AssetName: "BABA", AssetID: "A2", AssetAmount: 1500.00},
		Assets{AssetName: "JD", AssetID: "A3", AssetAmount: 1500.00},
		Assets{AssetName: "AMD", AssetID: "A4", AssetAmount: 1500.00},
		Assets{AssetName: "INTEL", AssetID: "A5", AssetAmount: 1500.00},
	}

	asset3 := []Assets{
		Assets{AssetName: "ATVI", AssetID: "A1", AssetAmount: 500.00},
		Assets{AssetName: "BABA", AssetID: "A2", AssetAmount: 500.00},
		Assets{AssetName: "JD", AssetID: "A3", AssetAmount: 500.00},
		Assets{AssetName: "AMD", AssetID: "A4", AssetAmount: 500.00},
		Assets{AssetName: "INTEL", AssetID: "A5", AssetAmount: 500.00},
	}

	asset4 := []Assets{
		Assets{AssetName: "ATVI", AssetID: "A1", AssetAmount: 1500.00},
		Assets{AssetName: "BABA", AssetID: "A2", AssetAmount: 2500.00},
		Assets{AssetName: "JD", AssetID: "A3", AssetAmount: 3500.00},
	}

	asset5 := []Assets{
		Assets{AssetName: "AMD", AssetID: "A4", AssetAmount: 3500.00},
		Assets{AssetName: "INTEL", AssetID: "A5", AssetAmount: 5500.00},
	}

	var clients = []Clients{
		Clients{Name: "CITI ADMINISTRATION", ClientID: "AD000001", ClientType: "Admin", SafeKeepingAccount: "348912452", Currency: "USD", Asset: asset1, Status: "active"},
		Clients{Name: "THE GREENWALL FOUNDATION", ClientID: "20180001", ClientType: "Regular", SafeKeepingAccount: "348912975", Currency: "USD", Asset: asset2, Status: "active"},
		Clients{Name: "SOLAR CAPITAL LTD", ClientID: "20180002", ClientType: "Regular", SafeKeepingAccount: "348912325", Currency: "USD", Asset: asset3, Status: "active"},
		Clients{Name: "PFPC-DFA FUNDS-IRISH", ClientID: "20180003", ClientType: "Regular", SafeKeepingAccount: "348912345", Currency: "USD", Asset: asset4, Status: "active"},
		Clients{Name: "ORBIS GROUP", ClientID: "20180004", ClientType: "Regular", SafeKeepingAccount: "348912790", Currency: "USD", Asset: asset5, Status: "active"},
	}

	i := 0
	for i < len(clients) {
		fmt.Println("i is ", i)
		clientAsBytes, _ := json.Marshal(clients[i])
		APIstub.PutState("CLIENT"+strconv.Itoa(i), clientAsBytes)
		fmt.Println("Added", clients[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "CLIENT0"
	endKey := "CLIENT999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

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

	fmt.Printf("- queryAllAsset:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) doInnerTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Seller client key//Buyer //Asset ID/Amount
	if len(args) != 4 {
		//return shim.Error('Incorrect number of arguments:'+len(args)+'. Expecting 4:[BankIDFrom,BankIDTo,AssetName,AssetAmount]')
	}
	startKey := "TRANC1"
	endKey := "TRANC999"
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	count := 1
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if queryResponse.Value != nil {
			count = count + 1
		} else {
			break
		}
	}
	var transaction_id string
	transaction_id = "TRANC" + strconv.Itoa(count)

	bankFAsBytes, _ := APIstub.GetState(args[0])
	bankTAsBytes, _ := APIstub.GetState(args[1])

	bankF := Clients{}
	bankT := Clients{}
	json.Unmarshal(bankFAsBytes, &bankF)
	json.Unmarshal(bankTAsBytes, &bankT)
	var asset_amount_from float64
	asset_from := Assets{}
	var asset_exist_in_to = false
	amount, _ := strconv.ParseFloat(args[3], 64)
	for i, asset := range bankF.Asset {
		if asset.AssetID == args[2] && asset.AssetAmount >= amount {
			asset_amount_from = asset.AssetAmount
			asset_from = asset
			bankF.Asset[i].AssetAmount = asset.AssetAmount - amount
		}
	}
	if asset_amount_from == 0 || asset_amount_from < amount {
		return shim.Error("Incorrect number of arguments. Expecting Correct Asset")
	} else {
		for i, asset_to := range bankT.Asset {
			if asset_to.AssetID == args[2] {
				asset_from = asset_to
				bankT.Asset[i].AssetAmount = asset_to.AssetAmount + amount
				asset_exist_in_to = true
			}
		}
	}
	if !asset_exist_in_to {
		bankT.Asset = append(bankT.Asset, Assets{AssetName: asset_from.AssetName, AssetID: asset_from.AssetID, AssetAmount: amount})
	}

	// registre transaction
	bankFAsBytes_after, _ := json.Marshal(bankF)
	bankTAsBytes_after, _ := json.Marshal(bankT)
	transactionAsBytes, _ := json.Marshal(Transactions{TsID: transaction_id, Seller: args[0], Buyer: args[1], Asset: Assets{AssetName: asset_from.AssetName, AssetID: asset_from.AssetID, AssetAmount: amount}, TimeStamp: string(time.Now().Unix())})

	APIstub.PutState(args[0], bankFAsBytes_after)
	APIstub.PutState(args[1], bankTAsBytes_after)
	APIstub.PutState(transaction_id, transactionAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) queryTransactions(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID like :CLIENT0")
	}
	findall := false
	if args[0] == "ALL" {
		findall = true
	}
	startKey := "TRANC1"
	endKey := "TRANC999"
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	var allResults []Transactions = make([]Transactions, 5)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		objres := Transactions{}
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(queryResponse.Value, &objres)
		if findall || objres.Seller == args[0] || objres.Buyer == args[0] {
			allResults = append(allResults, objres)
		}
	}
	transactionsAsBytes, _ := json.Marshal(allResults)
	return shim.Success(transactionsAsBytes)
}


func (s *SmartContract) setSwiftMessage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
    var swiftmsg=SwiftMessage{MsgID : args[0], MsgContent:args[1], MsgStatus: args[2], MsgType:args[3], TimeStamp: string(time.Now().Unix())}
	swiftmsg_bytes, err := json.Marshal(swiftmsg)
	if err != nil {
			return shim.Error(err.Error())
		}
	err2 :=APIstub.PutState(args[0], swiftmsg_bytes)
	if err2 != nil {
			return shim.Error(err2.Error())
		}
	return shim.Success(nil)
}

func (s *SmartContract) getSwiftMessage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	swiftmsg_bytes, err := APIstub.GetState(args[0])
	if err != nil {
			return shim.Error(err.Error())
		}
	if swiftmsg_bytes == nil {
			return shim.Error("Message is null")
		}	
	return shim.Success(swiftmsg_bytes)
}

func (s *SmartContract) updateSwiftMessage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	swiftmsg_bytes, _ := APIstub.GetState(args[0])
	msg := SwiftMessage{}

	json.Unmarshal(swiftmsg_bytes, &msg)
	msg.MsgStatus = args[1]

	swiftmsg_bytes, _ = json.Marshal(msg)
	APIstub.PutState(args[0], swiftmsg_bytes)

	return shim.Success(nil)
}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

