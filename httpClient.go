package main

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
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
		fmt.Println("Get request type error: ", reqType)
		err := fmt.Errorf("Get request type %s is not supported", reqType)
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-IIJmio-Developer", config.DeveloperId)
	req.Header.Add("X-IIJmio-Authorization", config.AccessToken)
	if err != nil {
		fmt.Println("HTTP GET request error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP GET request:")
		fmt.Printf("%+v\n\n", *req)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP GET response error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP GET response:")
		fmt.Printf("%+v\n\n", *resp)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTTP response body read error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP response body:")
		fmt.Printf("%+v\n\n", string(body))
	}

	return body, nil
}

func putCouponRequest(couponReqBytes []byte) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s", couponEndpoint)

	var bbuf bytes.Buffer
	bbuf.Write(couponReqBytes)
	req, err := http.NewRequest("PUT", url, &bbuf)
	req.Header.Add("X-IIJmio-Developer", config.DeveloperId)
	req.Header.Add("X-IIJmio-Authorization", config.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println("HTTP PUT request error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP PUT request:")
		fmt.Printf("%+v\n\n", *req)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP PUT response error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP PUT response:")
		fmt.Printf("%+v\n\n", *resp)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTTP body error: ", err)
		return nil, err
	}
	if debug == true {
		fmt.Printf("%s\n", "HTTP response body:")
		fmt.Printf("%+v\n\n", string(body))
	}

	return body, nil
}
