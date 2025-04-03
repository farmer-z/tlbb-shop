package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"sort"
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

var TotalIndex = 0

type ShopItem struct {
	TotalIndex       int    `json:"totalIndex"`
	Index            int    `json:"index"`
	ItemId           int    `json:"itemId"`
	ItemName         string `json:"itemName"`
	ItemCount        int    `json:"itemCount"`
	ItemPrice        int    `json:"itemPrice"`
	ItemDiscount     int    `json:"itemDiscount"`
	ItemDisplayColor string `json:"itemDisplayColor"`
	ItemSpecialType  int    `json:"itemSpecialType"`
}

type ItemInfo struct {
	ItemId   int    `json:"itemId"`
	ItemName string `json:"itemName"`
}

func init() {
	dir, _ := os.Getwd()

	txtFile := dir + "\\" + commonItemFile
	commonItemInfo, _ := readFile(txtFile)

	for _, row := range commonItemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[6]
	}
	txtFile = dir + "\\" + gemInfoFile
	gemInfo, _ := readFile(txtFile)

	for _, row := range gemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[7]
	}
	txtFile = dir + "\\" + equipBaseFile
	equipBaseInfo, _ := readFile(txtFile)
	for _, row := range equipBaseInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[10]
	}
}

func (a *App) GetShopTable() ([]int, error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 1. 读取 TXT 文件
	txtFile := dir + "\\" + shopTableFile
	rows, err := readFile(txtFile)
	if err != nil {
		return nil, err
	}
	TableData = make(map[int]map[int]map[int][]ShopItem)
	a.processData(rows)
	shopIds := make([]int, 0)
	for key, _ := range TableData {
		shopIds = append(shopIds, key)
	}
	sort.Ints(shopIds)
	return shopIds, nil
}
func (a *App) GetMenuItems(shopId int) ([]int, error) {
	if _, ok := TableData[shopId]; !ok {
		return nil, errors.New("不存在的商店")
	}
	menuIds := make([]int, 0)
	for key, _ := range TableData[shopId] {
		menuIds = append(menuIds, key)
	}
	sort.Ints(menuIds)
	return menuIds, nil
}
func (a *App) GetSubMenus(shopId int, menuId int) ([]int, error) {
	if _, ok := TableData[shopId]; !ok {
		return nil, errors.New("不存在的商店")
	}
	if _, ok := TableData[shopId][menuId]; !ok {
		return nil, errors.New("不存在的菜单")
	}
	subMenuIds := make([]int, 0)
	for key, _ := range TableData[shopId][menuId] {
		subMenuIds = append(subMenuIds, key)
	}
	sort.Ints(subMenuIds)
	return subMenuIds, nil
}
func (a *App) GetShopItems(shopId int, menuId int, subMenuId int) ([]ShopItem, error) {
	if _, ok := TableData[shopId]; !ok {
		return nil, errors.New("不存在的商店")
	}
	if _, ok := TableData[shopId][menuId]; !ok {
		return nil, errors.New("不存在的菜单")
	}
	if _, ok := TableData[shopId][menuId][subMenuId]; !ok {
		return nil, errors.New("不存在的子菜单")
	}
	result := TableData[shopId][menuId][subMenuId]

	remove := make(map[int]ShopItem)
	for _, item := range result {
		remove[item.Index] = item
	}
	newResult := make([]ShopItem, 0)
	for _, item := range remove {
		newResult = append(newResult, item)
	}
	sort.Slice(newResult, func(i, j int) bool {
		return newResult[i].Index < newResult[j].Index
	})
	return TableData[shopId][menuId][subMenuId], nil
}

func readFile(path string) ([][]string, error) {

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

	startRow := rows[2:]
	for _, row := range startRow {
		totalIndex, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		if totalIndex > TotalIndex {
			TotalIndex = totalIndex
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
			TableData[shopId][menuId][subMenuId] = make([]ShopItem, 0)
		}
		TableData[shopId][menuId][subMenuId] = append(TableData[shopId][menuId][subMenuId], ShopItem{
			TotalIndex:       totalIndex,
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

func (a *App) SearchItem(name string) ([]ItemInfo, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return []ItemInfo{}, nil
	}
	result := make([]ItemInfo, 0)
	fmt.Println("itemInfoMap:", ItemInfoMap)
	for itemId, itemName := range ItemInfoMap {
		if strings.Contains(itemName, name) {
			result = append(result, ItemInfo{
				ItemId:   itemId,
				ItemName: itemName,
			})
		}
	}
	return result, nil
}

func (a *App) MatchItemName(itemIdStr string) (string, error) {
	itemIdStr = strings.TrimSpace(itemIdStr)
	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return "", errors.New("itemIdStr is not int")
	}
	return ItemInfoMap[itemId], nil
}

// 新增文件追加方法
func (a *App) appendToFile(newData []string) error {
	// 使用管道符分隔，与原始格式一致
	line := strings.Join(newData, separator) + "\n"
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	shopTablePath := dir + "\\" + shopTableFile
	//先复制一份文件作为备份
	backupFile := dir + "\\" + shopTableFile + ".bak"
	err = a.CopyFile(shopTablePath, backupFile)
	if err != nil {
		return err
	}

	// 以追加模式打开文件
	f, err := os.OpenFile(shopTablePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 使用GBK编码写入
	encoder := simplifiedchinese.GBK.NewEncoder()
	writer := transform.NewWriter(f, encoder)

	_, err = writer.Write([]byte(line))
	if err != nil {
		return fmt.Errorf("写入失败: %v", err)
	}
	return nil
}

func (a *App) CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}
	return nil
}

// AddShopItem 修改AddShopItem方法
func (a *App) AddShopItem(shopId, menuId, subMenuId, itemId, itemCount, itemPrice int) error {
	fmt.Println("进入方法AddShopItem", shopId, menuId, subMenuId, itemId, itemCount, itemPrice)
	if _, ok := TableData[shopId]; !ok {
		return errors.New("不存在的商店")
	}
	if _, ok := TableData[shopId][menuId]; !ok {
		return errors.New("不存在的菜单")
	}
	if _, ok := TableData[shopId][menuId][subMenuId]; !ok {
		return errors.New("不存在的子菜单")
	}
	// 生成新行数据（根据XYJ_ShopTable.txt的列顺序）
	TotalIndex++
	shopTableData := TableData[shopId][menuId][subMenuId]
	index := 0
	for _, data := range shopTableData {
		if data.Index > index {
			index = data.Index
		}
	}
	newLine := []string{
		strconv.Itoa(TotalIndex),
		strconv.Itoa(shopId),
		strconv.Itoa(menuId),
		strconv.Itoa(subMenuId),
		strconv.Itoa(index + 1),
		strconv.Itoa(itemId),
		strconv.Itoa(itemCount),
		strconv.Itoa(itemPrice),
		strconv.Itoa(100),
		"",
		strconv.Itoa(0),
	}

	// 先写入文件
	if err := a.appendToFile(newLine); err != nil {
		return err
	}
	return nil
}

func (a *App) ReloadItemMap() {
	ItemInfoMap = make(map[int]string)
	dir, _ := os.Getwd()

	txtFile := dir + "\\" + commonItemFile
	commonItemInfo, _ := readFile(txtFile)

	for _, row := range commonItemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[6]
	}
	txtFile = dir + "\\" + gemInfoFile
	gemInfo, _ := readFile(txtFile)

	for _, row := range gemInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[7]
	}
	txtFile = dir + "\\" + equipBaseFile
	equipBaseInfo, _ := readFile(txtFile)
	for _, row := range equipBaseInfo {
		itemId, err := strconv.Atoi(row[0])
		if err != nil {
			// 非int类型
			continue
		}
		ItemInfoMap[itemId] = row[10]
	}
}
