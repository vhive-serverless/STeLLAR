package aws

// Security (unused at the moment)
//type Response struct {
//	Credentials struct {
//		AccessKeyID     string    `json:"AccessKeyId"`
//		SecretAccessKey string    `json:"SecretAccessKey"`
//		SessionToken    string    `json:"SessionToken"`
//		Expiration      time.Time `json:"Expiration"`
//	} `json:"Credentials"`
//}
//
//func (lambda Instance) setSessionToken() {
//	_, isSet := os.LookupEnv("AWS_SESSION_TOKEN")
//	if isSet {
//		return
//	}
//
//	if err := os.Unsetenv("AWS_SESSION_TOKEN"); err != nil {
//		log.Fatal(err.Error())
//	}
//
//	fmt.Print("Please enter virtual MFA device token code (ignore if not set): ")
//	var tokenCode int
//	_, err := fmt.Scanf("%d", &tokenCode)
//	if err != nil {
//		log.Fatalf("Could not parse virtual MFA device token code.")
//	}
//
//	cmd := exec.Command("/usr/local/bin/aws", "sts", "get-session-token",
//		"--serial-number", lambda.user, "--token-code", strconv.Itoa(tokenCode))
//	responseJSON := util.RunCommandAndLog(cmd)
//
//	var response Response
//	if err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
//		log.Fatalf(err.Error())
//	}
//
//	if err := os.Setenv("AWS_SESSION_TOKEN", response.Credentials.SessionToken); err != nil {
//		log.Fatal(err.Error())
//	}
//}
