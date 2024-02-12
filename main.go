package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Transaction struct {
	ID     string
	Fee    int
	Weight int
	Parent string
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
		fee, _ := strconv.Atoi(record[1])
		weight, _ := strconv.Atoi(record[2])
		transaction := &Transaction{
			ID:     record[0],
			Fee:    fee,
			Weight: weight,
			Parent: record[3],
		}
		mainData = append(mainData, transaction)
	}

	sort.Slice(mainData, func(i, j int) bool {
		return mainData[i].Fee > mainData[j].Fee
	})

	return getMaxBlockweightByFee(mainData)
}

func getMaxBlockweightByFee(transactions []*Transaction) []*Transaction {
	sumOfWeight := 0
	var cherryPicked []*Transaction
	for _, transaction := range transactions {
		if sumOfWeight < 4000000 {
			cherryPicked = append(cherryPicked, transaction)
			sumOfWeight += transaction.Weight
		}
	}
	return cherryPicked
}

func sortTrxByParents(transactions []*Transaction) []*Transaction {
	transactionMap := make(map[string]*Transaction)
	for _, trx := range transactions {
		transactionMap[trx.ID] = trx
	}

	var resArr []*Transaction
	for _, trx := range transactions {
		if trx.Parent != "" {
			parents := parseParents(trx.Parent)
			for _, parent := range parents {
				if parentTrx, ok := transactionMap[parent]; ok {
					if parentTrx != nil {
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

func parseParents(parentString string) ([]string, error) {
	return csv.NewReader(os.Stdin).Read(strings.NewReader(parentString)),nil
}
