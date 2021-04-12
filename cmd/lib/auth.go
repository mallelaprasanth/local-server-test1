package lib

import (
	"strings"

	"github.com/gin-gonic/gin"
	// arm "github.com/synspective/syns-arms/pkg/sdk"
	// "github.com/synspective/syns-arms/pkg/sdk/customer"
	// "google.golang.org/protobuf/types/known/wrapperspb"
)

// FilterUUIDsByService filters the given UUID respective to the defined service
func FilterUUIDsByService(c *gin.Context, service string, uuids []string) (newlist []string) {
	// ok, authtoken := CheckForBearerToken(c.GetHeader("Authorization"))
	// if !ok {
	// 	return
	// }

	// res, err := arm.Customer().GetUUIDs(arm.AuthContext(authtoken), &customer.UUIDRequest{
	// 	Service: wrapperspb.String(service),
	// })

	// if err != nil {
	// 	log.Printf("Failed to fetch UUIDs from ARM: %s", err)
	// 	return
	// }

	// // We need to cross-check the list, not just return back what ARM returned
	// for _, stru := range uuids {
	// 	uuid, err := strconv.ParseUint(stru, 10, 64)
	// 	if err != nil {
	// 		log.Fatalf("Failed parse UUID: %s", err)
	// 	} else {
	// 		for _, u := range res.Uuids {
	// 			if u == uuid {
	// 				newlist = append(newlist, stru)
	// 			}
	// 		}
	// 	}
	// }

	return
}

// CheckForBearerToken check whether the authorization header is available or not
func CheckForBearerToken(header string) (tokenexist bool, token string) {
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
		token = splitted[1]
		tokenexist = true
	}

	return tokenexist, token
}
