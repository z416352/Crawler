package request_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"

	"github.com/z416352/Database-api/pkg/responses"
)

// const base_mongodb_url = "http://localhost:8080/mongodb"
// const base_prices_url = "http://localhost:8080/prices"

var base_prices_url string
var base_mongodb_url string

func init() {
	db_API_ServiceName := os.Getenv("db_API_ServiceName")
	db_API_ServicePort := os.Getenv("db_API_ServicePort")
	db_API_Namespace := os.Getenv("db_API_Namespace")

	base_prices_url = fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/prices", db_API_ServiceName, db_API_Namespace, db_API_ServicePort)
	base_mongodb_url = fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/mongodb", db_API_ServiceName, db_API_Namespace, db_API_ServicePort)

	if os.Getenv("TEST") == "1" {
		base_mongodb_url = "http://localhost:8080/mongodb"
		base_prices_url = "http://localhost:8080/prices"
	}
}

// prices_url
func Get_GetNewestData(symbol, timeframe, num string) []*apiservice.BinanceAPI_Kline {
	url := fmt.Sprintf("%s/%s/%s/%s", base_prices_url, symbol, timeframe, num)
	logger.CrawlerLog.Debugf("Get method url: %s", url)

	response, err := http.Get(url)
	if err != nil {
		logger.CrawlerLog.Errorf("err")
	}
	defer response.Body.Close()

	// Read the response body
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)

	// Json string to data structure
	res := responses.UserResponse{}
	json.Unmarshal(body.Bytes(), &res)

	klines := []*apiservice.BinanceAPI_Kline{}
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		logger.CrawlerLog.Errorf("API returned an error: {'status': %d, 'message': %s}", res.Status, res.Message)
	} else {
		err = convertInterfaceToStruct(res.Data["klines"], &klines)
		if err != nil {
			logger.CrawlerLog.Errorf("err: %v", err)
		}
	}

	// logger.CrawlerLog.Infof("kline: %v", kline)

	return klines
}

// prices_url
func Post_InsertACrawlerData(kline_detail *apiservice.BinanceTypeRequestDetail, symbol string, timeframe string) {
	url := base_prices_url
	logger.CrawlerLog.Debugf("Post method url: %s", url)

	jsonString, err := json.Marshal(kline_detail)
	if err != nil {
		logger.CrawlerLog.Errorf("Error:", err)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonString))

	if err != nil {
		logger.CrawlerLog.Errorf("Error sending request: %v", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)

	// Json string to data structure
	res := responses.UserResponse{}
	if err = json.Unmarshal(body.Bytes(), &res); err != nil {
		logger.CrawlerLog.Errorf("%v", err)
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		logger.CrawlerLog.Errorf("API returned an error: {'status': %d, 'message': %s}", res.Status, res.Message)
	}
	// else {
	// 	// Print the response body
	// 	logger.CrawlerLog.Infof("API response: {'status': %d, 'message': %s}", res.Status, res.Message)
	// }
}

// mongodb_url
func Get_DBExistTimeframes(dbName string) []string {
	url := fmt.Sprintf("%s/%s", base_mongodb_url, dbName)
	logger.CrawlerLog.Debugf("Get method url: %s", url)

	response, err := http.Get(url)
	if err != nil {
		logger.CrawlerLog.Errorf("err")
	}
	defer response.Body.Close()

	// Read the response body
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)

	// Json string to data structure
	res := responses.UserResponse{}
	json.Unmarshal(body.Bytes(), &res)

	timeframes := []string{}

	// Check the response status code
	switch response.StatusCode {
	case http.StatusOK: // Get collections list
		// logger.CrawlerLog.Infof("API response: {'status': %d, 'message': %s}", res.Status, res.Message)

		err = convertInterfaceToStruct(res.Data["timeframes"], &timeframes)
		if err != nil {
			logger.CrawlerLog.Errorf("err: %v", err)
		}

	case http.StatusNotFound: // database not found
		// logger.CrawlerLog.Infof("Database '%s' isn't exist.", dbName)
		timeframes = []string{}

	case http.StatusInternalServerError:
		logger.CrawlerLog.Errorf("API returned an error: {'status': %d, 'message': %s}", res.Status, res.Message)
	}

	return timeframes
}
