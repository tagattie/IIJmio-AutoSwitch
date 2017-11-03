package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

// Packet-related data
type packetLog struct {
	Date          string `json:"date"`
	WithCoupon    int    `json:"withCoupon"`
	WithoutCoupon int    `json:"withoutCoupon"`
}
type pDataHdoInfo struct {
	HdoServiceCode string      `json:"hdoServiceCode"`
	PacketLog      []packetLog `json:"packetLog"`
}
type pDataHduInfo struct {
	HduServiceCode string      `json:"hduServiceCode"`
	PacketLog      []packetLog `json:"packetLog"`
}
type packetLogInfo struct {
	HddServiceCode string         `json:"hddServiceCode`
	Plan           string         `json:"plan"`
	HdoInfo        []pDataHdoInfo `json:"hdoInfo"`
	HduInfo        []pDataHduInfo `json:"hduInfo"`
}
type packetData struct {
	ReturnCode    string          `json:"returnCode"`
	PacketLogInfo []packetLogInfo `json:"packetLogInfo"`
}

type PacketLogs []packetLog

func (pls PacketLogs) Len() int {
	return len(pls)
}
func (pls PacketLogs) Less(i, j int) bool {
	return (pls[i].Date > pls[j].Date)
}
func (pls PacketLogs) Swap(i, j int) {
	pls[i], pls[j] = pls[j], pls[i]
}

type StringIntSlice struct {
	key   string
	value []int
}
type ListStringIntSlice []StringIntSlice

func (l ListStringIntSlice) Len() int {
	return len(l)
}
func (l ListStringIntSlice) Less(i, j int) bool {
	return (l[i].value[0] > l[j].value[0])
}
func (l ListStringIntSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func decodePacketDataJSON(packetBytes []byte) (*packetData, error) {
	var pd packetData

	if err := json.Unmarshal(packetBytes, &pd); err != nil {
		log.Println("Packet request JSON unmarshal error: ", err)
		return nil, err
	}

	if returnCode := pd.ReturnCode; returnCode != "OK" {
		log.Println("Packet request return code error: ", returnCode)
		err := fmt.Errorf("Packet request return code is %+v", returnCode)
		return &pd, err
	}

	return &pd, nil
}

func extractLatestPacketData(packetData *packetData) map[string][]int {
	latestPacketData := make(map[string][]int)
	// Extract latest (= today's) packet data
	plis := packetData.PacketLogInfo
	if plisLength := len(plis); plisLength > 0 {
		for i := range plis {
			hdois := plis[i].HdoInfo
			if hdoisLength := len(hdois); hdoisLength > 0 {
				for j, hdoInfo := range hdois {
					var pls PacketLogs
					pls = hdois[j].PacketLog
					if plLength := len(pls); plLength > 0 {
						// Sort data in descending order by date
						sort.Sort(pls)
						// Index 0 means today
						pl := make([]int, 2)
						pl[0] = pls[0].WithCoupon
						pl[1] = pls[0].WithoutCoupon
						latestPacketData[hdoInfo.HdoServiceCode] = pl
					}
				}
			}
		}
	}
	return latestPacketData
}

func extractLatestPacketDataSortedByAmount(packetData *packetData) ListStringIntSlice {
	latestPacketData := extractLatestPacketData(packetData)
	// Sort latest package data by coupon usage amount
	list := ListStringIntSlice{}
	for k, v := range latestPacketData {
		list = append(list, StringIntSlice{k, v})
	}
	sort.Sort(list)
	return list
}
