package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"12306-mcp-service/models"
)

// TicketService 车票服务
type TicketService struct {
	cityCodeMap models.CityCodeMap
	codeCityMap models.CodeCityMap
}

// NewTicketService 创建车票服务实例
func NewTicketService() *TicketService {
	return &TicketService{}
}

// QueryTickets 查询车票信息
func (s *TicketService) QueryTickets(fromStation, toStation, date string) ([]models.TicketInfo, error) {
	// 获取车票数据
	tickets, err := s.generateMockTickets(fromStation, toStation, date)
	if err != nil {
		return nil, fmt.Errorf("获取车票数据失败: %v", err)
	}

	return tickets, nil
}

// generateMockTickets 获取12306车票数据
func (s *TicketService) generateMockTickets(fromStation, toStation, date string) ([]models.TicketInfo, error) {
	// 1. 读取城市代码文件
	log.Printf("开始查询车票信息: 从 %s 到 %s，日期 %s", fromStation, toStation, date)
	err := s.loadCityCodeMap()
	cityData := s.cityCodeMap
	if err != nil {
		log.Printf("读取城市代码文件失败: %v", err)
		return nil, fmt.Errorf("读取城市代码文件失败: %v", err)
	}
	// 校验日期是否小于当天
	inputDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("日期格式错误: %v", err)
	}
	
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	
	if inputDate.Before(today) {
		return nil, fmt.Errorf("查询日期不能小于当天")
	}


	// 2. 查找城市代码
	fromCode, toCode, err := s.getCityCodes(cityData, fromStation, toStation)
	if err != nil {
		log.Printf("获取城市代码失败: %v", err)
		return nil, err
	}
	log.Printf("城市代码转换: %s -> %s, %s -> %s", fromStation, fromCode, toStation, toCode)

	// 3. 发送请求到12306
	tickets, err := s.queryTicketsFrom12306(fromCode, toCode, date)
	if err != nil {
		// 如果请求失败，使用模拟数据作为后备
		log.Printf("12306查询失败，使用模拟数据: %v", err)
		return nil, errors.New("获取数据异常,请稍后再试吧")
	}

	log.Printf("成功从12306获取到 %d 条车票信息", len(tickets))
	return tickets, nil
}

// loadCityCodeMap 加载城市代码映射
func (s *TicketService) loadCityCodeMap() error {
	m1 := make(map[string]string)
	m2 := make(map[string]string)

	//尝试从文件获取
	bs, err := os.ReadFile("./city.json")
	if err != nil {
		return fmt.Errorf("读取城市代码文件失败: %v", err)
	}

	for _, v := range strings.Split(string(bs), "@") {
		/*
			得到大概如下格式
			bjb|北京北|VAP|beijingbei|bjb|0|0357|北京|||
		*/
		if list := strings.Split(v, "|"); len(list) >= 3 {
			m1[list[1]] = list[2]
			m2[list[2]] = list[1]
		}

	}
	// 将城市代码映射保存到常量Map中
	s.cityCodeMap = m1
	s.codeCityMap = m2

	return nil
}

// getCityCodes 获取城市代码
func (s *TicketService) getCityCodes(cityMap models.CityCodeMap, fromCity, toCity string) (string, string, error) {
	fromCode := fromCity
	toCode := toCity

	if len(fromCity) != 3 || !s.isAllUpperCase(fromCity) {
		code, ok := cityMap[fromCity]
		if !ok {
			found := false
			for _, cityCode := range cityMap {
				if cityCode == fromCity {
					fromCode = cityCode
					found = true
					break
				}
			}
			if !found {
				return "", "", fmt.Errorf("未找到出发城市代码: %s", fromCity)
			}
		} else {
			fromCode = code
		}
	}

	if len(toCity) != 3 || !s.isAllUpperCase(toCity) {
		code, ok := cityMap[toCity]
		if !ok {
			found := false
			for _, cityCode := range cityMap {
				if cityCode == toCity {
					toCode = cityCode
					found = true
					break
				}
			}
			if !found {
				return "", "", fmt.Errorf("未找到到达城市代码: %s", toCity)
			}
		} else {
			toCode = code
		}
	}

	return fromCode, toCode, nil
}

// isAllUpperCase 检查字符串是否全部由大写字母组成
func (s *TicketService) isAllUpperCase(str string) bool {
	for _, r := range str {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

// queryTicketsFrom12306 从12306查询车票
func (s *TicketService) queryTicketsFrom12306(fromStation, toStation, date string) ([]models.TicketInfo, error) {
	queryURLs := []string{
		"https://kyfw.12306.cn/otn/leftTicket/queryG",
	}

	var lastErr error
	for _, baseURL := range queryURLs {
		tickets, err := s.queryFromURL(baseURL, fromStation, toStation, date)
		if err != nil {
			log.Printf("使用接口 %s 查询失败: %v", baseURL, err)
			lastErr = err
			continue
		}
		return tickets, nil
	}

	return nil, fmt.Errorf("所有查询接口均失败，最后错误: %v", lastErr)
}

// queryFromURL 从指定URL查询车票信息
func (s *TicketService) queryFromURL(baseURL, fromStation, toStation, date string) ([]models.TicketInfo, error) {
	params := url.Values{}
	params.Add("leftTicketDTO.train_date", date)
	params.Add("leftTicketDTO.from_station", fromStation)
	params.Add("leftTicketDTO.to_station", toStation)
	params.Add("purpose_codes", "ADULT")

	orderedKeys := []string{
		"leftTicketDTO.train_date",
		"leftTicketDTO.from_station",
		"leftTicketDTO.to_station",
		"purpose_codes",
	}
	var queryParts []string
	for _, key := range orderedKeys {
		value := params.Get(key)
		if value != "" {
			queryParts = append(queryParts, key+"="+value)
		}
	}
	queryStr := strings.Join(queryParts, "&")

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryStr)
	log.Printf("发送请求到12306: %s", requestURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")

	cookieStr := "_uab_collina=174469819159739930224404; JSESSIONID=34118AC2147BA68D78C247093B83D4BA; guidesStatus=off; highContrastMode=defaltMode; cursorStatus=off; _jc_save_wfdc_flag=dc; _jc_save_showIns=true; _jc_save_toStation=%u5B9C%u660C%2CYCN; _jc_save_fromDate=2025-04-30; _jc_save_fromStation=%u5357%u660C%2CNCG; BIGipServerpassport=786956554.50215.0000; route=495c805987d0f5c8c84b14f60212447d; BIGipServerotn=1675165962.24610.0000; _jc_save_toDate=2025-04-17"
	cookies := strings.Split(cookieStr, "; ")
	for _, c := range cookies {
		parts := strings.SplitN(c, "=", 2)
		if len(parts) == 2 {
			name, _ := url.QueryUnescape(parts[0])
			value, _ := url.QueryUnescape(parts[1])
			req.AddCookie(&http.Cookie{Name: name, Value: value})
		}
	}

	log.Printf("开始发送请求...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("收到响应，状态码: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	log.Printf("响应内容长度: %d 字节", len(body))

	if len(body) == 0 {
		return nil, fmt.Errorf("响应内容为空")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("JSON解析失败: %v, 内容: %s", err, string(body))
		return nil, fmt.Errorf("解析响应JSON失败: %v", err)
	}

	status, ok := response["status"].(bool)
	if !ok || !status {
		errorMsg := "未知错误"
		if msg, ok := response["messages"].(string); ok && msg != "" {
			errorMsg = msg
		} else if messages, ok := response["messages"].([]interface{}); ok && len(messages) > 0 {
			errorMsg = fmt.Sprintf("%v", messages[0])
		}
		return nil, fmt.Errorf("请求失败: %s", errorMsg)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("响应数据格式错误，未找到data字段")
	}

	result, ok := data["result"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("响应数据格式错误，未找到result字段")
	}

	log.Printf("解析到 %d 条车次信息", len(result))
	if len(result) == 0 {
		return nil, fmt.Errorf("没有找到符合条件的车票")
	}

	tickets := make([]models.TicketInfo, 0, len(result))
	for i, item := range result {
		info, ok := item.(string)
		if !ok {
			log.Printf("第 %d 条车次信息格式错误", i+1)
			continue
		}

		ticket, err := s.parseTrainInfo(info, date)
		if err != nil {
			log.Printf("解析第 %d 条车次信息失败: %v", i+1, err)
			continue
		}

		tickets = append(tickets, ticket)
	}

	if len(tickets) == 0 {
		return nil, fmt.Errorf("未找到符合条件的车票信息")
	}

	return tickets, nil
}

// parseTrainInfo 解析车次信息
func (s *TicketService) parseTrainInfo(info, date string) (models.TicketInfo, error) {
	parts := strings.Split(info, "|")
	if len(parts) < 40 {
		return models.TicketInfo{}, fmt.Errorf("车次信息格式错误")
	}

	departTimeStr := fmt.Sprintf("%s %s:00", date, parts[8])
	departTime, err := time.ParseInLocation("2006-01-02 15:04:05", departTimeStr, time.Local)
	if err != nil {
		return models.TicketInfo{}, fmt.Errorf("解析出发时间失败: %v", err)
	}

	arriveTime := s.calculateArriveTime(departTime, parts[10])

	seats := s.parseSeats(parts)

	return models.TicketInfo{
		TrainNo:     parts[3],
		FromStation: s.codeCityMap[parts[6]],
		ToStation:   s.codeCityMap[parts[7]],
		DepartTime:  departTime.Format("2006-01-02 15:04:05"),
		ArriveTime:  arriveTime.Format("2006-01-02 15:04:05"),
		Duration:    parts[10],
		Seats:       seats,
	}, nil
}

// calculateArriveTime 计算到达时间
func (s *TicketService) calculateArriveTime(departTime time.Time, duration string) time.Time {
	timeParts := strings.Split(duration, ":")
	if len(timeParts) != 2 {
		return departTime.Add(3 * time.Hour)
	}

	hours, _ := strconv.Atoi(timeParts[0])
	minutes, _ := strconv.Atoi(timeParts[1])

	return departTime.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
}

// parseSeats 解析座位信息
func (s *TicketService) parseSeats(parts []string) []models.SeatInfo {
	seats := make([]models.SeatInfo, 0)
	log.Printf("parseSeats: %v", parts[39])
	if parts[39] != "" {
		ticketInfo := s.parseTicketStockAndPrice(parts[39])
		if len(ticketInfo) > 0 {
			return ticketInfo
		}
	}

	return seats
}

// parseTicketStockAndPrice 解析座位票价信息
func (s *TicketService) parseTicketStockAndPrice(info string) []models.SeatInfo {
	seats := make([]models.SeatInfo, 0)

	for i := 0; i < len(info); i += 10 {
		if i+10 > len(info) {
			break
		}

		segment := info[i : i+10]

		seatType := "无座"
		if segment[0] != '3' {
			if t, ok := models.SeatTypeMap[string(segment[0])]; ok {
				seatType = t
			}
		}

		priceStr := fmt.Sprintf("%s.%s", segment[1:5], segment[5:6])
		price, _ := strconv.ParseFloat(priceStr, 64)

		count, _ := strconv.Atoi(segment[7:])

		seats = append(seats, models.SeatInfo{
			Type:  seatType,
			Count: count,
			Price: price,
		})
	}

	return seats
}
