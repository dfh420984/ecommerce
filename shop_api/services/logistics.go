package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shop_api/config"
)

// LogisticsService 物流服务
type LogisticsService struct {
	apiKey   string
	customer string
}

var logisticsService *LogisticsService

// TrackInfo 物流轨迹信息
type TrackInfo struct {
	Time        string `json:"time"`
	Status      string `json:"status"`
	Description string `json:"desc"`
}

// TrackResult 快递100返回结果
type TrackResult struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    []TrackInfo `json:"data"`
}

func GetLogisticsService() *LogisticsService {
	if logisticsService == nil {
		cfg := config.Get()
		logisticsService = &LogisticsService{
			apiKey:   cfg.Logistics.APIKey,
			customer: cfg.Logistics.Customer,
		}
	}
	return logisticsService
}

// QueryTrack 查询物流轨迹
func (s *LogisticsService) QueryTrack(expressCompany, expressNo string) ([]TrackInfo, error) {
	// 如果没有配置API Key，返回模拟数据
	if s.apiKey == "" {
		return s.getMockTrackData(), nil
	}

	// 调用快递100 API
	url := fmt.Sprintf("https://poll.kuaidi100.com/poll/query.do?com=%s&nu=%s&key=%s",
		expressCompany, expressNo, s.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求物流API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result TrackResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "200" {
		return nil, fmt.Errorf("物流查询失败: %s", result.Message)
	}

	return result.Data, nil
}

// getMockTrackData 获取模拟物流数据（开发环境使用）
func (s *LogisticsService) getMockTrackData() []TrackInfo {
	return []TrackInfo{
		{
			Time:        "2026-04-24 10:30:00",
			Status:      "已签收",
			Description: "您的快件已被签收，感谢使用顺丰速运，期待再次为您服务",
		},
		{
			Time:        "2026-04-24 08:15:00",
			Status:      "派送中",
			Description: "快件交给张快递员，正在派送途中（联系电话：138****8888）",
		},
		{
			Time:        "2026-04-23 22:45:00",
			Status:      "运输中",
			Description: "快件到达【北京朝阳区营业部】",
		},
		{
			Time:        "2026-04-23 15:20:00",
			Status:      "运输中",
			Description: "快件在【北京顺义集散中心】完成分拣，准备发往【北京朝阳区营业部】",
		},
		{
			Time:        "2026-04-22 18:30:00",
			Status:      "已发货",
			Description: "顺丰速运已收取快件",
		},
	}
}
