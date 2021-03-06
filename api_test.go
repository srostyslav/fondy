package fondy

import (
	"github.com/satori/go.uuid"
	"testing"
	"time"
	"fmt"
)

var testData map[string]map[string]interface{} = map[string]map[string]interface{}{
	"merchant": map[string]interface{}{
		"id": 1396424,
		"secret": "test",
	},
	"checkout_data": map[string]interface{}{
		"amount": 100,
		"currency": "USD",
	},
	"order_data": map[string]interface{}{
		"order_id": 14290,
	},
	"order_full_data": map[string]interface{}{
		"amount": "100",
		"currency": "RUB",
	},
	"payment_p2p": map[string]interface{}{
		"receiver_card_number": "4444555566661111",
		"currency": "RUB",
		"amount": "100",
	},
	"payment_pcidss_non3ds": map[string]interface{}{
		"currency": "RUB",
		"amount": "100",
		"card_number": "4444555511116666",
		"cvv2": "123",
		"expiry_date": "1224",
	},
	"payment_pcidss_3ds": map[string]interface{}{
		"currency": "RUB",
		"amount": "100",
		"card_number": "4444555566661111",
		"cvv2": "123",
		"expiry_date": "1224",
	},  
}

var api *Api = NewApi(&ApiOptions{MerchantID: int64(testData["merchant"]["id"].(int)), SecretKey: testData["merchant"]["secret"].(string)})

func TestUrl(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		OrderDesc: "test",
		OrderID: uuid.NewV4().String(),
	}
	if url, err := api.CheckoutUrl(data); err != nil {
		t.Error(err.Error())
	} else if url == "" {
		t.Error("url is empty")
	}
}

func TestToken(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		OrderDesc: "test",
		OrderID: uuid.NewV4().String(),
	}

	if token, err := api.CheckoutToken(data); err != nil {
		t.Error(err.Error())
	} else if token == "" {
		t.Error("token is empty")
	}
}

func TestSubscription(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		RecurringData: &Recurring{
			StartTime: "2028-11-11",
			Amount: 234324,
			Every: 40,
			Period: "day",
		},
		OrderDesc: "test",
		OrderID: uuid.NewV4().String(),
	}

	if url, err := api.CheckoutSubscription(data); err != nil {
		t.Error(err.Error())
	} else if url == "" {
		t.Error("url is empty")
	}
}

func TestVerification(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		OrderDesc: "test",
		OrderID: uuid.NewV4().String(),
	}

	if url, err := api.CheckoutVerification(data); err != nil {
		t.Error(err.Error())
	} else if url == "" {
		t.Error("url is empty")
	}
}

func TestPcidss(t *testing.T) {
	data := &PCIDSSOneStep{
		Amount: testData["payment_pcidss_non3ds"]["amount"].(string),
		Currency: testData["payment_pcidss_non3ds"]["currency"].(string),
		CardNumber: testData["payment_pcidss_non3ds"]["card_number"].(string),
		Cvv2: testData["payment_pcidss_non3ds"]["cvv2"].(string),
		ExpiryDate: testData["payment_pcidss_non3ds"]["expiry_date"].(string),
		Preauth: "Y",
		RequiredRectoken: "Y",
		OrderID: uuid.NewV4().String(),
		OrderDesc: "Pay for order",
	}

	resp, err := api.PcidssStep1(data)
	if err != nil {
		t.Error(err.Error())
		return
	} else if resp["acs_url"].(string) == "" {
		t.Error("acs_url is empty")
		return
	}

	if resp, err := api.PcidssStep2(&PCIDSSTwoStep{OrderID: data.OrderID, Pareq: resp["pareq"].(string), Md: resp["md"].(string)}); err != nil {
		t.Error(err.Error())
	} else if id := resp["order_id"].(string); id == "" {
		t.Error("order_id is empty")
	}
}

func TestP2Pcredit(t *testing.T) {
	a := NewApi(&ApiOptions{MerchantID: 1000, SecretKey: "testcredit"})
	data := &P2Pcredit{
		ReceiverCardNumber: testData["payment_p2p"]["receiver_card_number"].(string),
		Currency: testData["payment_p2p"]["currency"].(string),
		Amount: testData["payment_p2p"]["amount"].(string),
		OrderID: uuid.NewV4().String(),
		OrderDesc: "Pay for order",
	}

	if resp, err := a.P2Pcredit(data); err != nil {
		t.Error(err.Error())
	} else if status := resp["order_status"].(string); status == "" {
		t.Error("status is empty")
	}
}

func TestReports(t *testing.T) {
	if resp, err := api.GetReports(time.Now().Local().Add(-time.Minute * time.Duration(280)), time.Now()); err != nil {
		t.Error(err.Error())
	} else if len(resp) == 0 {
		t.Error("not item")
	}
}

func CreateOrder(orderID string) (map[string]interface{}, error) {
	data := &PCIDSSOneStep{
		Amount: testData["payment_pcidss_non3ds"]["amount"].(string),
		Currency: testData["payment_pcidss_non3ds"]["currency"].(string),
		CardNumber: testData["payment_pcidss_non3ds"]["card_number"].(string),
		Cvv2: testData["payment_pcidss_non3ds"]["cvv2"].(string),
		ExpiryDate: testData["payment_pcidss_non3ds"]["expiry_date"].(string),
		Preauth: "Y",
		RequiredRectoken: "Y",
		OrderID: orderID,
		OrderDesc: "Pay for order",
	}

	return api.PcidssStep1(data)
}

func TestRecurring(t *testing.T) {
	orderID := uuid.NewV4().String()
	if resp, err := CreateOrder(orderID); err != nil {
		t.Error(err.Error())
	} else {
		data := &RecurringBody{
			OrderID: orderID,
			OrderDesc: "Pay for order",
			Amount: fmt.Sprint(testData["checkout_data"]["amount"]),
			Currency: testData["checkout_data"]["currency"].(string),
			Rectoken: resp["rectoken"].(string),
		}
		if order, err := api.Recurring(data); err != nil {
			t.Error(err.Error())
		} else if status := order["order_status"].(string); status != "approved" {
			t.Error("order is not approve: " + status)
		}
	}
}

func TestSettlement(t *testing.T) {
	orderID := uuid.NewV4().String()
	data := &Settlement{
		OperationID: orderID,
		Receiver: []Receiver{
			Receiver{
				Requisites: &Requisites{
					Amount: 500,
					MerchantID: 600001,
				},
				Type: "merchant",
			},
			Receiver{
				Requisites: &Requisites{
					Amount: 500,
					MerchantID: 700001,
				},
				Type: "merchant",
			},
		},
		Amount: testData["order_full_data"]["amount"].(string),
		Currency: testData["order_full_data"]["currency"].(string),
	}

	dataCapture := &Capture{
		OrderID: orderID,
		Amount: data.Amount,
		Currency: data.Currency,
	}

	if _, err := CreateOrder(orderID); err != nil {
		t.Error(err.Error())
		return
	}

	if status, err := api.Capture(dataCapture); err != nil {
		t.Error(err.Error())
		return
	} else if status == "" {
		t.Error("status is empty")
		return
	}

	if resp, err := api.Settlement(data); err != nil {
		t.Error(err.Error())
	} else if resp == 0 {
		t.Error("payment not found")
	}
}

func TestReverse(t *testing.T) {
	orderID := uuid.NewV4().String()
	if _, err := CreateOrder(orderID); err != nil {
		t.Error(err.Error())
		return
	}

	data := &Reverse{
		Amount: testData["order_full_data"]["amount"].(string),
		Currency: testData["order_full_data"]["currency"].(string),
		OrderID: orderID,
	}

	if resp, err := api.Reverse(data); err != nil {
		t.Error(err.Error())
	} else if resp == "" {
		t.Error("status is not found")
	}
}

func TestOrderStatus(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		OrderID: uuid.NewV4().String(),
		OrderDesc: "test",
	}
	if url, err := api.CheckoutUrl(data); err != nil {
		t.Error(err.Error())
	} else if url == "" {
		t.Error("url is empty")
	}

	if order, err := api.GetOrderStatus(data.OrderID); err != nil {
		t.Error(err.Error())
	} else if status := order["order_status"].(string); status != "created" {
		t.Error("status: " + status)
	}
}

func TestTransactionList(t *testing.T) {
	data := &Checkout{
		Amount: int64(testData["checkout_data"]["amount"].(int)),
		Currency: testData["checkout_data"]["currency"].(string),
		OrderID: uuid.NewV4().String(),
		OrderDesc: "test",
	}
	if url, err := api.CheckoutUrl(data); err != nil {
		t.Error(err.Error())
	} else if url == "" {
		t.Error("url is empty")
	}

	if _, err := api.TransactionList(data.OrderID); err != nil {
		t.Error(err.Error())
	} 
}
