package openblive

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const getAdminApi = "https://api.live.bilibili.com/xlive/web-room/v1/roomAdmin/get_by_room?roomid=%d&page_size=99"

func getAdmins(roomId int) (adminNames []string) {
	uri := fmt.Sprintf(getAdminApi, roomId)
	val, err := resty.New().R().Get(uri)
	// test only, add me to admin
	adminNames = append(adminNames, "Aynakeya")
	if err != nil {
		return adminNames
	}
	valStr := val.String()
	if !gjson.Valid(valStr) {
		return adminNames
	}
	result := gjson.Parse(valStr)
	result.Get("data.data").ForEach(func(key, value gjson.Result) bool {
		adminNames = append(adminNames, value.Get("uname").String())
		return true
	})
	return adminNames
}
