package utils

import "context"

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsInt(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetRequestId (ctx context.Context) string {
	var requestID string
	request := ctx.Value("requestID")
	if request != nil {
		requestID = request.(string)
	} else {
		requestID = "no_context"
	}
	return requestID
}

