package hk

import (
	"fmt"
	"testing"

	"github.com/oneclickvirt/UnlockTests/utils"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := BilibiliHKMO(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println("BilibiliHKMO", ": ", res.Status, res.Region, res.UnlockType)

	res = MyTvSuper(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = ViuTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = NowE(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = HoyTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
}
