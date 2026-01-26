package main

import (
	"fmt"
	"net/url"
)

func urlParsing() {
	//[scheme://][userinfo@]host[:port][/path][?query][#fragment]

	rawURL := "https://example.com:8080/path?query=param#fragment"

	parsedURL, err := url.Parse(rawURL)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}

	fmt.Println("Scheme:", parsedURL.Scheme)
	fmt.Println("User info:", parsedURL.User)
	fmt.Println("Host:", parsedURL.Host)
	fmt.Println("Port:", parsedURL.Port())
	fmt.Println("Raw query:", parsedURL.RawQuery)
	fmt.Println("Path:", parsedURL.Path)
	fmt.Println("RawPath:", parsedURL.Fragment)

	rawURL1 := "https://example.com:8080/path?name=John&age=30"

	parseURL1, err := url.Parse(rawURL1)
	if err != nil {
		fmt.Println("Error:" , err)
		return
	}

	queryParams := parseURL1.Query()

	fmt.Println("Params: ", queryParams)
	fmt.Println("Name: ", queryParams.Get("name"))
	fmt.Println("Name: ", queryParams.Get("age"))

	for key, value := range queryParams{
		fmt.Println(key, value)
	}

	//building an URL

	baseURL := &url.URL{
		Scheme: "https",
		Host: "example.com",
		Path: "/path",
	}

	query := baseURL.Query()
	query.Set("name", "John")
	query.Set("age", "30")
	baseURL.RawQuery = query.Encode()

	fmt.Println("Build URL: ", baseURL.String())

	values := url.Values{}

	values.Add("name", "Jane")
	values.Add("age", "30")
	values.Add("city", "London")
	values.Add("country", "UK")

	encodedQuery := values.Encode()

	baseURL1 := "https://example.com/search"
	fullURL := baseURL1 + "?" + encodedQuery
	fmt.Println(encodedQuery)
	fmt.Println(fullURL)
	

}