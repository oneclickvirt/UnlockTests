package uk

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/utils"
	"testing"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := BBCiPlayer(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = BritBox(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Channel4(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Channel5(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = ITVX(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = SkyGo(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = DiscoveryPlus(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
}

// func Test(t *testing.T) {
// 	req, _ := utils.ParseInterface("", "", "tcp4")

// 	res := Channel5(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
// }
