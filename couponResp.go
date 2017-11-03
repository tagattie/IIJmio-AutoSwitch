package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Response to coupon request
type couponResp struct {
	ReturnCode string `json:"returnCode"`
}

func decodeCouponRespJSON(couponRespBytes []byte) (*couponResp, error) {
	var cresp couponResp

	if err := json.Unmarshal(couponRespBytes, &cresp); err != nil {
		log.Println("Coupon response JSON unmarshal error: ", err)
		return nil, err
	}

	if returnCode := cresp.ReturnCode; returnCode != "OK" {
		log.Println("Coupon response return code error: ", returnCode)
		err := fmt.Errorf("Coupon response return code is %+v", returnCode)
		return nil, err
	}

	return &cresp, nil
}
