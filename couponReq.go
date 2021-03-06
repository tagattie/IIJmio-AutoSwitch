package main

import (
	"encoding/json"
	"log"
	"strings"
)

// Coupon request-related data
type cReqHdoInfo struct {
	HdoServiceCode string `json:"hdoServiceCode"`
	CouponUse      bool   `json:"couponUse"`
}
type cReqHduInfo struct {
	HduServiceCode string `json:"hduServiceCode"`
	CouponUse      bool   `json:"couponUse"`
}
type cReqCouponInfo struct {
	HdoInfo []cReqHdoInfo `json:"hdoInfo"`
	HduInfo []cReqHduInfo `json:"hduInfo"`
}
type couponReq struct {
	CouponInfo []cReqCouponInfo `json:"couponInfo"`
}

func encodeCouponReqJSON(couponReqInfo map[string]bool) ([]byte, error) {
	creq := couponReq{}

	crhdois := []cReqHdoInfo{}
	crhduis := []cReqHduInfo{}
	for k, v := range couponReqInfo {
		if strings.Index(k, "hdo") == 0 {
			crhdoi := cReqHdoInfo{
				HdoServiceCode: k,
				CouponUse:      v,
			}
			crhdois = append(crhdois, crhdoi)
		} else if strings.Index(k, "hdu") == 0 {
			crhdui := cReqHduInfo{
				HduServiceCode: k,
				CouponUse:      v,
			}
			crhduis = append(crhduis, crhdui)
		} else {
			log.Println("Service code error: ", k)
		}
	}
	crci := cReqCouponInfo{
		HdoInfo: crhdois,
		HduInfo: crhduis,
	}
	crcis := []cReqCouponInfo{}
	crcis = append(crcis, crci)
	creq.CouponInfo = crcis

	couponReqBytes, err := json.Marshal(creq)
	if err != nil {
		log.Println("Coupon request JSON marshal error: ", err)
		return nil, err
	}

	return couponReqBytes, nil
}
