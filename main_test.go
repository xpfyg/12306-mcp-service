package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"12306-mcp-service/handlers"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestQueryTicketsHandler(t *testing.T) {
	testCases := []struct {
		name          string
		fromStation   string
		toStation     string
		date          string
		expectedError bool
	}{
		{
			name:          "正常查询",
			fromStation:   "杭州",
			toStation:     "武汉",
			date:          "2025-04-26",
			expectedError: false,
		},
	}

	handler := handlers.NewTicketHandler()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"from_station": tc.fromStation,
						"to_station":   tc.toStation,
						"date":         tc.date,
					},
				},
			}

			result, err := handler.QueryTickets(ctx, request)
			if tc.expectedError {
				if err == nil {
					t.Errorf("期望出现错误,但是没有错误")
				}
			} else {
				if err != nil {
					t.Errorf("未期望出现错误,但是出现错误: %v", err)
				}

				// 打印查询结果
				t.Logf("查询结果: %+v", result)
			}

		})
	}
}

func TestGetStationName(t *testing.T) {

	stationMap, err := GetStationName()
	if err != nil {
		t.Errorf("获取车站名称失败: %v", err)
		return
	}

	// 打印车站名称映射
	for station, code := range stationMap {
		fmt.Printf("车站: %s, 代码: %s\n", station, code)
	}

}

func GetStationName() (map[string]string, error) {

	m := make(map[string]string)
	keuValid := "valid"

	//尝试从文件获取
	bs, err := os.ReadFile("./city.json")
	if err == nil {
		var validTime int64
		if json.Unmarshal(bs, &m) == nil {
			if validTime, err = strconv.ParseInt(m[keuValid], 10, 64); err == nil && time.Now().Unix()-validTime < 0 {
				delete(m, keuValid)
				return m, nil
			}
		}
	}
	/*
		响应格式如下
		var station_names ='@bjb|北京北|VAP|beijingbei|bjb|0|0357|北京|||@bjd|北京东|BOP|beijingdong|bjd|1|0357|北京|||@...'
	*/

	for _, v := range strings.Split(string(bs), "@") {
		/*
			得到大概如下格式
			bjb|北京北|VAP|beijingbei|bjb|0|0357|北京|||
		*/
		if list := strings.Split(v, "|"); len(list) >= 3 {
			m[list[1]] = list[2]
		}

	}

	return m, nil
}
