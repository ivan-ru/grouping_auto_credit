package main

import (
	"time"
)

var (
	chunkedDataTempPerDebitSourceAccount       = [][]account{}
	chunkedDataTempPerDebitSourceAccountLength = 0
	chunkedDataTemp                            = [][][]account{}
	chunkedDataWrapper                         = [][][]account{}
	loopCount                                  int
	lastLoop                                   bool
)

// chunkList ...
func chunkList(ungroupedAccountList []account) (chunkedDataWrapper [][][]account, dataToStore []accountGroup) {
	accountListGrouped := make(map[string][]account)
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

		// fmt.Println("(len(val) + chunkedDataTempPerDebitSourceAccountLength)")
		// fmt.Println((len(val) + chunkedDataTempPerDebitSourceAccountLength))
		if ((len(val) + chunkedDataTempPerDebitSourceAccountLength) > maxJob) || lastLoop {
			if lastLoop {
				chunkedDataTempPerDebitSourceAccount = append(chunkedDataTempPerDebitSourceAccount, val)
			}
			if chunkedDataTempPerDebitSourceAccountLength == maxJob {
				chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTempPerDebitSourceAccount)
				// fmt.Println("chunkedDataWrapperrrrrrrrrrrrrrr")
				// fmt.Println(chunkedDataWrapper)
			} else {
				// fmt.Println("asdasdasd")
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
							// fmt.Println("chunkedDataTemp[key]")
							// fmt.Println(chunkedDataTemp[key])
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
			// fmt.Printf("length : %d\n", len(val))
			// fmt.Println("chunkedDataTemp")
			// fmt.Println(chunkedDataTemp)
			chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTemp...)
			break
		}

		chunkedDataTempPerDebitSourceAccount = append(chunkedDataTempPerDebitSourceAccount, val) //chunked data to process auto credit
		chunkedDataTempPerDebitSourceAccountLength += len(val)
		loopCount++
		// fmt.Printf("length : %d\n", len(val))
		// fmt.Println("chunkedDataTempPerDebitSourceAccount")
		// fmt.Println(chunkedDataTempPerDebitSourceAccount)
	}
	return
}
