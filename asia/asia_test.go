package asia

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"sync"
	"testing"
)

func Test(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(6)
	results := make(chan model.Result, 6)

	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := StarPlus(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := TLCGO(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := BilibiliMainland(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := HBOGO(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := Hotstar(req)
		results <- res
	}()
	go func() {
		defer wg.Done()
		req, _ := utils.ParseInterface("", "", "tcp4")
		res := MolaTV(req)
		results <- res
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf(res.Name + ": ")
		if res.Err != nil {
			fmt.Printf(res.Err.Error() + " ")
		}
		fmt.Println(res.Status, res.Region)
	}
}
