package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Coupon-related data
type coupon struct {
	Volume int    `json:"volume"`
	Expire string `json:"expire"`
	Type   string `json:"type"`
}
type cDataHdoInfo struct {
	HdoServiceCode string   `json:"hdoServiceCode"`
	Number         string   `json:"number"`
	Iccid          string   `json:"iccid"`
	Regulation     bool     `json:"regulation"`
	Sms            bool     `json:"sms"`
	Voice          bool     `json:"voice"`
	CouponUse      bool     `json:"couponUse"`
	Coupon         []coupon `json:"coupon"`
}
type cDataHduInfo struct {
	HduServiceCode string   `json:"hduServiceCode"`
	Number         string   `json:"number"`
	Iccid          string   `json:"iccid"`
	Regulation     bool     `json:"regulation"`
	Sms            bool     `json:"sms"`
	Voice          bool     `json:"voice"`
	CouponUse      bool     `json:"couponUse"`
	Coupon         []coupon `json:"coupon"`
}
type couponInfo struct {
	HddServiceCode string         `json:"hddServiceCode"`
	Plan           string         `json:"plan"`
	HdoInfo        []cDataHdoInfo `json:"hdoInfo"`
	HduInfo        []cDataHduInfo `json:"hduInfo"`
	Coupon         []coupon       `json:"coupon"`
}
type couponData struct {
	ReturnCode string       `json:"returnCode"`
	CouponInfo []couponInfo `json:"couponInfo"`
}

func decodeCouponDataJSON(couponBytes []byte) (*couponData, error) {
	var cd couponData

	if err := json.Unmarshal(couponBytes, &cd); err != nil {
		log.Println("Coupon request JSON unmarshal error: ", err)
		return nil, err
	}

	if returnCode := cd.ReturnCode; returnCode != "OK" {
		log.Println("Coupon request return code error: ", returnCode)
		err := fmt.Errorf("Coupon request return code is %+v", returnCode)
		return nil, err
	}

	return &cd, nil
}

func getCouponStateAndAmount(couponData *couponData) (map[string]bool, int) {
	cState := make(map[string]bool)
	cAmount := 0

	cis := couponData.CouponInfo
	if cisLength := len(cis); cisLength > 0 {
		for i := range cis {
			hdois := cis[i].HdoInfo
			if hdoisLength := len(hdois); hdoisLength > 0 {
				for _, hdoInfo := range hdois {
					cState[hdoInfo.HdoServiceCode] = hdoInfo.CouponUse
				}
			}
			coupons := cis[i].Coupon
			if couponsLength := len(coupons); couponsLength > 0 {
				for _, coupon := range coupons {
					cAmount += coupon.Volume
				}
			}
		}
	}
	return cState, cAmount
}
