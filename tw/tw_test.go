package tw

import (
	"fmt"
	"testing"

	"github.com/oneclickvirt/UnlockTests/utils"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := Tw4gtv(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = HamiVideo(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = LiTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = LineTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = MyVideo(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)

	res = BilibiliTW(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println("BilibiliTW", ": ", res.Status, res.Region)

	res = KKTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region)
}

// func Test(t *testing.T) {
// 	req, _ := utils.ParseInterface("", "", "tcp4")

// 	res := BahamutAnime(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region)

// 	res = Catchplay(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region)
// }
