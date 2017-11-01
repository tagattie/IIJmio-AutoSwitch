package main

import (
	"fmt"
)

func couponChangeByIdBasedCouponUsage(
	packetData *packetData,
	couponData *couponData) (
	map[string][]int,
	map[string]bool,
	int,
	map[string]bool) {
	// Extract latest packet usage data
	latestPacketData := extractLatestPacketData(packetData)
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Latest packet usage (in MB):")
		fmt.Printf("%+v\n\n", latestPacketData)
	}

	// Extract current coupon state and availability
	couponState, couponAmount := getCouponStateAndAmount(couponData)
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Current coupon state and amount:")
		fmt.Printf("%+v %+v\n\n", couponState, couponAmount)
	}

	// Make coupon state change request based on
	// latest packet data and current coupon state, and amount
	couponReqInfo := make(map[string]bool)
	for k, _ := range latestPacketData {
		// latestPacketData[k][0]: Packet data amount with coupon
		// latestPacketData[k][1]: Packet data amount without coupon
		if latestPacketData[k][0] >= config.Mio.MaxDailyAmount &&
			couponState[k] == true {
			couponReqInfo[k] = false
		} else if latestPacketData[k][0] < config.Mio.MaxDailyAmount &&
			couponState[k] == false {
			// Only when there is still coupon amount available
			if couponAmount > 0 {
				couponReqInfo[k] = true
			}
		}
	}

	return latestPacketData, couponState, couponAmount, couponReqInfo
}
