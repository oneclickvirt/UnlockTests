package vn

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/utils"
	"testing"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := BilibiliVN(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println("BilibiliVN", ": ", res.Status, res.Region)

	res = KPLUS(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	//res = TV360(req)
	//if res.Err != nil {
	//	fmt.Println(res.Err)
	//}
	//fmt.Println(res.Name, ": ", res.Status, res.Region)
}
