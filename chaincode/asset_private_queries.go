package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)



// RequestToBuyExists returns true when asset Price exists on shared collection so we dont redefine it
func (s *SmartContract) RequestToBuyExists(ctx contractapi.TransactionContextInterface, assetID string,sharedCollection string) (bool, error) {

	requestToBuyKey, err := ctx.GetStub().CreateCompositeKey(requestToBuyObjectType, []string{assetID})
	if err != nil {
		return false, fmt.Errorf("failed to create composite key: %v", err)
	}
	requestJSON, err := ctx.GetStub().GetPrivateData(sharedCollection, requestToBuyKey) // Get price from shared collection
	if err != nil {
		return false, fmt.Errorf("failed to read RequestToBuyObject: %v", err)
	}

	return requestJSON != nil, nil
}


// ReadRequestToBuy gets the buyer's identity and buyers MSP from the transfer request from collection
func (s *SmartContract) ReadRequestToBuy(ctx contractapi.TransactionContextInterface, assetID string, sharedCollection string) (*RequestToBuyObject, error) {
	//log.Printf("ReadRequestToBuy: collection %v, ID %v", assetCollection, assetID)
	// composite key for RequestToBuyObject of this asset
	transferAgreeKey, err := ctx.GetStub().CreateCompositeKey(requestToBuyObjectType, []string{assetID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}

	log.Printf("ReadRequestToBuy: collection %v, ID %v", sharedCollection, assetID)
	requestJSON, err := ctx.GetStub().GetPrivateData(sharedCollection, transferAgreeKey) // Get the state from world state
	if err != nil {
		return nil, fmt.Errorf("failed to read RequestToBuyObject: %v", err)
	}
	if requestJSON == nil {
		log.Printf("RequestToBuyObject for %v does not exist", assetID)
		return nil, nil
	}

	var assetBuyRequestObj RequestToBuyObject
	err = json.Unmarshal(requestJSON, &assetBuyRequestObj)
	if err != nil {
		return nil, err
	}

	request := &RequestToBuyObject{
		ID:      assetBuyRequestObj.ID,
		BuyerID: 	 assetBuyRequestObj.BuyerID,
		BuyerMSP:  assetBuyRequestObj.BuyerMSP,
	}
	return request, nil
}


// AssetPriceExists returns true when asset Price exists on shared collection so we dont redefine it
func (s *SmartContract) AssetPriceExists(ctx contractapi.TransactionContextInterface, assetID string,sharedCollection string) (bool, error) {

	assetPriceKey, err := ctx.GetStub().CreateCompositeKey(typeAssetForSale, []string{assetID})
	if err != nil {
		return false, fmt.Errorf("failed to create composite key: %v", err)
	}
	priceJSON, err := ctx.GetStub().GetPrivateData(sharedCollection, assetPriceKey) // Get price from shared collection
	if err != nil {
		return false, fmt.Errorf("failed to read RequestToBuyObject: %v", err)
	}



	return priceJSON != nil, nil
}

// ReadAssetPrice get the price that the seller put on shared collection.Should return an unmarshalled object with properties of asset_id,price,trade_id in GO ID,Price,TradeID
func (s *SmartContract) ReadAssetPrice(ctx contractapi.TransactionContextInterface, assetID string,sharedCollection string) (*assetPriceTransientInput, error) {
	//create the price key of asset
	assetPriceKey, err := ctx.GetStub().CreateCompositeKey(typeAssetForSale, []string{assetID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}
	log.Printf("Read price from : Collection %v, ID %v", sharedCollection, assetID)
	priceJSON, err := ctx.GetStub().GetPrivateData(sharedCollection, assetPriceKey) // Get price from shared collection
	if err != nil {
		return nil, fmt.Errorf("failed to read RequestToBuyObject: %v", err)
	}

	if priceJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", assetID)
	}
	type assetPriceTemp struct {
		ID       string `json:"asset_id"`
		Price 	 int 	`json:"price"`
		TradeID  string `json:"trade_id"`
		Salt  	 string `json:"salt"`
	}
	var assetPriceInput assetPriceTemp
	err = json.Unmarshal(priceJSON, &assetPriceInput)
	if err != nil {
		return nil, err
	}
	request := &assetPriceTransientInput{
		ID:      assetPriceInput.ID,
		Price: 	 assetPriceInput.Price,
		TradeID: assetPriceInput.TradeID,
		Salt: 	 assetPriceInput.Salt}
	return request, nil
}


/*=========================Phase 3 =========================================*/

func (s *SmartContract) GetAssetSalesPrice(ctx contractapi.TransactionContextInterface, assetID string) (string, error) {
	return getAssetPrice(ctx, assetID, typeAssetForSale)
}

// GetAssetBidPrice returns the bid price
func (s *SmartContract) GetAssetBidPrice(ctx contractapi.TransactionContextInterface, assetID string) (string, error) {
	return getAssetPrice(ctx, assetID, typeAssetBid)
}

// getAssetPrice gets the bid or ask price from caller's implicit private data collection
func getAssetPrice(ctx contractapi.TransactionContextInterface, assetID string, priceType string) (string, error) {
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return "",fmt.Errorf("TransferAsset cannot be performed: Error %v", err)
	}

	collection, err := buildCollectionName(ctx)
	if err != nil {
		return "", err
	}

	assetPriceKey, err := ctx.GetStub().CreateCompositeKey(priceType, []string{assetID})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	price, err := ctx.GetStub().GetPrivateData(collection, assetPriceKey)
	if err != nil {
		return "", fmt.Errorf("failed to read asset price from implicit private data collection: %v", err)
	}
	if price == nil {
		return "", fmt.Errorf("asset price does not exist: %s", assetID)
	}

	return string(price), nil
}





/*========================================END OF PHASE 3===================================*/
// =======Rich queries =========================================================================
// Two examples of rich queries are provided below (parameterized query and ad hoc query).
// Rich queries pass a query string to the state database.
// Rich queries are only supported by state database implementations
//  that support rich query (e.g. CouchDB).
// The query string is in the syntax of the underlying state database.
// With rich queries there is no guarantee that the result set hasn't changed between
//  endorsement time and commit time, aka 'phantom reads'.
// Therefore, rich queries should not be used in update transactions, unless the
// application handles the possibility of result set changes between endorsement and commit time.
// Rich queries can be used for point-in-time queries against a peer.
// ============================================================================================

// ===== Example: Parameterized rich query =================================================

// QueryAssetByOwner queries for assets based on assetType, owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (s *SmartContract) QueryAssetByOwner(ctx contractapi.TransactionContextInterface, assetType string, owner string) ([]*Asset, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"%v\",\"owner\":\"%v\"}}", assetType, owner)

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QueryAssets uses a query string to perform a query for assets.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryAssetByOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
func (s *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {

	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// getQueryResultForQueryString executes the passed in query string.
func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(assetCollection, queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*Asset{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset *Asset

		err = json.Unmarshal(response.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		results = append(results, asset)
	}
	return results, nil
}

