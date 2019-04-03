package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type (
	account struct {
		AccountNumber      string `json:"account_number"`
		DebitSourceAccount string `json:"debit_source_account"`
		AccountOpeningDate string `json:"account_opening_date"`
	}

	accountGroup struct {
		DebitSourceAccount string    `json:"debit_source_account"`
		Account            []account `json:"account"`
	}
)

const (
	numberOfAccount            = 1000000
	numberOfDebitSourceAccount = 1000
)

var (
	timeElapsed            time.Duration
	timeElapsedSeconds     float64
	timeElapsedNanoSeconds int64
	dataToStore            = []accountGroup{}
	accountListGrouped     = make(map[string][]account)
)

func main() {
	// get mock data of accounts sorted by account opening date
	ungroupedAccountList := getData()
	for _, val := range ungroupedAccountList {
		accountListGrouped[val.DebitSourceAccount] = append(accountListGrouped[val.DebitSourceAccount], val)
	}

	start := time.Now() // start time count to record process time after get mock data
	for key, val := range accountListGrouped {
		newAccountGroup := accountGroup{
			DebitSourceAccount: key,
			Account:            val,
		}
		dataToStore = append(dataToStore, newAccountGroup)
	}
	json.Marshal(dataToStore)
	// dataToStoreByte, _ := json.Marshal(dataToStore)
	// fmt.Println(string(dataToStoreByte))

	timeElapsed = time.Since(start)
	timeElapsedSeconds = timeElapsed.Seconds()
	timeElapsedNanoSeconds = timeElapsed.Nanoseconds()
	fmt.Printf("Total time elapsed for %d grouped data => %f seconds => %d nanoseconds\n", numberOfAccount, timeElapsedSeconds, timeElapsedNanoSeconds)

	// // write result to file
	// f, _ := os.Create("account_grouped.json")

	// defer f.Close()
	// w := bufio.NewWriter(f)
	// w.WriteString(string(dataToStoreByte))

	// w.Flush()

}

// getData() mock data from database
// to get the actual data we will run these query ==>
// nb: change 'hari ini' to today's date

// select * from account_auto_credit
// where auto_credit_status='A'
// AND account_status=1
// and auto credit date = 'hari ini'
// order by account_opening_date asc
func getData() (accountList []account) {
	for i := 0; i < numberOfAccount; i++ {
		newAccount := account{
			AccountNumber:      strconv.Itoa(i + 1),
			AccountOpeningDate: time.Now().Local().Add(time.Second * time.Duration(i)).Format("2006-01-02 15:04:05"),
		}
		newAccount.DebitSourceAccount = strconv.Itoa(randInt(1, numberOfDebitSourceAccount))
		accountList = append(accountList, newAccount)
		// fmt.Println(accountList[i])
	}
	return
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
