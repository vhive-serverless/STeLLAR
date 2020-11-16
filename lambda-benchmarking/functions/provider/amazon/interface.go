package amazon

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
	appName    string
	region     string
	cloneAPIID string
	stage      string
}

//Initialize will create a new AWS Instance to interact with.
func Initialize() *Instance {
	return &Instance{
		appName:    "benchmarking",
		region:     "us-west-1",
		cloneAPIID: "hjnwqihyo1",
		stage:      "prod",
	}
}

func (amazon Instance) getAPIID(i int) string {
	cmd := exec.Command("/usr/local/bin/aws", "apigateway", "get-rest-apis", "--query",
		fmt.Sprintf("items[?name==`%s-API-%v`].id", amazon.appName, i), "--output", "text",
		"--region", amazon.region)
	apiID := util.RunCommandAndLog(cmd)
	apiID, _ = strconv.Unquote(strings.ReplaceAll(strconv.Quote(apiID), `\n`, ""))
	log.Infof("API ID of %s-API-%v is %s", amazon.appName, i, apiID)
	return apiID
}
