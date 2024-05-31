package africa

import (
	"fmt"
	"sync"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func Test(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	results := make(chan model.Result, 3)
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