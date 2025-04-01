package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"strconv"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {

	fmt.Println("App实例已创建")
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

const (
	shopTableFile  = "XYJ_ShopTable.txt"
	separator      = "\t" // 定义分隔符
	commonItemFile = "CommonItem.txt"
	gemInfoFile    = "GemInfo.txt"
	equipBaseFile  = "EquipBase.txt"
)

var TableData = make(map[int]map[int]map[int][]ShopItem)

var ItemInfoMap = make(map[int]string)

type ShopItem struct {
	Index            int    `json:"index"`
	ItemId           int    `json:"itemId"`
	ItemName         string `json:"itemName"`
	ItemCount        int    `json:"itemCount"`
	ItemPrice        int    `json:"itemPrice"`
	ItemDiscount     int    `json:"itemDiscount"`
	ItemDisplayColor string `json:"itemDisplayColor"`
	ItemSpecialType  int    `json:"itemSpecialType"`
}

func (a *App) GetShopTable() (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// 1. 读取 TXT 文件
	txtFile := dir + "\\" + shopTableFile
	rows, err := a.readFile(txtFile)
	if err != nil {
		return "", err
	}
	a.processData(rows)
	// 序列化为 JSON 字符串
	jsonData, err := json.Marshal(TableData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (a *App) readFile(path string) ([][]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(path + "不存在")
	}
	defer file.Close()
	// 显式使用 GBK 解码器
	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(file, decoder)
	// 2. 逐行解析 TXT
	var rows [][]string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// 假设字段用逗号分隔（可替换为其他分隔符，如\t）
		fields := strings.Split(line, separator)
		rows = append(rows, fields)
	}

	return rows, nil
}

func (a *App) processData(rows [][]string) {
	dir, _ := os.Getwd()

	txtFile := dir + "\\" + commonItemFile
	commonItemInfo, _ := a.readFile(txtFile)

	for _, row := range commonItemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[6]
	}
	txtFile = dir + "\\" + gemInfoFile
	gemInfo, _ := a.readFile(txtFile)

	for _, row := range gemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[7]
	}
	txtFile = dir + "\\" + equipBaseFile
	equipBaseInfo, _ := a.readFile(txtFile)
	for _, row := range equipBaseInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[10]
	}
	startRow := rows[2:]
	for _, row := range startRow {
		_, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		shopId, _ := strconv.Atoi(row[1])
		menuId, _ := strconv.Atoi(row[2])
		subMenuId, _ := strconv.Atoi(row[3])
		itemIndex, _ := strconv.Atoi(row[4])
		ItemId, _ := strconv.Atoi(row[5])
		ItemCount, _ := strconv.Atoi(row[6])
		ItemPrice, _ := strconv.Atoi(row[7])
		ItemDiscount, _ := strconv.Atoi(row[8])
		ItemSpecialType, _ := strconv.Atoi(row[10])
		if _, ok := TableData[shopId]; !ok {
			TableData[shopId] = make(map[int]map[int][]ShopItem)
		}
		if _, ok := TableData[shopId][menuId]; !ok {
			TableData[shopId][menuId] = make(map[int][]ShopItem)
		}
		if _, ok := TableData[shopId][menuId][subMenuId]; !ok {
			TableData[shopId][menuId][subMenuId] = append(TableData[shopId][menuId][subMenuId], ShopItem{
				Index:            itemIndex,
				ItemId:           ItemId,
				ItemName:         ItemInfoMap[ItemId],
				ItemCount:        ItemCount,
				ItemPrice:        ItemPrice,
				ItemDiscount:     ItemDiscount,
				ItemDisplayColor: row[9],
				ItemSpecialType:  ItemSpecialType,
			})
		}

	}

}
