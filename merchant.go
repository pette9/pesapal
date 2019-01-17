package pesapal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

//IPesapalMerchant ...
type IPesapalMerchant interface {
	GenerateOrderURL(callbackURL string, orderDetails map[string]interface{}) (string, error)
	QueryPaymentStatus(merchantReference string, transactionTrackingID string) (map[string]string, error)
	IPNListener(merchantReference string, transactionTrackingID string) (map[string]interface{}, error)
	QueryPaymentDetails(merchantReference string, transactionTrackingID string) (map[string]string, error)
}

//Merchant ...
type Merchant struct{}

//GenerateOrderURL ...
func (merchant *Merchant) GenerateOrderURL(callbackURL string, orderDetails map[string]interface{}) (string, error) {
	postXML, err := generatePostXML(orderDetails)
	if err != nil {
		return "", err
	}
	params := setPostParameters(callbackURL, viper.GetString("pesapal_consumer_key"), postXML)
	oauthSignature, err := generateOauthSignature(
		"GET",
		viper.GetString("pesapal_url"),
		viper.GetString("pesapal_direct_order_v4_path"),
		params,
		viper.GetString("pesapal_consumer_secret"),
		"",
	)
	if err != nil {
		return "", err
	}
	params.Add("oauth_signature", oauthSignature)
	queryString := generateEncodedParamsQueryString(params)
	orderURL := fmt.Sprintf("%s%s?%s", viper.GetString("pesapal_url"), viper.GetString("pesapal_direct_order_v4_path"), queryString)
	return orderURL, nil
}

//QueryPaymentDetails ...
func (merchant *Merchant) QueryPaymentDetails(merchantReference string, transactionTrackingID string) (map[string]string, error) {
	var paymentDetails map[string]string
	params := setDetailsParameters(viper.GetString("pesapal_consumer_key"), merchantReference, transactionTrackingID)
	oauthSignature, err := generateOauthSignature(
		"GET",
		viper.GetString("pesapal_url"),
		viper.GetString("pesapal_payment_details_path"),
		params,
		viper.GetString("pesapal_consumer_secret"),
		"",
	)
	if err != nil {
		return paymentDetails, err
	}
	params.Add("oauth_signature", oauthSignature)
	queryString := generateEncodedParamsQueryString(params)
	queryPaymentDetailsURL := fmt.Sprintf("%s%s?%s", viper.GetString("pesapal_url"), viper.GetString("pesapal_payment_details_path"), queryString)
	client := &http.Client{}

	req, err := http.NewRequest("GET", queryPaymentDetailsURL, nil)

	if err != nil {
		return paymentDetails, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return paymentDetails, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return paymentDetails, err
	}
	responseString := strings.Split(string(body), "=")
	s := strings.Split(responseString[1], ",")
	paymentDetails = map[string]string{
		"method":                  s[1],
		"status":                  s[2],
		"merchant_reference":      s[3],
		"transaction_tracking_id": s[0],
	}
	return paymentDetails, nil
}

//QueryPaymentStatus ...
func (merchant *Merchant) QueryPaymentStatus(merchantReference string, transactionTrackingID string) (map[string]string, error) {
	var paymentStatus map[string]string
	params := setStatusParameters(viper.GetString("pesapal_consumer_key"), merchantReference, transactionTrackingID)
	oauthSignature, err := generateOauthSignature(
		"GET",
		viper.GetString("pesapal_url"),
		viper.GetString("pesapal_payment_status_path"),
		params,
		viper.GetString("pesapal_consumer_secret"),
		"",
	)
	if err != nil {
		return paymentStatus, err
	}
	params.Add("oauth_signature", oauthSignature)
	queryString := generateEncodedParamsQueryString(params)
	queryPaymentStatusURL := fmt.Sprintf("%s%s?%s", viper.GetString("pesapal_url"), viper.GetString("pesapal_payment_status_path"), queryString)
	client := &http.Client{}
	req, err := http.NewRequest("GET", queryPaymentStatusURL, nil)
	if err != nil {
		return paymentStatus, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return paymentStatus, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return paymentStatus, err
	}
	responseString := strings.Split(string(body), "=")
	paymentStatus = map[string]string{
		"status": responseString[1],
	}
	return paymentStatus, nil
}

//IPNListener ...
func (merchant *Merchant) IPNListener(merchantReference string, transactionTrackingID string) (map[string]interface{}, error) {
	notificationType := "CHANGE"
	var output map[string]interface{}
	status, err := merchant.QueryPaymentStatus(merchantReference, transactionTrackingID)
	if err != nil {
		return output, err
	}

	output = map[string]interface{}{"status": status["status"], "response": nil}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("pesapal_notification_type=%s", notificationType))
	sb.WriteString(fmt.Sprintf("&pesapal_transaction_tracking_id=%s", transactionTrackingID))
	sb.WriteString(fmt.Sprintf("&pesapal_merchant_reference=%s", merchantReference))

	if status["status"] == "COMPLETED" || status["status"] == "FAILED" {
		output["response"] = sb.String()
	}
	return output, nil
}
