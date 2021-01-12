package util

import (
	"encoding/json"
	"net/http"
)

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Remove(slice []string, items []string) []string {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	for _, item := range items {
		_, ok := set[item]
		if ok {
			delete(set, item)
		}
	}

	var newSlice []string
	for k, _ := range set {
		newSlice = append(newSlice, k)
	}
	return newSlice
}

func SetResponse(res http.ResponseWriter, outString string, outJson interface{}, status int) http.ResponseWriter {

	//set Cors
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Max-Age", "3628800")
	res.Header().Set("Access-Control-Expose-Headers", "Content-Type, X-Requested-With, Accept, Authorization, Referer, User-Agent")

	//set StatusCode
	res.WriteHeader(status)

	//set Out
	if outJson != nil {
		res.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(outJson)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.Write(js)
	} else {
		res.Write([]byte(outString))
	}
	return res
}
