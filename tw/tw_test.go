package tw

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/utils"
	"testing"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := Tw4gtv(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = HamiVideo(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = LiTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = LineTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = MyVideo(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = BilibiliTW(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println("BilibiliTW", ": ", res.Status, res.Region, res.UnlockType)

	res = KKTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = BahamutAnime(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
}

// func Test(t *testing.T) {
// 	req, _ := utils.ParseInterface("", "", "tcp4")

// 	res := BahamutAnime(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

// 	res = Catchplay(req)
// 	if res.Err != nil {
// 		fmt.Println(res.Err)
// 	}
// 	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
// }
