package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testNumberOfAccounts            = 100
	testNumberOfDebitSourceAccounts = 10
)

func Test_chunkList(t *testing.T) {
	assert := assert.New(t)
	ungroupedAccountList := getData(testNumberOfAccounts, testNumberOfDebitSourceAccounts)
	start := startCountProcessTime()
	chunkedData, dataToStore := chunkList(ungroupedAccountList)
	json.Marshal(dataToStore)
	json.Marshal(chunkedData)
	timeElapsed := endCountProcessTime(start)
	printProcessTime(timeElapsed, testNumberOfAccounts)

	// chunkedData================
	numberOfAccountsChunked := 0
	lastOpeningDate := time.Time{}
	lastDebitSourceAccount := ""
	layout := "2006-01-02 15:04:05"
	for _, a := range chunkedData {
		for _, b := range a {
			numberOfAccountsChunked += len(b)
			for _, c := range b {
				timeSortedProperly := true
				accountOpeningDateInTime, _ := time.Parse(layout, c.AccountOpeningDate)
				if !lastOpeningDate.Before(accountOpeningDateInTime) {
					timeSortedProperly = false
				}
				if lastDebitSourceAccount != c.DebitSourceAccount {
					timeSortedProperly = true
				}
				// fmt.Println("========================")
				// fmt.Println("lastOpeningDate")
				// fmt.Println(lastOpeningDate)
				// fmt.Println("accountOpeningDateInTime")
				// fmt.Println(accountOpeningDateInTime)
				// fmt.Println("timeSortedProperly")
				// fmt.Println(timeSortedProperly)
				// fmt.Println("========================")
				assert.Equal(true, timeSortedProperly, "lastOpeningDate should be less than accountOpeningDateInTime")
				lastOpeningDate = accountOpeningDateInTime
				lastDebitSourceAccount = c.DebitSourceAccount
			}
		}
	}
	// chunkedData================

	// dataToStore================
	numberOfAccountsToStore := 0
	for _, a := range dataToStore {
		numberOfAccountsToStore += len(a.Account)
	}
	// dataToStore================

	assert.Equal(numberOfAccountsToStore, testNumberOfAccounts, "account list length on dataToStore should be equal")
	assert.Equal(numberOfAccountsChunked, testNumberOfAccounts, "account list length on chunkedData should be equal")
}
