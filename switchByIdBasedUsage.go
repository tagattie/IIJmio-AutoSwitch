package main

import (
	"log"
)

func switchByIdBasedUsage(
	packetData *packetData,
	couponData *couponData) (
	map[string][]int,
	map[string]bool,
	int,
	map[string]bool) {
	// Extract current coupon state and availability
	couponState, couponAmount := getCouponStateAndAmount(couponData)
	log.Println("Current coupon state and amount:")
	log.Printf("%+v %+v\n", couponState, couponAmount)

	// Extract latest packet usage data
	latestPacketData := extractLatestPacketData(packetData)
	log.Println("Latest packet usage (in MB):")
	log.Printf("%+v\n", latestPacketData)

	// Make coupon state change request based on
	// latest packet data and current coupon state, and amount
	couponReqInfo := make(map[string]bool)
	for k := range latestPacketData {
		// latestPacketData[k][0]: Packet data amount with coupon
		// latestPacketData[k][1]: Packet data amount without coupon
		if latestPacketData[k][0] >= config.Mio.MaxDailyAmount && couponState[k] {
			couponReqInfo[k] = false
		} else if latestPacketData[k][0] < config.Mio.MaxDailyAmount && !couponState[k] {
			// Only when there is still coupon amount available
			if couponAmount > 0 {
				couponReqInfo[k] = true
			}
		}
	}
	log.Println("Coupon request info:")
	log.Printf("%v\n", couponReqInfo)

	return latestPacketData, couponState, couponAmount, couponReqInfo
}
