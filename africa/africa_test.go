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
	wg.Add(3)
	results := make(chan model.Result, 3)
	go func() {
		defer wg.Done()
		c, _ := utils.ParseInterface("", "", "tcp4")
		res := Showmax(c)
		results <- res
	}()
	go func() {
		defer wg.Done()
		c, _ := utils.ParseInterface("", "", "tcp4")
		res := DSTV(c)
		results <- res
	}()
	go func() {
		defer wg.Done()
		c, _ := utils.ParseInterface("", "", "tcp4")
		res := BeinConnect(c)
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
		fmt.Println(res.Status, res.Region, res.Info, res.UnlockType)
	}
}
