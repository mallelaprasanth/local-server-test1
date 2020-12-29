package lib

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	arm "github.com/synspective/syns-arms/pkg/sdk"
	"github.com/synspective/syns-arms/pkg/sdk/customer"
	"google.golang.org/grpc/codes"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// return error codes
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

// check whether the authorization header is available or not
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

// check whether uuid is accessible for user or not through arm connect
func CanAccessUUID(struuid string, authtoken string) (canAccess bool) {
	uuid, err := strconv.ParseUint(struuid, 10, 64)
	if err != nil {
		canAccess = false
		return canAccess
	}
	res, err := arm.Customer().GetUUIDs(arm.AuthContext(authtoken), &customer.UUIDRequest{})
	if err != nil {
		canAccess = false
		return canAccess
	}
	for _, u := range res.Uuids {
		if u == uuid {
			canAccess = true
			break
		}
	}
	return canAccess
}
