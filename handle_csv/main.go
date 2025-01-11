package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// checkURL 检查URL是否可以访问。
func checkURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// processURLs 处理一批URL，并将结果发送到通道。
func processURLs(urls []string, wg *sync.WaitGroup, goodUrls, badUrls chan string) {
	defer wg.Done() 
	for _, url := range urls {
		if checkURL(url) {
			goodUrls <- url 
		} else {
			badUrls <- url 
		}
	}
}

// writeCSV 将数据写入CSV文件。
func writeCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// 记录开始时间
	startTime := time.Now()

	// 读取CSV文件
	data, err := ioutil.ReadFile("./data.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(strings.NewReader(string(data)))
	urls, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 初始化通道
	goodUrls := make(chan string, len(urls))
	badUrls := make(chan string, len(urls))

	// 初始化工作池
	var wg sync.WaitGroup
	chunks := 10 
	numWorkers := len(urls) / chunks
	if len(urls)%chunks > 0 {
		numWorkers++
	}

	
	for i := 0; i < numWorkers; i++ {
		start := i * chunks
		end := (i + 1) * chunks
		if end > len(urls) {
			end = len(urls)
		}
		wg.Add(1)
		go processURLs(flattenURLs(urls[start:end]), &wg, goodUrls, badUrls)
	}

	// 等待所有工作goroutine完成并关闭通道
	wg.Wait()
	close(goodUrls)
	close(badUrls)

	// 收集结果
	var good, bad [][]string
	for url := range goodUrls {
		good = append(good, []string{url})
	}
	for url := range badUrls {
		bad = append(bad, []string{url})
	}

	if err := writeCSV("good.csv", good); err != nil {
		log.Fatal(err)
	}
	if err := writeCSV("bad.csv", bad); err != nil {
		log.Fatal(err)
	}

	
	endTime := time.Now()

	
	fmt.Printf("处理完成。总共耗时：%v\n", endTime.Sub(startTime))
}


func flattenURLs(urls [][]string) []string {
	var result []string
	for _, sublist := range urls {
		result = append(result, sublist[0]) 
	}
	return result
}