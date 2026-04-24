package services

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"shop_api/config"
	"shop_api/database"
	"shop_api/models"
	"shop_api/types"
	"shop_api/utils"

	"github.com/google/uuid"
)

type WechatPayService struct {
	cfg *config.WechatConfig
}

type WechatUnifiedOrderRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Body           string `xml:"body"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
}

type WechatUnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
}

type WechatPayNotify struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	TradeType     string `xml:"trade_type"`
	TradeState    string `xml:"trade_state"`
	TotalFee      int    `xml:"total_fee"`
	CashFee       int    `xml:"cash_fee"`
}

func NewWechatPayService() *WechatPayService {
	return &WechatPayService{
		cfg: &config.Get().Wechat,
	}
}

func (s *WechatPayService) UnifiedOrder(order *models.Order) (string, error) {
	req := WechatUnifiedOrderRequest{
		AppID:          s.cfg.AppID,
		MchID:          s.cfg.MchID,
		NonceStr:       uuid.New().String(),
		Body:           "订单支付-" + order.OrderNo,
		OutTradeNo:     order.OrderNo,
		TotalFee:       int(order.PayAmount * 100),
		SpbillCreateIP: "127.0.0.1",
		NotifyURL:      s.cfg.NotifyURL,
		TradeType:      "NATIVE",
	}

	req.Sign = s.Sign(req)

	xmlData, _ := xml.Marshal(req)
	xmlData = []byte(xml.Header + string(xmlData))

	resp, err := s.httpPost("https://api.mch.weixin.qq.com/pay/unifiedorder", xmlData)
	if err != nil {
		return "", err
	}

	var result WechatUnifiedOrderResponse
	if err := xml.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	if result.ReturnCode != "SUCCESS" || result.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("wechat pay error: %s", result.ReturnMsg)
	}

	return result.CodeURL, nil
}

func (s *WechatPayService) Notify(data []byte) (*WechatPayNotify, error) {
	var notify WechatPayNotify
	if err := xml.Unmarshal(data, &notify); err != nil {
		return nil, err
	}

	if notify.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("notify error: %s", notify.ReturnMsg)
	}

	orderNo := notify.OutTradeNo
	var order models.Order
	if err := database.GetDB().Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}

	if order.PayStatus != models.PayStatusUnpaid {
		return &notify, nil
	}

	order.PayStatus = models.PayStatusPaid
	order.PayType = models.PayTypeWechat
	now := types.Now()
	order.PayTime = &now
	order.OrderStatus = models.OrderStatusPaid

	tx := database.GetDB().Begin()

	for _, item := range order.Items {
		database.GetDB().Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", database.GetDB().Raw("stock - ?", item.Quantity))
		database.GetDB().Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("sales", database.GetDB().Raw("sales + ?", item.Quantity))
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	payLog := models.PayLog{
		OrderID:    order.ID,
		OrderNo:    order.OrderNo,
		TradeNo:    notify.TransactionID,
		PayType:    models.PayTypeWechat,
		PayStatus:  1,
		PayAmount:  order.PayAmount,
		NotifyData: string(data),
	}
	database.GetDB().Create(&payLog)

	tx.Commit()

	return &notify, nil
}

func (s *WechatPayService) Sign(req WechatUnifiedOrderRequest) string {
	var keys []string
	values := make(map[string]string)

	values["appid"] = req.AppID
	values["mch_id"] = req.MchID
	values["nonce_str"] = req.NonceStr
	values["body"] = req.Body
	values["out_trade_no"] = req.OutTradeNo
	values["total_fee"] = strconv.Itoa(req.TotalFee)
	values["spbill_create_ip"] = req.SpbillCreateIP
	values["notify_url"] = req.NotifyURL
	values["trade_type"] = req.TradeType

	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for _, k := range keys {
		if values[k] != "" {
			signStr.WriteString(k + "=" + values[k] + "&")
		}
	}
	signStr.WriteString("key=" + s.cfg.APIKey)

	return strings.ToUpper(utils.MD5(signStr.String()))
}

func (s *WechatPayService) httpPost(url string, data []byte) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type AlipayService struct {
	cfg *config.AlipayConfig
}

type AlipayTradePreparePayRequest struct {
	OutTradeNo  string `json:"out_trade_no"`
	TotalAmount string `json:"total_amount"`
	Subject     string `json:"subject"`
}

type AlipayTradePreparePayResponse struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	OutTradeNo string `json:"out_trade_no"`
	QRCode     string `json:"qr_code"`
}

func NewAlipayService() *AlipayService {
	return &AlipayService{
		cfg: &config.Get().Alipay,
	}
}

func (s *AlipayService) TradePagePay(order *models.Order) (string, error) {
	return fmt.Sprintf("https://openapi.alipay.com/gateway.do?out_trade_no=%s&total_amount=%.2f&subject=订单%s",
		order.OrderNo, order.PayAmount, order.OrderNo), nil
}

func (s *AlipayService) Notify(data map[string]string) error {
	outTradeNo := data["out_trade_no"]
	tradeStatus := data["trade_status"]

	var order models.Order
	if err := database.GetDB().Where("order_no = ?", outTradeNo).First(&order).Error; err != nil {
		return err
	}

	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		if order.PayStatus == models.PayStatusUnpaid {
			order.PayStatus = models.PayStatusPaid
			order.PayType = models.PayTypeAlipay
			now := types.Now()
			order.PayTime = &now
			order.OrderStatus = models.OrderStatusPaid

			tx := database.GetDB().Begin()

			for _, item := range order.Items {
				database.GetDB().Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", database.GetDB().Raw("stock - ?", item.Quantity))
				database.GetDB().Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("sales", database.GetDB().Raw("sales + ?", item.Quantity))
			}

			if err := tx.Save(&order).Error; err != nil {
				tx.Rollback()
				return err
			}

			tx.Commit()
		}
	}

	return nil
}

type PayService struct {
	wechat *WechatPayService
	alipay *AlipayService
}

func NewPayService() *PayService {
	return &PayService{
		wechat: NewWechatPayService(),
		alipay: NewAlipayService(),
	}
}

func (s *PayService) GetWechatPayURL(order *models.Order) (string, error) {
	return s.wechat.UnifiedOrder(order)
}

func (s *PayService) GetAlipayURL(order *models.Order) (string, error) {
	return s.alipay.TradePagePay(order)
}

// IsWechatConfigured 检查微信支付是否已配置
func (s *PayService) IsWechatConfigured() bool {
	cfg := config.Get()
	if cfg == nil {
		return false
	}
	// 检查必要字段是否还是默认占位符
	return cfg.Wechat.AppID != "" &&
		cfg.Wechat.AppID != "your-wechat-appid" &&
		cfg.Wechat.MchID != "" &&
		cfg.Wechat.MchID != "your-wechat-mchid" &&
		cfg.Wechat.APIKey != "" &&
		cfg.Wechat.APIKey != "your-wechat-apikey"
}

// IsAlipayConfigured 检查支付宝是否已配置
func (s *PayService) IsAlipayConfigured() bool {
	cfg := config.Get()
	if cfg == nil {
		return false
	}
	// 检查必要字段是否还是默认占位符
	return cfg.Alipay.AppID != "" &&
		cfg.Alipay.AppID != "your-alipay-appid" &&
		cfg.Alipay.PrivateKey != "" &&
		cfg.Alipay.PrivateKey != "your-alipay-private-key"
}

func (s *PayService) WechatNotify(data []byte) error {
	_, err := s.wechat.Notify(data)
	return err
}

func (s *PayService) AlipayNotify(data map[string]string) error {
	return s.alipay.Notify(data)
}

func (s *PayService) QueryOrder(orderNo string) error {
	return nil
}

func (s *PayService) Refund(order *models.Order) error {
	return nil
}

// WechatRefund 微信退款
func (s *PayService) WechatRefund(order *models.Order, refundAmount float64) (string, error) {
	if !s.IsWechatConfigured() {
		// 未配置，返回模拟交易号
		return "mock_wechat_refund_" + order.OrderNo, nil
	}

	// TODO: 实现真实的微信退款逻辑
	// 这里需要调用微信支付V3 API的退款接口
	return "", fmt.Errorf("微信退款功能待实现")
}

// AlipayRefund 支付宝退款
func (s *PayService) AlipayRefund(order *models.Order, refundAmount float64) (string, error) {
	if !s.IsAlipayConfigured() {
		// 未配置，返回模拟交易号
		return "mock_alipay_refund_" + order.OrderNo, nil
	}

	// TODO: 实现真实的支付宝退款逻辑
	return "", fmt.Errorf("支付宝退款功能待实现")
}

var payService *PayService
var payOnce sync.Once

func GetPayService() *PayService {
	payOnce.Do(func() {
		payService = NewPayService()
	})
	return payService
}
