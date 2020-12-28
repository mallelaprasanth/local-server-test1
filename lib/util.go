package lib

import (
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func grpcErrorCodeToHTTP(c codes.Code) int {
	switch c {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func CheckForBearerToken(header string) (tokenexist bool, tokenHeader string) {
	if header == "" {
		tokenexist = false
		return tokenexist, ""
	}
	splitted := strings.Split(header, " ")
	if len(splitted) != 2 {
		tokenexist = false
	} else if splitted[0] != "Bearer" {
		tokenexist = false
	} else {
		tokenHeader = splitted[1]
		tokenexist = true
	}
	return tokenexist, tokenHeader
}
