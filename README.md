# This is the 3rd scenario where 3 orgs have a  shared collection betwenn org1-org2  and org2-3. Data like price and buy request are sent through chaincode.
 


I used attribute based access control for creation of product and create a simple example of using that.
The roles are farming and retailer and it's easier if every organization has different roles. 

			Org1 => farming
			Org2 => retailing 
			Org3 => supermarket


In this phase i'm going to add private data collections for Org1,Org2 and Org3
It is going to be consisted of:

			* Implicit Collection for Org1
			* Implicit Collection for Org2
			* Implicit Collection for Org3
			* Priv collection between orgs 12 and 23
			*No data are sent through external means,only through chaincode as private data
			* Buy request has also buyersMSP , and price goes into implicit and the shared collection
			for Buyer to be able to read it and accept that price

	 
INITIALIZATION HAS TO FOLLOW THESE STEPS

If you decide to have separate file for main.go and the chaincode folder you have to go on the 
project folder ex phase3 where you have folder chaincode and file main.go 
and enter command *go mod init phase3* . The main file has to look like this

		package main

		import (
			"log"

			"github.com/hyperledger/fabric-contract-api-go/contractapi"
			"phase3/chaincode"
			)

		func main() {
			assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
			if err != nil {
				log.Panicf("Error creating asset-transfer-private-data chaincode: %v", err)
			}

			if err := assetChaincode.Start(); err != nil {
				log.Panicf("Error starting asset-transfer-private-data chaincode: %v", err)
			}
		}
				
				


