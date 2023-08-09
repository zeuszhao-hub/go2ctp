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
	"log"
	"time"

	"github.com/pseudocodes/go2ctp/ctp"
	"github.com/pseudocodes/go2ctp/ctp_tts"
	"github.com/pseudocodes/go2ctp/thost"
)

/*
Simnow是上期技术提供的CTP程序测试、模拟、学习的模拟平台。

7x24环境-电信：
交易前置: tcp://180.168.146.187:10130
行情前置: tcp://180.168.146.187:10131

仿真环境1-电信：交易时段同实盘
交易前置: tcp://180.168.146.187:10201
行情前置: tcp://180.168.146.187:10211

仿真环境2-电信：交易时段同实盘
交易前置: tcp://180.168.146.187:10202
行情前置: tcp://180.168.146.187:10212

仿真环境3-移动：交易时段同实盘
交易前置: tcp://218.202.237.33:10203
行情前置: tcp://218.202.237.33:10213
*/
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
	ctp_tts.BaseMdSpi
	// ctp.BaseMdSpi
	mdapi ctp_tts.MdApi
}

func CreateBaseSpi() *baseSpi {
	s := &baseSpi{}
	return s
}

func (s *baseSpi) OnFrontConnected() {
	log.Printf("OnFrontConnected\n")

	loginR := &thost.CThostFtdcReqUserLoginField{}
	copy(loginR.BrokerID[:], "9999")
	copy(loginR.UserID[:], "2011")

	ret := s.mdapi.ReqUserLogin(loginR, 1)

	log.Printf("user log: %v\n", ret)
}

func (s *baseSpi) OnHeartBeatWarning(timelapse int) {
	log.Printf("OnHeartBeatWarning: %v\n", timelapse)
}

func (s *baseSpi) OnFrontDisconnected(nReason int) {
	log.Printf("OnFrontDisconnected: %v\n", nReason)
}

func (s *baseSpi) OnRspUserLogin(pRspUserLogin *thost.CThostFtdcRspUserLoginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("RspUserLogin: %+v\nRspInfo: %+v\n", pRspUserLogin, nil)
	s.mdapi.SubscribeMarketData("ag2310")
}

func (s *baseSpi) OnRspSubMarketData(pSpecificInstrument *thost.CThostFtdcSpecificInstrumentField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Printf("instrumentID: %+v\n  RspInfo: %+v\n", pSpecificInstrument, nil)
}

func (s *baseSpi) OnRtnDepthMarketData(pDepthMarketData *thost.CThostFtdcDepthMarketDataField) {
	log.Printf("OnRtnDeptMarketData\n")
}

func (s *baseSpi) OnRspError(pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("RspInfo: %+v\n", pRspInfo.ErrorID)

}

type baseSpi2 struct {
	ctp.BaseMdSpi
	mdapi ctp.MdApi
}

func CreateBaseSpi2() *baseSpi2 {
	s := &baseSpi2{}

	s.OnFrontConnectedCallback = func() {
		log.Printf("OnFrontConnected\n")

		loginR := &thost.CThostFtdcReqUserLoginField{}
		copy(loginR.BrokerID[:], "9999")
		copy(loginR.UserID[:], "2011")

		ret := s.mdapi.ReqUserLogin(loginR, 1)

		log.Printf("user log: %v\n", ret)
	}
	s.OnFrontDisconnectedCallback = func(nReason int) {
		log.Printf("OnFrontDisconnected: %v\n", nReason)
	}
	s.OnRspUserLoginCallback = func(pRspUserLogin *thost.CThostFtdcRspUserLoginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
		log.Printf("RspUserLogin: %+v\nRspInfo: %+v\n", pRspUserLogin, nil)
		s.mdapi.SubscribeMarketData("ag2310")
	}
	s.OnRtnDepthMarketDataCallback = func(pDepthMarketData *thost.CThostFtdcDepthMarketDataField) {
		// log.Printf("tick {%+v}\n", quote)
		log.Printf("OnRtnDeptMarketData 222\n")
	}
	return s

}

func sample1() {

	mdapi := ctp_tts.CreateMdApi(ctp_tts.MdDynamicLibPath("../../ctp_tts/lib/v6.6.9_20220920/mac_arm64/libthostmduserapi_se.dylib"), ctp_tts.MdFlowPath("./data/"), ctp_tts.MdUsingUDP(false), ctp_tts.MdMultiCast(false))
	baseSpi := CreateBaseSpi()
	baseSpi.mdapi = mdapi
	mdapi.RegisterSpi(baseSpi)

	mdapi.RegisterFront("tcp://121.37.80.177:20004")
	mdapi.Init()

	println(mdapi.GetApiVersion())
	println(mdapi.GetTradingDay())

	// mdapi.Join()
	for {
		time.Sleep(10 * time.Second)
	}

}

func sample2() {
	mdapi := ctp.CreateMdApi(ctp.MdFlowPath("./data/"), ctp.MdUsingUDP(false), ctp.MdMultiCast(false))
	baseSpi2 := CreateBaseSpi2()
	baseSpi2.mdapi = mdapi
	mdapi.RegisterSpi(baseSpi2)

	mdapi.RegisterFront("tcp://140.206.244.33:11616")
	mdapi.Init()

	println(mdapi.GetApiVersion())
	println(mdapi.GetTradingDay())

	// mdapi.Join()
	for {
		time.Sleep(10 * time.Second)
	}

}

func main() {

	sample1()
	// sample2()
}
