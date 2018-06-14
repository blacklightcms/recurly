package webhooks

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/blacklightcms/recurly"
)

// Webhook notification constants.
const (
	// Account notifications.
	NewAccount              = "new_account_notification"
	UpdatedAccount          = "updated_account_notification"
	CanceledAccount         = "canceled_account_notification"
	BillingInfoUpdated      = "billing_info_updated_notification"
	BillingInfoUpdateFailed = "billing_info_update_failed_notification"

	// Subscription notifications.
	NewSubscription      = "new_subscription_notification"
	UpdatedSubscription  = "updated_subscription_notification"
	RenewedSubscription  = "renewed_subscription_notification"
	ExpiredSubscription  = "expired_subscription_notification"
	CanceledSubscription = "canceled_subscription_notification"
	ReactivatedAccount   = "reactivated_account_notification"

	// Invoice notifications.
	NewInvoice        = "new_invoice_notification"
	PastDueInvoice    = "past_due_invoice_notification"
	ProcessingInvoice = "processing_invoice_notification"
	ClosedInvoice     = "closed_invoice_notification"

	// Payment notifications.
	SuccessfulPayment = "successful_payment_notification"
	FailedPayment     = "failed_payment_notification"
	VoidPayment       = "void_payment_notification"
	SuccessfulRefund  = "successful_refund_notification"
	ScheduledPayment  = "scheduled_payment_notification"
	ProcessingPayment = "processing_payment_notification"

	// Dunning Event notifications.
	NewDunningEvent = "new_dunning_event_notification"
)

type notificationName struct {
	XMLName xml.Name
}

// Account represents the account object sent in webhooks.
type Account struct {
	XMLName     xml.Name `xml:"account"`
	Code        string   `xml:"account_code,omitempty"`
	Username    string   `xml:"username,omitempty"`
	Email       string   `xml:"email,omitempty"`
	FirstName   string   `xml:"first_name,omitempty"`
	LastName    string   `xml:"last_name,omitempty"`
	CompanyName string   `xml:"company_name,omitempty"`
	Phone       string   `xml:"phone,omitempty"`
}

// Transaction represents the transaction object sent in webhooks.
type Transaction struct {
	XMLName           xml.Name         `xml:"transaction"`
	UUID              string           `xml:"id,omitempty"`
	InvoiceNumber     int              `xml:"invoice_number,omitempty"`
	SubscriptionUUID  string           `xml:"subscription_id,omitempty"`
	Action            string           `xml:"action,omitempty"`
	AmountInCents     int              `xml:"amount_in_cents,omitempty"`
	Status            string           `xml:"status,omitempty"`
	Message           string           `xml:"message,omitempty"`
	GatewayErrorCodes string           `xml:"gateway_error_codes,omitempty"`
	FailureType       string           `xml:"failure_type,omitempty"`
	Reference         string           `xml:"reference,omitempty"`
	Source            string           `xml:"source,omitempty"`
	Test              recurly.NullBool `xml:"test,omitempty"`
	Voidable          recurly.NullBool `xml:"voidable,omitempty"`
	Refundable        recurly.NullBool `xml:"refundable,omitempty"`
}

// Invoice represents the invoice object sent in webhooks.
type Invoice struct {
	XMLName             xml.Name         `xml:"invoice,omitempty"`
	SubscriptionUUID    string           `xml:"subscription_id,omitempty"`
	UUID                string           `xml:"uuid,omitempty"`
	State               string           `xml:"state,omitempty"`
	InvoiceNumberPrefix string           `xml:"invoice_number_prefix,omitempty"`
	InvoiceNumber       int              `xml:"invoice_number,omitempty"`
	PONumber            string           `xml:"po_number,omitempty"`
	VATNumber           string           `xml:"vat_number,omitempty"`
	TotalInCents        int              `xml:"total_in_cents,omitempty"`
	Currency            string           `xml:"currency,omitempty"`
	CreatedAt           recurly.NullTime `xml:"date,omitempty"`
	CreatedAtNew        recurly.NullTime `xml:"created_at,omitempty"`
	UpdatedAt           recurly.NullTime `xml:"updated_at,omitempty"`
	ClosedAt            recurly.NullTime `xml:"closed_at,omitempty"`
	NetTerms            recurly.NullInt  `xml:"net_terms,omitempty"`
	CollectionMethod    string           `xml:"collection_method,omitempty"`
}

// Transaction constants.
const (
	TransactionFailureTypeDeclined  = "declined"
	TransactionFailureTypeDuplicate = "duplicate_transaction"
)

// Account types.
type (
	// NewAccountNotification is sent when a customer creates a new account.
	// https://dev.recurly.com/page/webhooks#section-new-account
	NewAccountNotification struct {
		Account Account `xml:"account"`
	}

	// UpdatedAccountNotification is sent when a customer updates account information.
	// https://dev.recurly.com/page/webhooks#section-updated-account
	UpdatedAccountNotification struct {
		Account Account `xml:"account"`
	}

	// CanceledAccountNotification is sent when a customer closes their account.
	// https://dev.recurly.com/page/webhooks#section-closed-account
	CanceledAccountNotification struct {
		Account Account `xml:"account"`
	}

	// BillingInfoUpdatedNotification is sent when a customer updates or adds billing information.
	// https://dev.recurly.com/page/webhooks#section-updated-billing-information
	BillingInfoUpdatedNotification struct {
		Account Account `xml:"account"`
	}

	// BillingInfoUpdateFailedNotification is sent when a customer's billing update fails.
	// https://dev.recurly.com/page/webhooks#section-failed-billing-information-update
	BillingInfoUpdateFailedNotification struct {
		Account Account `xml:"account"`
	}
)

// Subscription types.
type (
	// NewSubscriptionNotification is sent when a new subscription is created.
	// https://dev.recurly.com/page/webhooks#section-new-subscription
	NewSubscriptionNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}

	// UpdatedSubscriptionNotification is sent when a subscription is upgraded or downgraded.
	// https://dev.recurly.com/page/webhooks#section-updated-subscription
	UpdatedSubscriptionNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}

	// RenewedSubscriptionNotification is sent when a subscription renew.
	// https://dev.recurly.com/page/webhooks#section-renewed-subscription
	RenewedSubscriptionNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}

	// ExpiredSubscriptionNotification is sent when a subscription is no longer valid.
	// https://dev.recurly.com/v2.4/page/webhooks#section-expired-subscription
	ExpiredSubscriptionNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}

	// CanceledSubscriptionNotification is sent when a subscription is canceled.
	// https://dev.recurly.com/page/webhooks#section-canceled-subscription
	CanceledSubscriptionNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}

	// ReactivatedAccountNotification is sent when a subscription is reactivated after having been canceled
	// https://dev.recurly.com/v2.6/page/webhooks#section-reactivated-subscription
	ReactivatedAccountNotification struct {
		Account      Account              `xml:"account"`
		Subscription recurly.Subscription `xml:"subscription"`
	}
)

// Invoice types.
type (
	// NewInvoiceNotification is sent when an invoice generated.
	// https://dev.recurly.com/page/webhooks#section-new-invoice
	NewInvoiceNotification struct {
		Account Account `xml:"account"`
		Invoice Invoice `xml:"invoice"`
	}

	// PastDueInvoiceNotification is sent when an invoice is past due.
	// https://dev.recurly.com/v2.4/page/webhooks#section-past-due-invoice
	PastDueInvoiceNotification struct {
		Account Account `xml:"account"`
		Invoice Invoice `xml:"invoice"`
	}

	// ProcessingInvoiceNotification is sent if an invoice is paid with ACH or a PayPal eCheck.
	// https://dev.recurly.com/page/webhooks#section-processing-invoice-automatic-only-for-ach-and-paypal-echeck-payments-
	ProcessingInvoiceNotification struct {
		Account Account `xml:"account"`
		Invoice Invoice `xml:"invoice"`
	}

	// ClosedInvoiceNotification is sent when an invoice is closed.
	// https://dev.recurly.com/page/webhooks#section-closed-invoice
	ClosedInvoiceNotification struct {
		Account Account `xml:"account"`
		Invoice Invoice `xml:"invoice"`
	}
)

// Payment types.
type (
	// SuccessfulPaymentNotification is sent when a payment is successful.
	// https://dev.recurly.com/v2.4/page/webhooks#section-successful-payment
	SuccessfulPaymentNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}

	// FailedPaymentNotification is sent when a payment fails.
	// https://dev.recurly.com/v2.4/page/webhooks#section-failed-payment
	FailedPaymentNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}

	// VoidPaymentNotification is sent when a successful payment is voided.
	// https://dev.recurly.com/page/webhooks#section-void-payment
	VoidPaymentNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}

	// SuccessfulRefundNotification is sent when an amount is refunded.
	// https://dev.recurly.com/page/webhooks#section-successful-refund
	SuccessfulRefundNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}

	// ScheduledPaymentNotification is sent when Recurly initiates an ACH payment from a customer entering payment or the renewal process.
	// https://dev.recurly.com/page/webhooks#section-scheduled-payment-only-for-ach-payments-
	ScheduledPaymentNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}

	// ProcessingPaymentNotification is sent when an ACH or PayPal eCheck payment moves from the scheduled state to the processing state.
	// https://dev.recurly.com/page/webhooks#section-processing-payment-only-for-ach-and-paypal-echeck-payments-
	ProcessingPaymentNotification struct {
		Account     Account     `xml:"account"`
		Transaction Transaction `xml:"transaction"`
	}
)

// Dunning Event types.
type (
	// NewDunningEventNotification is sent when an invoice enters and remains in dunning.
	// https://dev.recurly.com/page/webhooks#dunning-event-notifications
	NewDunningEventNotification struct {
		Account      Account              `xml:"account"`
		Invoice      Invoice              `xml:"invoice"`
		Subscription recurly.Subscription `xml:"subscription"`
		Transaction  Transaction          `xml:"transaction"`
	}
)

// ErrUnknownNotification is used when the incoming webhook does not match a
// predefined notification type. It implements the error interface.
type ErrUnknownNotification struct {
	name string
}

// Error implements the error interface.
func (e ErrUnknownNotification) Error() string {
	return fmt.Sprintf("unknown notification: %s", e.name)
}

// Name returns the name of the unknown notification.
func (e ErrUnknownNotification) Name() string {
	return e.name
}

// Parse parses an incoming webhook and returns the notification.
func Parse(r io.Reader) (interface{}, error) {
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	notification, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var n notificationName
	if err := xml.Unmarshal(notification, &n); err != nil {
		return nil, err
	}

	var dst interface{}
	switch n.XMLName.Local {
	case NewAccount:
		dst = &NewAccountNotification{}
	case UpdatedAccount:
		dst = &UpdatedAccountNotification{}
	case CanceledAccount:
		dst = &CanceledAccountNotification{}
	case BillingInfoUpdated:
		dst = &BillingInfoUpdatedNotification{}
	case BillingInfoUpdateFailed:
		dst = &BillingInfoUpdateFailedNotification{}
	case NewSubscription:
		dst = &NewSubscriptionNotification{}
	case UpdatedSubscription:
		dst = &UpdatedSubscriptionNotification{}
	case RenewedSubscription:
		dst = &RenewedSubscriptionNotification{}
	case ExpiredSubscription:
		dst = &ExpiredSubscriptionNotification{}
	case CanceledSubscription:
		dst = &CanceledSubscriptionNotification{}
	case ReactivatedAccount:
		dst = &ReactivatedAccountNotification{}
	case NewInvoice:
		dst = &NewInvoiceNotification{}
	case PastDueInvoice:
		dst = &PastDueInvoiceNotification{}
	case ProcessingInvoice:
		dst = &ProcessingInvoiceNotification{}
	case ClosedInvoice:
		dst = &ClosedInvoiceNotification{}
	case SuccessfulPayment:
		dst = &SuccessfulPaymentNotification{}
	case FailedPayment:
		dst = &FailedPaymentNotification{}
	case VoidPayment:
		dst = &VoidPaymentNotification{}
	case SuccessfulRefund:
		dst = &SuccessfulRefundNotification{}
	case ScheduledPayment:
		dst = &ScheduledPaymentNotification{}
	case ProcessingPayment:
		dst = &ProcessingPaymentNotification{}
	case NewDunningEvent:
		dst = &NewDunningEventNotification{}
	default:
		return nil, ErrUnknownNotification{name: n.XMLName.Local}
	}

	if err := xml.Unmarshal(notification, dst); err != nil {
		return nil, err
	}

	return dst, nil
}
