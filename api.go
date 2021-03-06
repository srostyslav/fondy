package fondy

import (
	"github.com/satori/go.uuid"
	"encoding/base64"
	"encoding/json"
	"crypto/sha1"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"errors"
	"time"
	"fmt"
	"os"
	"io"
)

type Response struct {
	ResponseStatus	string 	`json:"response_status"`
	ErrorCode	int	`json:"error_code"`
	ErrorMessage	string 	`json:"error_message"`
}

func (r *Response) GetError() error {
	if r.ErrorCode > 0 {
		return errors.New(fmt.Sprintf("%d: %s", r.ErrorCode, r.ErrorMessage))
	}
	return nil
}

type ApiOptions struct {
	MerchantID		int64		// Merchant id numeric
	SecretKey		string		// Secret key string 	
	RequestType		string		// request type allowed json, xml, form
	ApiDomain		string 		// api domain
	ApiProtocol		string 		// allowed protocols 1.0, 2.0
}

type Api struct {
	Options			*ApiOptions
	ApiUrl 			string
}

func (a *Api) headers() map[string]string {
	return map[string]string{
		"User-Agent": "Go SDK",
		"Content-Type": "application/json; charset=utf-8",
	}
}

func (a *Api) ToB64(data interface{}) (string, error) {
	var output []byte

	switch v := data.(type) {
	case string:
		output = []byte(v)
	default:
		if value, err := json.Marshal(data); err != nil {
			return "", err
		} else {
			output = value
		}
	}

	return base64.StdEncoding.EncodeToString(output), nil
}

func (a *Api) GetSignature(data string) string {

	s, h := strings.Join([]string{a.Options.SecretKey, data}, "|"), sha1.New()
	io.WriteString(h, s)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (a *Api) prepereData(body interface{}) ([]byte, error) {
	output, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, err
	}

	if v, ok := data["merchant_id"]; !ok || v == "" {
		data["merchant_id"] = a.Options.MerchantID
	}

	if v, ok := data["reservation_data"]; ok {
		if value, err := a.ToB64(v); err != nil {
			return nil, err
		} else {
			data["reservation_data"] = value
		} 
	}

	if b64Data, err := a.ToB64(map[string]interface{}{"order": data}); err != nil {
		return nil, err
	} else {
		return json.Marshal(map[string]interface{}{
			"request": map[string]interface{}{
				"data": b64Data,
				"version": a.Options.ApiProtocol,
				"signature": a.GetSignature(b64Data),
			},
		})
	}
}

func (a *Api) CheckSignature(response map[string]interface{}) error {
	switch data := response["data"].(type) {
	case string:
		switch signature := response["signature"].(type) {
		case string:
			if sign := a.GetSignature(data); sign == signature {
				return nil
			} else {
				return errors.New(fmt.Sprintf("Signature does not match: %s != %s", sign, signature))				
			}
		}
	}
	switch err := response["error_message"].(type) {
	case string:
		return errors.New(fmt.Sprintf("%v: %s", response["error_code"], err))
		
	}
	return errors.New("Signature does not match")
}

func (a *Api) GetResponse(content []byte, obj interface{}, checkSignature bool) error {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return err
	}

	switch response := data["response"].(type)  {
	case map[string]interface{}:
		if checkSignature {
			if err := a.CheckSignature(response); err != nil {
				return err
			}
		}
		switch dataString := response["data"].(type) {
		case string:
			if sd, err := base64.StdEncoding.DecodeString(dataString); err != nil {
				return err
			} else {
				return json.Unmarshal(sd, &obj)
			}
		default:
			if output, err := json.Marshal(response); err != nil {
				return err
			} else {
				return json.Unmarshal(output, &obj)
			}
		}
	case []interface{}:
		if output, err := json.Marshal(response); err != nil {
			return err
		} else {
			return json.Unmarshal(output, &obj)
		}	
	}
	return errors.New(fmt.Sprintf("Response body is empty: %v", data))
}

func (a *Api) post(path string, body interface{}, obj interface{}, checkSignature bool) error {

	output, err := a.prepereData(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.ApiUrl + path, strings.NewReader(string(output)))
	if err != nil {
		return err
	}

	for k, v := range a.headers() {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		return err 
	} else {
		defer resp.Body.Close()
		
		if content, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		} else {

			if resp.StatusCode != 200 && resp.StatusCode != 201 {
				return errors.New(fmt.Sprintf("Response code is: %d; content: %s", resp.StatusCode, string(content)))
			}

			return a.GetResponse(content, obj, checkSignature)
		}
	}
}

func (a *Api) checkout(data *Checkout, typ string, checkSignature bool) (string, error) {

	var resp struct {
		Response
		Order struct {
			CheckoutUrl	string 	`json:"checkout_url"`
		} `json:"order"`
		Token 			string 	`json:"token"`
	}
	if err := a.post("/checkout/" + typ + "/", data, &resp, checkSignature); err != nil {
		return "", err
	} else if typ == "url" {
		return resp.Order.CheckoutUrl, resp.GetError()
	} else {
		return resp.Token, resp.GetError()
	}
}

func (a *Api) CheckoutUrl(data *Checkout) (string, error) {
	return a.checkout(data, "url", true)
}

func (a *Api) CheckoutToken(data *Checkout) (string, error) {
	return a.checkout(data, "token", false)
}

func (a *Api) CheckoutVerification(data *Checkout) (string, error) {
	data.Verification = "Y"
	if data.VerificationType == "" {
		data.VerificationType = "code"
	}
	return a.CheckoutUrl(data)
}

func (a *Api) CheckoutSubscription(data *Checkout) (string, error) {
	data.Verification = "Y"
	return a.CheckoutUrl(data)
}

func (a *Api) PcidssStep1(data *PCIDSSOneStep) (map[string]interface{}, error) {
	var resp map[string]interface{}

	if err := a.post("/3dsecure_step1/", data, &resp, true); err != nil {
		return resp, err
	} else if code, ok := resp["error_code"]; ok && code.(int) > 0 {
		return resp, errors.New(fmt.Sprintf("%d: %s", code, resp["error_message"]))
	}
	return resp, nil
}

func (a *Api) PcidssStep2(data *PCIDSSTwoStep) (map[string]interface{}, error) {
	var resp struct {
		Response
		Order map[string]interface{}	`json:"order"`
	}
	
	if err := a.post("/3dsecure_step2/", data, &resp, true); err != nil {
		return resp.Order, err
	} else {
		return resp.Order, resp.GetError()
	}
}

func (a *Api) P2Pcredit(data *P2Pcredit) (map[string]interface{}, error) {
	var resp struct {
		Response
		Order map[string]interface{}	`json:"order"`
	}
	
	if err := a.post("/p2pcredit/", data, &resp, true); err != nil {
		return resp.Order, err
	} else {
		return resp.Order, resp.GetError()
	}
}

func (a *Api) GetReports(dateFrom, dateTo time.Time) ([]map[string]interface{}, error) {
	var resp []map[string]interface{}
	
	dateFmt := "02.01.2006 15:04:05"
	data := map[string]interface{}{
		"date_from": dateFrom.Format(dateFmt),
		"date_to": dateTo.Format(dateFmt),
	} 

	if err := a.post("/reports/", data, &resp, false); err != nil {
		return resp, err
	}
	return resp, nil
}

func (a *Api) Recurring(data *RecurringBody) (map[string]interface{}, error) {
	var resp struct {
		Response
		Order map[string]interface{}	`json:"order"`
	}
	
	if err := a.post("/recurring/", data, &resp, true); err != nil {
		return resp.Order, err
	} else {
		return resp.Order, resp.GetError()
	}
}

func (a *Api) Settlement(data *Settlement) (int64, error) {

	if data.OrderID == "" {
		data.OrderID = uuid.NewV4().String()
	}

	var resp struct {
		Response
		Order struct {
			PaymentID  int64  `json:"payment_id"`
		} `json:"order"`
	}
	if err := a.post("/settlement/", data, &resp, true); err != nil {
		return 0, err
	} else {
		return resp.Order.PaymentID, resp.GetError()
	}
}

func (a *Api) Capture(data *Capture) (string, error) {
	var resp struct {
		Response
		Order struct {
			CaptureStatus 	string 	`json:"capture_status"`
		} `json:"order"`
	}
	if err := a.post("/capture/", data, &resp, true); err != nil {
		return "", err
	} else {
		return resp.Order.CaptureStatus, resp.GetError()
	}
}

func (a *Api) Reverse(data *Reverse) (string, error) {
	var resp struct {
		Response
		Order struct {
			ReverseStatus 	string 	`json:"reverse_status"`
		} `json:"order"`
	}
	if err := a.post("/reverse/order_id/", data, &resp, true); err != nil {
		return "", err
	} else {
		return resp.Order.ReverseStatus, resp.GetError()
	}
}

func (a *Api) GetOrderStatus(orderID string) (map[string]interface{}, error) {
	var resp struct {
		Response
		Order map[string]interface{} `json:"order"`
	}
	if err := a.post("/status/order_id/", map[string]interface{}{"order_id": orderID}, &resp, true); err != nil {
		return resp.Order, err
	} else {
		return resp.Order, resp.GetError()
	}
}

func (a *Api) TransactionList(orderID string) ([]map[string]interface{}, error) {
	var resp []map[string]interface{}
	if err := a.post("/transaction_list/", map[string]interface{}{"order_id": orderID}, &resp, true); err != nil {
		return resp, err
	} else {
		return resp, nil
	}
}

func (a *Api) AtolLogs(orderID string) (interface{}, error) {
	var resp struct {
		Response
		Order interface{} `json:"order"`
	}
	if err := a.post("/get_atol_logs/", map[string]interface{}{"order_id": orderID}, &resp, true); err != nil {
		return resp.Order, err
	} else {
		return resp.Order, resp.GetError()
	}
}

func NewApi(options *ApiOptions) *Api {

	if options.RequestType == "" {
		options.RequestType = "json"
	}
	
	if options.MerchantID == 0 || options.SecretKey == "" {
		if id, err := strconv.Atoi(os.Getenv("CLOUDIPSP_MERCHANT_ID")); err != nil {
			panic("Incorrect 'CLOUDIPSP_MERCHANT_ID' env variable: " + err.Error())
		} else {
			options.MerchantID = int64(id)
		}
		options.SecretKey = os.Getenv("CLOUDIPSP_SECRETKEY")
	}

	if options.ApiDomain == "" {
		options.ApiDomain = "api.fondy.eu"
	}

	if options.ApiProtocol == "" {
		options.ApiProtocol = "2.0"
	}

	if options.ApiProtocol != "2.0" {
		panic("Incorrect protocol version")
	}

	if options.RequestType != "json" {
		panic("Only 'json' encoding allowed")
	}

	return &Api{Options: options, ApiUrl: fmt.Sprintf("https://%s/api", options.ApiDomain)}
}
