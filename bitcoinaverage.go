package bitcoinaverage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	publicKey string = "OTk4NDczMDEyYTg3NGFlMDkxNzdiYTkyYTljOGMxZTU"
	secretKey string = "MzAwZDAxYWY0M2RmNGE4MTg2NTk3YmY5NDVjZjQ1MTM4YTYxMWNlM2U0NDk0ZmRjODQ1N2RhMTlmNzA2OGNkNw"
	url       string = "https://apiv2.bitcoinaverage.com/"
)

//CurrentPrice parses JSON response
type CurrentPrice struct {
	Success bool    `json:"success"`
	Time    string  `json:"time"`
	Price   float64 `json:"price"`
}

//Convert takes two currencies and returns a CurrentPrice struct
func Convert(from string, to string) CurrentPrice {
	sig := generateHeader()
	action := "convert/"
	amount := "1"

	client := &http.Client{}
	req, e := http.NewRequest("GET", url+action+"global?from="+from+"&to="+to+"&amount="+amount, nil)
	checkErr(e)
	req.Header.Set("X-signature", sig)
	resp, e := client.Do(req)
	checkErr(e)

	//read Body and Unmarshal response into struct
	body, e := ioutil.ReadAll(resp.Body)
	checkErr(e)
	cp := new(CurrentPrice)
	e = json.Unmarshal([]byte(body), &cp)
	checkErr(e)
	return *cp
}

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func generateDigest(p string) string {
	m := []byte(p)
	k := []byte(secretKey)
	mac := hmac.New(sha256.New, k)
	mac.Write(m)
	return hex.EncodeToString(mac.Sum(nil))
}

func generateHeader() string {
	payload := generatePayload()
	digest := generateDigest(payload)
	return payload + "." + digest
}

func generatePayload() string {
	timestamp := time.Now().Unix()
	return strconv.Itoa(int(timestamp)) + "." + publicKey
}
