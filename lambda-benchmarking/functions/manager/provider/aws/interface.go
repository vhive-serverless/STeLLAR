package aws

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	user        = "arn:aws:iam::335329526041:mfa/theodor"
	name        = "benchmarking"
	region      = "us-west-1"
	cloneAPIID  = "hjnwqihyo1"
	usagePlanID = "8hutzb"
	stage       = "prod"
)

type Interface struct{}

func (lambda Interface) GetAPIID(i int) string {
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "get-rest-apis", "--query",
		fmt.Sprintf("items[?name==`%s-API-%v`].id", name, i), "--output", "text",
		"--region", region)
	apiID := runCommandAndReturnOutput(cmd)
	apiID, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(apiID), `\n`, ""))
	log.Printf("API ID of %s-API-%v is %s", name, i, apiID)
	return apiID
}