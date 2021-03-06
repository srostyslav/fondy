package fondy

type Recurring struct {
	StartTime				string 		`json:"start_time,omitempty"`  	// start date of the recurring order ('YYYY-MM-DD')
	Amount					int64 		`json:"amount,omitempty"`	// amount of the recurring order (int)
	Period					string 		`json:"period,omitempty"`	// period of the recurring order ('day', 'month', 'year')
	Every					int64 		`json:"every,omitempty"`	// frequency of the recurring order (int)
	Readonly				string 		`json:"readonly,omitempty"`	// possibility to change parameters of the recurring order by user ('y', 'n')
	State					string 		`json:"state,omitempty"`	// default state of the recurring order after opening url of the order ('y', 'n')
}

type Checkout struct {
	OrderID 				string 		`json:"order_id"`
	OrderDesc				string 		`json:"order_desc"`
	Amount					int64 		`json:"amount"`
	Currency				string 		`json:"currency"`
	ResponseUrl				string 		`json:"response_url,omitempty"`
	ServerCallbackUrl			string 		`json:"server_callback_url,omitempty"`
	PaymentSystems				string 		`json:"payment_systems,omitempty"`
	DefaultPaymentSystem			string 		`json:"default_payment_system,omitempty"`
	Lifetime				int64 		`json:"lifetime,omitempty"`
	MerchantData				string 		`json:"merchant_data,omitempty"`
	Preauth					string 		`json:"preauth,omitempty"`
	SenderEmail				string 		`json:"sender_email,omitempty"`
	Descriptor				string 		`json:"descriptor,omitempty"`
	Delayed					string 		`json:"delayed,omitempty"`
	Lang					string 		`json:"lang,omitempty"`
	ProductID 				string 		`json:"product_id,omitempty"`
	RequiredRectoken			string 		`json:"required_rectoken,omitempty"`
	Verification				string 		`json:"verification,omitempty"`
	VerificationType			string 		`json:"verification_type,omitempty"`
	Rectoken				string 		`json:"rectoken,omitempty"`
	ReceiverRectoken			string 		`json:"receiver_rectoken,omitempty"`
	DesignID 				int64 		`json:"design_id,omitempty"`
	Subscription				string 		`json:"subscription,omitempty"`
	SubscriptionCallbackUrl			string 		`json:"subscription_callback_url,omitempty"`
	RecurringData				*Recurring	`json:"recurring_data,omitempty"`
}

type Requisites struct {
	Amount					float64		`json:"amount"`
	SettlementDescription			string		`json:"settlement_description,omitempty"`
	MerchantID				int64		`json:"merchant_id,omitempty"`
	Okpo					int64		`json:"okpo,omitempty"`
	JurName					string		`json:"jur_name,omitempty"`
	Account					int64		`json:"account,omitempty"`
	Rectoken				string		`json:"rectoken,omitempty"`
	CardNumber				int64		`json:"card_number,omitempty"`
}

type Receiver struct {
	Type 					string		`json:"type"`
	Requisites				*Requisites	`json:"requisites"`
}

type Settlement struct {
	ServerCallbackUrl			string 		`json:"server_callback_url,omitempty"`
	Rectoken				string 		`json:"rectoken,omitempty"`
	Currency				string 		`json:"currency"`
	Amount					string 		`json:"amount"`
	OrderType				string 		`json:"order_type"`
	ResponseUrl				string 		`json:"response_url,omitempty"`
	OrderID 				string 		`json:"order_id"`
	OperationID 				string 		`json:"operation_id"`
	OrderDesc				string 		`json:"order_desc,omitempty"`
	Receiver				[]Receiver	`json:"receiver"`
}

type PCIDSSOneStep struct {
	OrderID 				string 			`json:"order_id,omitempty"`
	OrderDesc				string 			`json:"order_desc,omitempty"`
	Amount					string 			`json:"amount,omitempty"`
	Currency				string 			`json:"currency,omitempty"`
	CardNumber				string 			`json:"card_number,omitempty"`
	Cvv2					string 			`json:"cvv2,omitempty"`
	ExpiryDate				string 			`json:"expiry_date,omitempty"`
	ClientIP 				string 			`json:"client_ip,omitempty"`
	Container				string 			`json:"container,omitempty"`
	RequiredRectoken			string 			`json:"required_rectoken"`
	Preauth					string 			`json:"preauth"`
}

type PCIDSSTwoStep struct {
	OrderID 				string 			`json:"order_id,omitempty"`
	Pareq					string 			`json:"pareq"`
	Md					string 			`json:"md"`
}

type P2Pcredit struct {
	ReceiverCardNumber			string			`json:"receiver_card_number,omitempty"`
	ReceiverRectoken			string 			`json:"receiver_rectoken,omitempty"`
	OrderID 				string 			`json:"order_id"`
	OrderDesc				string 			`json:"order_desc"`
	Currency				string 			`json:"currency"`
	Amount					string 			`json:"amount"`
}

type RecurringBody struct {
	OrderID 				string 			`json:"order_id"`
	OrderDesc				string 			`json:"order_desc"`
	Currency				string 			`json:"currency"`
	Amount					string 			`json:"amount"`
	Rectoken				string 			`json:"rectoken"`
}

type Capture struct {
	OrderID 				string 			`json:"order_id"`
	Amount					string 			`json:"amount"`
	Currency				string 			`json:"currency"`
}

type Reverse struct {
	OrderID 				string 			`json:"order_id"`
	Amount					string 			`json:"amount"`
	Currency				string 			`json:"currency"`
	Comment					string 			`json:"comment"`
}
