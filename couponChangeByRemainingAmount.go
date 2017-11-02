package main

import (
	"fmt"
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
	if !silent || debug {
		fmt.Println("Current coupon state and amount:")
		fmt.Printf("%+v %+v\n\n", couponState, couponAmount)
	}

	// Get today's date
	today := time.Now().Day()
	if debug {
		fmt.Println("Today's day:")
		fmt.Printf("%+v\n\n", today)
	}
	// Calculate remaining coupon quotas
	couponAmountQuota := []int{}
	for i := 1; i <= len(couponState); i++ {
		couponAmountQuota = append(couponAmountQuota,
			config.Mio.StartingCouponAmount-
				((today-1)*config.Mio.MaxDailyAmount*len(couponState))-
				(i*config.Mio.MaxDailyAmount))
	}
	if !silent || debug {
		fmt.Println("Remaining coupon amount quota (in MB):")
		fmt.Printf("%+v\n\n", couponAmountQuota)
	}

	// Extract latest packet usage data
	latestPacketDataSortedByAmount := extractLatestPacketDataSortedByAmount(packetData)
	if !silent || debug {
		fmt.Println("Latest packet usage (in MB):")
		fmt.Printf("%+v\n\n", latestPacketDataSortedByAmount)
	}
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
	fmt.Println("Coupon request info:")
	fmt.Printf("%v\n\n", couponReqInfo)

	return latestPacketData, couponState, couponAmount, couponReqInfo
}
