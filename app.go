package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

const (
	dataFile  = "XYJ_ShopTable.txt"
	separator = "\t" // 定义分隔符
)

type ShopTable struct {
	ctx context.Context
}

func (a *App) GetShopTable() [][]string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return nil
	}
	fmt.Println("当前工作目录:", dir)
	// 1. 读取 TXT 文件
	txtFile := dir + "\\" + dataFile
	file, err := os.Open(txtFile)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return nil
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
	return rows
}

// 写入 TXT 文件
func (a *App) WriteData(rows [][]string) error {
	os.Rename(dataFile, "old_"+dataFile)
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range rows {
		_, err := writer.WriteString(strings.Join(row, separator) + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
