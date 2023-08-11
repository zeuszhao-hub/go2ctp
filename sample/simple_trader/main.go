// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/pseudocodes/go2ctp/ctp_tts"
	"github.com/pseudocodes/go2ctp/thost"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var SimnowEnv map[string]map[string]string = map[string]map[string]string{
	"td": {
		"7x24":      "tcp://180.168.146.187:10130",
		"telesim1":  "tcp://180.168.146.187:10201",
		"telesim2":  "tcp://180.168.146.187:10202",
		"moblesim3": "tcp://218.202.237.33:10203",
	},
	"md": {
		"7x24":      "tcp://180.168.146.187:10131",
		"telesim1":  "tcp://180.168.146.187:10211",
		"telesim2":  "tcp://180.168.146.187:10212",
		"moblesim3": "tcp://218.202.237.33:10213",
	},
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type baseSpi struct {
	brokerID   string
	investorID string
	password   string
	appid      string
	authCode   string

	requestID atomic.Int32

	ctp_tts.BaseTraderSpi
	tdapi ctp_tts.TraderApi
}

func CreateBaseSpi() *baseSpi {
	s := &baseSpi{
		// tdapi: tdapi,
		brokerID:   "9999",
		investorID: os.Getenv("SIMNOW_USER_ID"),       // <- 环境变量设置
		password:   os.Getenv("SIMNOW_USER_PASSWORD"), // <- 环境变量设置

		appid:    "simnow_client_test",
		authCode: "0000000000000000",
	}
	return s
}

func (s *baseSpi) OnFrontDisconnected(nReason int) {
	log.Printf("OnFrontDissconnected: %v\n", nReason)
}

func (p *baseSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Println("(OnHeartBeatWarning) nTimerLapse=", nTimeLapse)
}

func (s *baseSpi) OnFrontConnected() {
	var ret int
	log.Printf("OnFrontConnected\n")
	req := &thost.CThostFtdcReqAuthenticateField{}
	copy(req.BrokerID[:], []byte(s.brokerID))
	copy(req.UserID[:], []byte(s.investorID))
	ret = s.tdapi.ReqAuthenticate(req, int(s.requestID.Add(1)))

	log.Printf("user auth: %v\n", ret)
}

func (s *baseSpi) OnRspAuthenticate(pRspAuthenticateField *thost.CThostFtdcRspAuthenticateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	req := &thost.CThostFtdcReqUserLoginField{}
	copy(req.BrokerID[:], []byte(s.brokerID))
	copy(req.UserID[:], []byte(s.investorID))
	copy(req.Password[:], []byte(s.password))

	ret := s.tdapi.ReqUserLogin(req, int(s.requestID.Add(1)))
	log.Printf("user login: %v\n", ret)
}

func (s *baseSpi) OnRspUserLogin(pRspUserLogin *thost.CThostFtdcRspUserLoginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("OnRspUserLogin\n")
	req := &thost.CThostFtdcSettlementInfoConfirmField{}
	copy(req.BrokerID[:], []byte(s.brokerID))
	copy(req.InvestorID[:], []byte(s.investorID))

	ret := s.tdapi.ReqSettlementInfoConfirm(req, int(s.requestID.Add(1)))
	log.Printf("req_settlement_info_confirm : %v\n", ret)

}

// OnRspSettlementInfoConfirm 发送投资者结算单确认响应
func (s *baseSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm *thost.CThostFtdcSettlementInfoConfirmField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	// req := &goctp.QryInstrumentField{}
	// ret := s.tdapi.ReqQryInstrument(req, 3)
	// log.Printf("user qry ins: %v\n", ret)

	req := &thost.CThostFtdcQryTradingAccountField{}
	copy(req.BrokerID[:], []byte(s.brokerID))
	copy(req.InvestorID[:], []byte(s.investorID))
	ret := s.tdapi.ReqQryTradingAccount(req, int(s.requestID.Add(1)))
	if ret != 0 {
		log.Printf("req_qry_trading_account failed %v\n", ret)
	}

}

func (s *baseSpi) OnRspQryInstrument(pInstrument *thost.CThostFtdcInstrumentField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
}

// OnRspQryTradingAccount 请求查询资金账户响应
func (s *baseSpi) OnRspQryTradingAccount(pTradingAccount *thost.CThostFtdcTradingAccountField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	if bIsLast && !s.isErrorRspInfo(pRspInfo) {
		accountID := bytesToString2(pTradingAccount.AccountID[:])
		balance := pTradingAccount.Balance
		log.Printf("Account[%v] Balance[%.2f]\n", accountID, balance)
	}
}

// // 合约交易状态通知
// func (s *baseSpi) OnRtnInstrumentStatus(pInstrumentStatus *goctp.InstrumentStatusField) {
// 	// dump.P(pInstrumentStatus)
// }

// func (s *baseSpi) OnRspQryOrder(pOrder *goctp.OrderField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

// }

// func (s *baseSpi) OnRspQryInvestorPosition(pInvestorPosition *goctp.InvestorPositionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {

// }

// 错误应答
func (s *baseSpi) OnRspError(pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	s.isErrorRspInfo(pRspInfo)
}

// // 报单操作错误回报
// func (s *baseSpi) OnErrRtnOrderAction(pOrderAction *goctp.OrderActionField, pRspInfo *goctp.RspInfoField) {
// 	s.isErrorRspInfo(pRspInfo)
// }

// // 报单操作请求响应（撤单失败会触发）
// func (s *baseSpi) OnRspOrderAction(pInputOrderAction *goctp.InputOrderActionField, pRspInfo *goctp.RspInfoField, nRequestID int, bIsLast bool) {
// 	s.isErrorRspInfo(pRspInfo)
// }

// OnRtnTrade 成交通知（委托单在交易所成交了）
func (s *baseSpi) OnRtnTrade(pTrade *thost.CThostFtdcTradeField) {
}

// OnRtnOrder 报单通知（委托单）
func (s *baseSpi) OnRtnOrder(pOrder *thost.CThostFtdcOrderField) {

}

func (s *baseSpi) isErrorRspInfo(pRspInfo *thost.CThostFtdcRspInfoField) bool {

	// 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
	if pRspInfo == nil {
		return false
	}
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.ErrorID != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.ErrorID, Bytes2StringGBK(pRspInfo.ErrorMsg[:]))
	}
	return bResult

}

func bytesToString2(b []byte) string {
	before, _, _ := bytes.Cut(b, []byte{'\x00'})
	if len(before) > 0 {
		return string(before)
	}
	return ""
}

func Bytes2StringGBK(t []byte) string {
	msg, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(bytes.Split(t, []byte{'\x00'})[0])
	return strings.Trim(string(msg), "\u0000")
}

func sample1() {
	var tdapi ctp_tts.TraderApi
	if runtime.GOOS == "darwin" {
		tdapi = ctp_tts.CreateTraderApi(ctp_tts.TraderDynamicLibPath("../../ctp_tts/lib/v6.6.9_20220920/mac_arm64/thosttraderapi_se.dylib"), ctp_tts.TraderFlowPath("./data/"))
	} else if runtime.GOOS == "linux" {
		tdapi = ctp_tts.CreateTraderApi(ctp_tts.TraderDynamicLibPath("../../ctp_tts/lib/v6.6.9_20220920/lin64/thosttraderapi_se.so"), ctp_tts.TraderFlowPath("./data/"))
	}

	baseSpi := CreateBaseSpi()
	baseSpi.tdapi = tdapi

	tdapi.RegisterSpi(baseSpi)
	tdapi.RegisterFront("tcp://121.37.90.193:20002")

	tdapi.Init()

	println(tdapi.GetTradingDay())
	println(tdapi.GetApiVersion())

	tdapi.Join()
	select {}
}

func main() {
	sample1()
}
