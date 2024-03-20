package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// 调查房子价值 1 2 3 可以并行，4 5 依赖前面的结果
// 1. 找物业，查到住户可以是多人 Owners
// 2. 评价房产&淘宝法拍 Price
// 3. 查银行账户 BankAccount，是否有贷款
// 4. 查住户可以拿到的钱 GetPrice
// 5. 存入数据库

// 房屋信息
type HouseInfo struct {
	ID          int      // 房子ID
	Owners      []string // 住户
	Price       int      // 价格
	BankAccount []int    // 银行账户
	GetPrice    int      // 售出可得，可以是负数
}

type Response struct {
	data map[string]any
	err  error
}

const (
	Owners      = "Owners"
	Price       = "Price"
	BankAccount = "BankAccount"
)

func sellHouseInfo(id int) (*HouseInfo, error) {
	resChan := make(chan Response, 3)
	wg := &sync.WaitGroup{}

	wg.Add(3)

	go tenementInfo(id, resChan, wg)
	go evaluateHouse(id, resChan, wg)
	go loanInfo(id, resChan, wg)

	wg.Wait()

	close(resChan)

	houseInfo := &HouseInfo{}

	resMap := make(map[string]any, 3)

	for res := range resChan {
		if res.err != nil {
			return nil, res.err
		}
		// 将 data 中的数据提取到 resMap 中
		for key, value := range res.data {
			resMap[key] = value
		}
	}

	// 将 resMap 中的数据提取到 houseInfo 中
	// 因为 any 类型的数据不能直接赋值给其他类型，所以需要类型断言

	if _, ok := resMap[Owners]; !ok {
		return nil, fmt.Errorf("没有找到住户信息")
	}
	houseInfo.Owners = resMap[Owners].([]string)

	if _, ok := resMap[Price]; !ok {
		return nil, fmt.Errorf("没有找到房屋价格")
	}
	houseInfo.Price = resMap[Price].(int)

	if _, ok := resMap[BankAccount]; !ok {
		return nil, fmt.Errorf("没有找到银行账户")
	}
	houseInfo.BankAccount = resMap[BankAccount].([]int)

	return houseInfo, nil
}

func tenementInfo(id int, resChan chan Response, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Millisecond * 2000)

	users := []string{"张三", "李四", "王五"}

	resChan <- Response{
		map[string]any{
			Owners: users,
		},
		nil,
	}
}

func evaluateHouse(id int, resChan chan Response, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Millisecond * 2000)

	resChan <- Response{
		map[string]any{
			Price: 1000000,
		},
		nil,
	}
}

func loanInfo(id int, resChan chan Response, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Millisecond * 2000)

	cardIDs := []int{123456, 654321}

	resChan <- Response{
		map[string]any{
			BankAccount: cardIDs,
		},
		nil,
	}
}

func unPayLoan(cards []int) int {
	return len(cards) * 200000
}

func storeData(house *HouseInfo) error {
	fmt.Printf("%+v\n", house)
	fmt.Println("存入数据库")
	return nil
}

func main() {
	start := time.Now()
	houseInfo, err := sellHouseInfo(14234)
	if err != nil {
		log.Fatal(err)
	}
	houseInfo.GetPrice = houseInfo.Price - unPayLoan(houseInfo.BankAccount)
	if err := storeData(houseInfo); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("耗时：%v\n", time.Since(start))
}
