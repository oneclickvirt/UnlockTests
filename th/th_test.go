package th

import (
	"fmt"
	"testing"

	"github.com/oneclickvirt/UnlockTests/utils"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := BilibiliTH(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println("BilibiliTH", ": ", res.Status, res.Region)

	res = AISPlay(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = TrueID(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)
}

// func Test(t *testing.T) {
// 	req, _ := utils.ParseInterface("", "", "tcp4")

// 	res := AISPlay(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region)
// }
