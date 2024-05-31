package africa

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"sync"
	"testing"
)

func Test(t *testing.T) {
	var wg sync.WaitGroup
	// 启动并发请求
	wg.Add(3)
	// 创建一个channel来传递信息
	results := make(chan model.Result, 3)
	// 启动并发请求
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := Showmax(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := DSTV(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := BeinConnect(req)
		results <- res
	}()
	// 确保所有请求完成后关闭channel
	go func() {
		wg.Wait()
		close(results)
	}()
	// 从channel中依次取出结果并打印
	for res := range results {
		if res.Err != nil {
			fmt.Println(res.Err)
		}
		fmt.Println(res.Name, res.Status, res.Region)
	}
}
