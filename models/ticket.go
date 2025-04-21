package models

import (
	"encoding/json"
)

// TicketInfo 定义车票信息结构
type TicketInfo struct {
	TrainNo     string     `json:"车次"`   // 车次
	FromStation string     `json:"出发站"`  // 出发站
	ToStation   string     `json:"到达站"`  // 到达站
	DepartTime  string     `json:"出发时间"` // 出发时间
	ArriveTime  string     `json:"到达时间"` // 到达时间
	Duration    string     `json:"历时"`   // 历时
	Seats       []SeatInfo `json:"余票信息"` // 余票信息
}

// String 实现Stringer接口，美化输出
func (t TicketInfo) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

// SeatInfo 定义座位信息
type SeatInfo struct {
	Type  string  `json:"座位类型"` // 座位类型
	Count int     `json:"余票数量"` // 剩余数量
	Price float64 `json:"价格"`   // 价格
}

// CityCodeMap 城市代码映射
type CityCodeMap map[string]string

type CodeCityMap map[string]string

// SeatTypeMap 座位类型映射
var SeatTypeMap = map[string]string{
	"A": "高级动卧", "B": "混编硬座", "C": "混编硬卧", "D": "优选一等座", "E": "特等软座", "F": "动卧",
	"G": "二人软包", "H": "一人软包", "I": "一等卧", "J": "二等卧", "K": "混编软座", "L": "混编软卧",
	"M": "一等座", "O": "二等座", "P": "特等座", "Q": "多功能座", "S": "二等包座", "W": "无座",
	"0": "棚车", "1": "硬座", "2": "软座", "3": "硬卧", "4": "软卧", "5": "包厢硬卧", "6": "高级软卧",
	"7": "一等软座", "8": "二等软座", "9": "商务座",
}
