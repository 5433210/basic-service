package microservice

import (
	"strings"

	"wailik.com/internal/pkg/log"
)

func GetNodeNameByNodeId(id string) string {
	sli := strings.Split(id, "/")
	log.Debugf("%+v", len(sli))

	nodeName := sli[len(sli)-2]
	log.Debugf("%+v", nodeName)

	return nodeName
}
