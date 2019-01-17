package pesapal

import (
	"fmt"
	"net/url"
	"strconv"
	"swiftbay-api/common/utils"
	"time"
)

func setDetailsParameters(consumerKey string, merchantReference string, transactionTrackingID string) url.Values {
	timestamp := time.Now().Unix()
	params := url.Values{}
	params.Add("oauth_consumer_key", consumerKey)
	params.Add("oauth_nonce", fmt.Sprintf("%s%s", strconv.FormatInt(timestamp, 10), utils.RandStringBytesMaskImprSrc(12)))
	params.Add("oauth_signature_method", "HMAC-SHA1")
	params.Add("oauth_timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("oauth_version", "1.0")
	params.Add("pesapal_merchant_reference", merchantReference)
	params.Add("pesapal_transaction_tracking_id", transactionTrackingID)
	return params
}
