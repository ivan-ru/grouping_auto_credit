package main

import (
	"time"
)

var (
	chunkedDataTempPerDebitSourceAccount       = [][]autoCreditHistory{}
	chunkedDataTempPerDebitSourceAccountLength = 0
	chunkedDataTemp                            = [][][]autoCreditHistory{}
	chunkedDataWrapper                         = [][][]autoCreditHistory{}
	loopCount                                  int
	lastLoop                                   bool
)

// chunkList ...
func chunkList(ungroupedAccountList []autoCreditHistory) (chunkedDataWrapper [][][]autoCreditHistory, dataToStore []accountGroup) {
	accountListGrouped := make(map[string][]autoCreditHistory)
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
				if len(val)+len(chunkedDataTempPerDebitSourceAccount) < maxJob {
					chunkedDataTempPerDebitSourceAccount = append(chunkedDataTempPerDebitSourceAccount, val)
				} else {
					chunkedDataTemp = append(chunkedDataTemp, [][]autoCreditHistory{val})
				}
			}
			if chunkedDataTempPerDebitSourceAccountLength == maxJob {
				chunkedDataWrapper = append(chunkedDataWrapper, chunkedDataTempPerDebitSourceAccount)
				// fmt.Println("chunkedDataWrapperrrrrrrrrrrrrrr")
				// fmt.Println(chunkedDataWrapper)
			} else {
				// fmt.Println("asdasdasd")
				if len(chunkedDataTemp) == 0 || (len(val)+chunkedDataTempPerDebitSourceAccountLength) != maxJob {
					if len(chunkedDataTempPerDebitSourceAccount) != 0 {
						chunkedDataTemp = append(chunkedDataTemp, chunkedDataTempPerDebitSourceAccount)
					}
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
							var accountGroupPerDebitSourceTemp [][]autoCreditHistory
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
			// chunkedDataTempWrapperByte, _ := json.Marshal(chunkedDataTemp)
			// fmt.Println(string(chunkedDataTempWrapperByte))
			for _, val := range chunkedDataTemp {
				chunkedDataWrapper = append(chunkedDataWrapper, val)
			}
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
