package pesapal

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"swiftbay-api/common/utils"
	"time"

	"github.com/beevik/etree"

	"github.com/shopspring/decimal"
)

func generatePostXML(details map[string]interface{}) (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)
	orderInfo := doc.CreateElement("PesapalDirectOrderInfo")
	orderInfo.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	orderInfo.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	orderInfo.CreateAttr("Amount", details["amount"].(decimal.Decimal).String())
	orderInfo.CreateAttr("Description", details["description"].(string))
	orderInfo.CreateAttr("Type", details["type"].(string))
	orderInfo.CreateAttr("Reference", details["reference"].(string))
	orderInfo.CreateAttr("FirstName", details["first_name"].(string))
	orderInfo.CreateAttr("LastName", details["last_name"].(string))
	orderInfo.CreateAttr("Email", details["email"].(string))
	orderInfo.CreateAttr("PhoneNumber", "")
	orderInfo.CreateAttr("Currency", details["currency"].(string))
	orderInfo.CreateAttr("xmlns", "http://www.pesapal.com")
	buf := new(bytes.Buffer)
	doc.WriteTo(buf)
	return buf.String(), nil
}

func setPostParameters(callbackURL string, consumerKey string, postXML string) url.Values {
	timestamp := time.Now().Unix()
	params := url.Values{}
	params.Add("oauth_callback", callbackURL)
	params.Add("oauth_consumer_key", consumerKey)
	params.Add("oauth_nonce", fmt.Sprintf("%s%s", strconv.FormatInt(timestamp, 10), utils.RandStringBytesMaskImprSrc(12)))
	params.Add("oauth_signature_method", "HMAC-SHA1")
	params.Add("oauth_timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("oauth_version", "1.0")
	params.Add("pesapal_request_data", postXML)
	return params
}
