package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	ID     string
	Fee    int
	Weight int
	Parents []string
}

func main() {
	transactions := getTransactions()
	fmt.Println(transactions)
}

func getTransactions() []*Transaction {
	filePath := "mempool.csv"
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer readFile.Close()

	reader := csv.NewReader(readFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil
	}
	
	var mainData []*Transaction

	for _, record := range records {
		
		row:=strings.Split(record[0], ",")
		
		fee, _ := strconv.Atoi(row[1])
		weight, _ := strconv.Atoi(row[2])
		var parents []string
		if row[3] != "" {
			parents = strings.Split(row[3], ",")
		}
		transaction := &Transaction{
			ID:     row[0],
			Fee:    fee,
			Weight: weight,
			Parents: parents,
		}
		
		mainData = append(mainData, transaction)
	}
	
	sort.Slice(mainData, func(i, j int) bool {
		return mainData[i].Fee > mainData[j].Fee
	})

	sortedByParent:= sortTrxByParents(mainData)

	filteredTransactions:= getMaxBlockweightByFee(sortedByParent)
	return filteredTransactions
	
}

func getMaxBlockweightByFee(transactions []*Transaction) []*Transaction {
	sumOfWeight := 0
	var filteredByweight []*Transaction
	for _, transaction := range transactions {
		if sumOfWeight < 4000000 {
			filteredByweight = append(filteredByweight, transaction)
			sumOfWeight += transaction.Weight
		}
	}
	return filteredByweight
}

func sortTrxByParents(transactions []*Transaction) []*Transaction {
	transactionMap := make(map[string]*Transaction)
	for _, trx := range transactions {
		transactionMap[trx.ID] = trx
	}
	
	
	var resArr []*Transaction
	for _, trx := range transactions {
		
		if len(trx.Parents) > 0 {
			
			for _, parent := range trx.Parents {
				if parentTrx, ok := transactionMap[parent]; ok {
					fmt.Println(parentTrx)
					
					if ok  {
						//check if the parent already exists on resArr
						resArr = append(resArr, parentTrx)
						
						resArr = append(resArr, trx)
						break
					}
				}
			}
		} else {
			resArr = append(resArr, trx)
		}
	}
	
	 return resArr
}


