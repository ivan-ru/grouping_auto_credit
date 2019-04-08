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
		AutoCreditDate     time.Time `json:"auto_credit_date"`
		Account            []account `json:"account"`
	}
)

const (
	numberOfAccount            = 1000000
	numberOfDebitSourceAccount = 100000
	maxJob                     = 10
)

var (
	timeElapsed            time.Duration
	timeElapsedSeconds     float64
	timeElapsedNanoSeconds int64
)

func main() {
	// get mock data of accounts sorted by account opening date
	ungroupedAccountList := getData(numberOfAccount, numberOfDebitSourceAccount)
	start := startCountProcessTime() // start time count to record process time after get mock data
	chunkedData, dataToStore := chunkList(ungroupedAccountList)
	json.Marshal(dataToStore)
	// dataToStoreByte, _ := json.Marshal(dataToStore)
	// fmt.Println(string(dataToStoreByte))

	// // write result to file
	// f, _ := os.Create("account_grouped.json")

	// defer f.Close()
	// w := bufio.NewWriter(f)
	// w.WriteString(string(dataToStoreByte))

	// w.Flush()

	// a := 0
	// b := 0
	// for _, val := range chunkedDataWrapper {
	// 	for _, z := range val {
	// 		fmt.Println(len(z))
	// 		a += len(z)
	// 		b++
	// 	}
	// }
	// fmt.Println(a / b)

	// fmt.Println("chunkedData==========")
	json.Marshal(chunkedData)
	// chunkedDataWrapperByte, _ := json.Marshal(chunkedData)
	// fmt.Println(string(chunkedDataWrapperByte))
	timeElapsed := endCountProcessTime(start)
	printProcessTime(timeElapsed, numberOfAccount)
}

// getData() mock data from database
// to get the actual data we will run these query ==>
// nb: change 'hari ini' to today's date

// select * from account_auto_credit
// where auto_credit_status='A'
// AND account_status=1
// and auto credit date = 'hari ini'
// order by account_opening_date asc
func getData(numberOfAccountData int, numberOfDebitSourceAccount int) (accountList []account) {
	start := startCountProcessTime()
	for i := 0; i < numberOfAccountData; i++ {
		newAccount := account{
			AccountNumber:      strconv.Itoa(i + 1),
			AccountOpeningDate: time.Now().Local().Add(time.Second * time.Duration(i)).Format("2006-01-02 15:04:05"),
		}
		newAccount.DebitSourceAccount = strconv.Itoa(randInt(1, numberOfDebitSourceAccount))
		// fmt.Print(" = ")
		// fmt.Println(newAccount.DebitSourceAccount)
		accountList = append(accountList, newAccount)
	}
	timeElapsed := endCountProcessTime(start)
	timeElapsedSeconds = timeElapsed.Seconds()
	timeElapsedNanoSeconds = timeElapsed.Nanoseconds()
	fmt.Printf("Total time elapsed for getting random data => %f seconds => %d nanoseconds\n", timeElapsedSeconds, timeElapsedNanoSeconds)
	return
}

func randInt(min int, max int) (randomNumber int) {
	rand.Seed(time.Now().UnixNano())
	randomNumber = min + rand.Intn(max-min)
	// fmt.Print(randomNumber)
	return
}

func startCountProcessTime() time.Time {
	return time.Now()
}

func endCountProcessTime(start time.Time) time.Duration {
	return time.Since(start)
}

func printProcessTime(timeElapsed time.Duration, numberOfAccount int) {
	timeElapsedSeconds = timeElapsed.Seconds()
	timeElapsedNanoSeconds = timeElapsed.Nanoseconds()
	fmt.Printf("Total time elapsed for %d grouped data => %f seconds => %d nanoseconds\n", numberOfAccount, timeElapsedSeconds, timeElapsedNanoSeconds)
}
