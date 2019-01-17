package pesapal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
)

func generateEncodedParamsQueryString(params url.Values) string {
	return params.Encode()
}

func generateSignableEncodedParamsQueryString(params url.Values) string {
	params.Del("oauth_signature")
	return generateEncodedParamsQueryString(params)
}
func generateSignatureBaseString(
	httpMethod string,
	base string,
	path string,
	params url.Values) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("Malformed URL: %s", err.Error())
	}
	baseURL.Path += path
	encodedQuery := url.QueryEscape(generateSignableEncodedParamsQueryString(params))
	//logrus.Infof("encodedquery: %s", encodedQuery)
	return fmt.Sprintf("%s&%s&%s", httpMethod, url.QueryEscape(baseURL.String()), encodedQuery), nil
}

func generateSigningKey(consumerSecret string, tokenSecret string) string {
	encodedConsumerSecret := url.QueryEscape(consumerSecret)
	// encodedTokenSecret := ""
	/*if len(tokenSecret) > 0 {
		encodedTokenSecret = url.QueryEscape(tokenSecret)
	}*/

	return fmt.Sprintf("%s&", encodedConsumerSecret)
}

func generateOauthSignature(httpMethod string, base string, path string, params url.Values, consumerSecret string, tokenSecret string) (string, error) {
	signatureBaseString, err := generateSignatureBaseString(httpMethod, base, path, params)
	if err != nil {
		return "", err
	}
	signingKey := generateSigningKey(consumerSecret, tokenSecret)
	h := hmac.New(sha1.New, []byte(signingKey))
	h.Write([]byte(signatureBaseString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
