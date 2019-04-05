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
	numberOfAccount            = 30
	numberOfDebitSourceAccount = 10
)

var (
	timeElapsed                                time.Duration
	timeElapsedSeconds                         float64
	timeElapsedNanoSeconds                     int64
	maxJob                                     = 10
	accountListGrouped                         = make(map[string][]account)
	dataToStore                                = []accountGroup{}
	chunkedDataTempPerDebitSourceAccount       = [][]account{}
	chunkedDataTempPerDebitSourceAccountLength = 0
	chunkedDataTemp                            = [][][]account{}
	chunkedDataWrapper                         = [][][]account{}
	loopCount                                  int
	lastLoop                                   bool
)

func main() {
	// get mock data of accounts sorted by account opening date
	ungroupedAccountList := getData()
	start := time.Now() // start time count to record process time after get mock data
	for _, val := range ungroupedAccountList {
		accountListGrouped[val.DebitSourceAccount] = append(accountListGrouped[val.DebitSourceAccount], val)
	}

	for key, val := range accountListGrouped {
		if loopCount == len(accountListGrouped)-1 {
			lastLoop = true
		}
		newAccountGroup := accountGroup{
			DebitSourceAccount: key,
			AutoCreditDate:     time.Now(),
			Account:            val,
		}
		dataToStore = append(dataToStore, newAccountGroup) //data to store to db later as history

		fmt.Println("(len(val) + chunkedDataTempPerDebitSourceAccountLength)")
		fmt.Println((len(val) + chunkedDataTempPerDebitSourceAccountLength))
		if ((len(val) + chunkedDataTempPerDebitSourceAccountLength) > maxJob) || lastLoop {
			if lastLoop {
				chunkedDataTempPerDebitSourceAccount = append(chunkedDataTempPerDebitSourceAccount, val)
			}
			if chunkedDataTempPerDebitSourceAccountLength == maxJob {
				chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTempPerDebitSourceAccount)
				fmt.Println("chunkedDataWrapperrrrrrrrrrrrrrr")
				fmt.Println(chunkedDataWrapper)
			} else {
				fmt.Println("asdasdasd")
				if len(chunkedDataTemp) == 0 || (len(val)+chunkedDataTempPerDebitSourceAccountLength) != maxJob {
					chunkedDataTemp = append(chunkedDataTemp, chunkedDataTempPerDebitSourceAccount)
				} else {
					for key, val := range chunkedDataTemp {
						if len(val) == maxJob {
							chunkedDataWrapper = append(chunkedDataWrapper, val)
							continue
						}
						if (len(val) + chunkedDataTempPerDebitSourceAccountLength) == maxJob {
							chunkedDataTemp[key] = append(chunkedDataTemp[key], chunkedDataTempPerDebitSourceAccount...)
							fmt.Println("chunkedDataTemp[key]")
							fmt.Println(chunkedDataTemp[key])
							var accountGroupPerDebitSourceTemp [][]account
							accountGroupPerDebitSourceTemp = chunkedDataTemp[key]
							chunkedDataWrapper = append(chunkedDataWrapper, accountGroupPerDebitSourceTemp)
							// chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTemp[key])
							chunkedDataTemp[key] = nil
							break
						}
					}
				}
			}
			chunkedDataTempPerDebitSourceAccount = nil
			chunkedDataTempPerDebitSourceAccountLength = 0
		}
		if lastLoop {
			fmt.Printf("length : %d\n", len(val))
			fmt.Println("chunkedDataTemp")
			fmt.Println(chunkedDataTemp)
			chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTemp...)
			break
		}

		chunkedDataTempPerDebitSourceAccount = append(chunkedDataTempPerDebitSourceAccount, val) //chunked data to process auto credit
		chunkedDataTempPerDebitSourceAccountLength += len(val)
		loopCount++
		fmt.Printf("length : %d\n", len(val))
		fmt.Println("chunkedDataTempPerDebitSourceAccount")
		fmt.Println(chunkedDataTempPerDebitSourceAccount)
	}
	json.Marshal(dataToStore)
	// dataToStoreByte, _ := json.Marshal(dataToStore)
	// fmt.Println(string(dataToStoreByte))

	countProcessTime(start)

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

	fmt.Println("chunkedDataWrapper==========")
	chunkedDataWrapperByte, _ := json.Marshal(chunkedDataWrapper)
	fmt.Println(string(chunkedDataWrapperByte))
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
		// fmt.Print(" = ")
		// fmt.Println(newAccount.DebitSourceAccount)
		accountList = append(accountList, newAccount)
	}
	return
}

func randInt(min int, max int) (randomNumber int) {
	rand.Seed(time.Now().UnixNano())
	randomNumber = min + rand.Intn(max-min)
	// fmt.Print(randomNumber)
	return
}

func countProcessTime(start time.Time) {
	timeElapsed = time.Since(start)
	timeElapsedSeconds = timeElapsed.Seconds()
	timeElapsedNanoSeconds = timeElapsed.Nanoseconds()
	fmt.Printf("Total time elapsed for %d grouped data => %f seconds => %d nanoseconds\n", numberOfAccount, timeElapsedSeconds, timeElapsedNanoSeconds)
}
