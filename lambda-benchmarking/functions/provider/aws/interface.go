package aws

import (
	"fmt"
	"functions/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

//Instance is an object used to interact with AWS through the methods it exports.
type Instance struct {
	userID       string
	username     string
	user         string
	familiarName string
	region       string
	cloneAPIID   string
	stage        string
}

//Initialize will create a new AWS Instance to interact with.
func Initialize() *Instance {
	return &Instance{
		userID:       "335329526041",
		username:     "theodor",
		user:         "arn:aws:iam::335329526041:mfa/theodor",
		familiarName: "benchmarking",
		region:       "us-west-1",
		cloneAPIID:   "hjnwqihyo1",
		stage:        "prod",
	}
}

func (lambda Instance) getAPIID(i int) string {
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "get-rest-apis", "--query",
		fmt.Sprintf("items[?name==`%s-API-%v`].id", lambda.familiarName, i), "--output", "text",
		"--region", lambda.region)
	apiID := util.RunCommandAndLog(cmd)
	apiID, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(apiID), `\n`, ""))
	log.Infof("API ID of %s-API-%v is %s", lambda.familiarName, i, apiID)
	return apiID
}
