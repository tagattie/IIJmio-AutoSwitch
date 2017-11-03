package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getData(reqType string) ([]byte, error) {
	var url string
	// Get packet usage data
	client := &http.Client{}
	switch reqType {
	case "packet":
		url = fmt.Sprintf("%s", packetEndpoint)
	case "coupon":
		url = fmt.Sprintf("%s", couponEndpoint)
	default:
		log.Println("Get request type error: ", reqType)
		err := fmt.Errorf("Get request type %+v is not supported", reqType)
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-IIJmio-Developer", config.Mio.DeveloperId)
	req.Header.Add("X-IIJmio-Authorization", config.Mio.AccessToken)
	if err != nil {
		log.Println("HTTP GET request error: ", err)
		return nil, err
	}
	log.Println("HTTP GET request:")
	log.Printf("%+v\n", *req)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP GET response error: ", err)
		return nil, err
	}
	log.Println("HTTP GET response:")
	log.Printf("%+v\n", *resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("HTTP response body read error: ", err)
		return nil, err
	}
	log.Println("HTTP response body:")
	log.Printf("%+v\n", string(body))

	return body, nil
}

func putCouponRequest(couponReqBytes []byte) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s", couponEndpoint)

	var bbuf bytes.Buffer
	bbuf.Write(couponReqBytes)
	req, err := http.NewRequest("PUT", url, &bbuf)
	req.Header.Add("X-IIJmio-Developer", config.Mio.DeveloperId)
	req.Header.Add("X-IIJmio-Authorization", config.Mio.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("HTTP PUT request error: ", err)
		return nil, err
	}
	log.Println("HTTP PUT request:")
	log.Printf("%+v\n", *req)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP PUT response error: ", err)
		return nil, err
	}
	log.Println("HTTP PUT response:")
	log.Printf("%+v\n", *resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("HTTP body error: ", err)
		return nil, err
	}
	log.Println("HTTP response body:")
	log.Printf("%+v\n", string(body))

	return body, nil
}
