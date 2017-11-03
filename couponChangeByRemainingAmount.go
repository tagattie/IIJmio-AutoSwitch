package main

import (
	"log"
	"time"
)

func couponChangeByRemainingCouponAmount(
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

	// Get today's date
	today := time.Now().Day()
	log.Println("Today's day:")
	log.Printf("%+v\n", today)
	// Calculate coupon quotas to apply
	couponAmountQuota := []int{}
	for i := 1; i <= len(couponState); i++ {
		couponAmountQuota = append(couponAmountQuota,
			config.Mio.StartingAmount-
				((today-1)*config.Mio.MaxDailyAmount*len(couponState))-
				(i*config.Mio.MaxDailyAmount))
	}
	log.Println("Remaining coupon amount quotas (in MB):")
	log.Printf("%+v\n", couponAmountQuota)

	// Extract latest packet usage data
	latestPacketDataSortedByAmount := extractLatestPacketDataSortedByAmount(packetData)
	log.Println("Latest packet usage (in MB):")
	log.Printf("%+v\n", latestPacketDataSortedByAmount)
	latestPacketData := make(map[string][]int)
	for _, v := range latestPacketDataSortedByAmount {
		latestPacketData[v.key] = v.value
	}

	// Make coupon state change request based on remaining coupon amount
	couponReqInfo := make(map[string]bool)
	for i := range latestPacketDataSortedByAmount {
		// fmt.Printf("latestPacketDataSortedByAmount[%d]: ", i)
		// fmt.Printf("%v\n", latestPacketDataSortedByAmount[i])
		// fmt.Printf("couponState[latestPacketDataSortedByAmount[%d].key]: ", i)
		// fmt.Printf("%v\n", couponState[latestPacketDataSortedByAmount[i].key])
		if couponAmount < couponAmountQuota[i] &&
			couponState[latestPacketDataSortedByAmount[i].key] {
			couponReqInfo[latestPacketDataSortedByAmount[i].key] = false
		} else if couponAmount > couponAmountQuota[i] &&
			!couponState[latestPacketDataSortedByAmount[i].key] {
			// Only when there is still coupon amount available
			if couponAmount > 0 {
				couponReqInfo[latestPacketDataSortedByAmount[i].key] = true
			}
		}
	}
	log.Println("Coupon request info:")
	log.Printf("%v\n", couponReqInfo)

	return latestPacketData, couponState, couponAmount, couponReqInfo
}
