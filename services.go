package recurly

import (
	"bytes"
	"time"
)

// Params are used to send parameters with the request.
type Params map[string]interface{}

// AccountsService represents the interactions available for accounts.
type AccountsService interface {
	List(params Params) (*Response, []Account, error)
	Get(code string) (*Response, *Account, error)
	LookupAccountBalance(code string) (*Response, *AccountBalance, error)
	Create(a Account) (*Response, *Account, error)
	Update(code string, a Account) (*Response, *Account, error)
	Close(code string) (*Response, error)
	Reopen(code string) (*Response, error)
	ListNotes(code string) (*Response, []Note, error)
}

// AdjustmentsService represents the interactions available for adjustments.
type AdjustmentsService interface {
	List(accountCode string, params Params) (*Response, []Adjustment, error)
	Get(uuid string) (*Response, *Adjustment, error)
	Create(accountCode string, a Adjustment) (*Response, *Adjustment, error)
	Delete(uuid string) (*Response, error)
}

// AddOnsService represents the interactions available for add ons.
type AddOnsService interface {
	List(planCode string, params Params) (*Response, []AddOn, error)
	Get(planCode string, code string) (*Response, *AddOn, error)
	Create(planCode string, a AddOn) (*Response, *AddOn, error)
	Update(planCode string, code string, a AddOn) (*Response, *AddOn, error)
	Delete(planCode string, code string) (*Response, error)
}

// BillingService represents the interactions available for billing.
type BillingService interface {
	Get(accountCode string) (*Response, *Billing, error)
	Create(accountCode string, b Billing) (*Response, *Billing, error)
	CreateWithToken(accountCode string, token string) (*Response, *Billing, error)
	Update(accountCode string, b Billing) (*Response, *Billing, error)
	UpdateWithToken(accountCode string, token string) (*Response, *Billing, error)
	Clear(accountCode string) (*Response, error)
}

// CouponsService represents the interactions available for coupons.
type CouponsService interface {
	List(params Params) (*Response, []Coupon, error)
	Get(code string) (*Response, *Coupon, error)
	Create(c Coupon) (*Response, *Coupon, error)
	Delete(code string) (*Response, error)
}

// InvoicesService represents the interactions available for invoices.
type InvoicesService interface {
	List(params Params) (*Response, []Invoice, error)
	ListAccount(accountCode string, params Params) (*Response, []Invoice, error)
	Get(invoiceNumber int) (*Response, *Invoice, error)
	GetPDF(invoiceNumber int, language string) (*Response, *bytes.Buffer, error)
	Preview(accountCode string) (*Response, *Invoice, error)
	Create(accountCode string, invoice Invoice) (*Response, *Invoice, error)
	Collect(invoiceNumber int) (*Response, *Invoice, error)
	MarkPaid(invoiceNumber int) (*Response, *Invoice, error)
	MarkFailed(invoiceNumber int) (*Response, *Invoice, error)
	RefundVoidOpenAmount(invoiceNumber int, amountInCents int, refundMethod string) (*Response, *Invoice, error)
	VoidCreditInvoice(invoiceNumber int) (*Response, *Invoice, error)
	RecordPayment(offlinePayment OfflinePayment) (*Response, *Transaction, error)
}

// PlansService represents the interactions available for plans.
type PlansService interface {
	List(params Params) (*Response, []Plan, error)
	Get(code string) (*Response, *Plan, error)
	Create(p Plan) (*Response, *Plan, error)
	Update(code string, p Plan) (*Response, *Plan, error)
	Delete(code string) (*Response, error)
}

// PurchasesService represents the interactions available for a purchase
// involving at least one adjustment or one subscription.
type PurchasesService interface {
	Create(p Purchase) (*Response, *InvoiceCollection, error)
	Preview(p Purchase) (*Response, *InvoiceCollection, error)
}

// RedemptionsService represents the interactions available for redemptions.
type RedemptionsService interface {
	GetForAccount(accountCode string) (*Response, *Redemption, error)
	GetForInvoice(invoiceNumber string) (*Response, *Redemption, error)
	Redeem(code string, accountCode string, currency string) (*Response, *Redemption, error)
	Delete(accountCode string) (*Response, error)
}

// ShippingAddressesService represents the interactions available for shipping addresses.
type ShippingAddressesService interface {
	ListAccount(accountCode string, params Params) (*Response, []ShippingAddress, error)
	Create(accountCode string, address ShippingAddress) (*Response, *ShippingAddress, error)
	Update(accountCode string, shippingAddressID int64, address ShippingAddress) (*Response, *ShippingAddress, error)
	Delete(accountCode string, shippingAddressID int64) (*Response, error)
	GetSubscriptions(accountCode string, shippingAddressID int64) (*Response, []Subscription, error)
}

// SubscriptionsService represents the interactions available for subscriptions.
type SubscriptionsService interface {
	List(params Params) (*Response, []Subscription, error)
	ListAccount(accountCode string, params Params) (*Response, []Subscription, error)
	Get(uuid string) (*Response, *Subscription, error)
	Create(sub NewSubscription) (*Response, *NewSubscriptionResponse, error)
	Preview(sub NewSubscription) (*Response, *Subscription, error)
	Update(uuid string, sub UpdateSubscription) (*Response, *Subscription, error)
	UpdateNotes(uuid string, n SubscriptionNotes) (*Response, *Subscription, error)
	PreviewChange(uuid string, sub UpdateSubscription) (*Response, *Subscription, error)
	Cancel(uuid string) (*Response, *Subscription, error)
	Reactivate(uuid string) (*Response, *Subscription, error)
	TerminateWithPartialRefund(uuid string) (*Response, *Subscription, error)
	TerminateWithFullRefund(uuid string) (*Response, *Subscription, error)
	TerminateWithoutRefund(uuid string) (*Response, *Subscription, error)
	Postpone(uuid string, dt time.Time, bulk bool) (*Response, *Subscription, error)
	Pause(uuid string, cycles int) (*Response, *Subscription, error)
	Resume(uuid string) (*Response, *Subscription, error)
}

// TransactionsService represents the interactions available for transactions.
type TransactionsService interface {
	List(params Params) (*Response, []Transaction, error)
	ListAccount(accountCode string, params Params) (*Response, []Transaction, error)
	Get(uuid string) (*Response, *Transaction, error)
	Create(t Transaction) (*Response, *Transaction, error)
}

// CreditPaymentsService represents the interactions available for credit payments.
type CreditPaymentsService interface {
	List(params Params) (*Response, []CreditPayment, error)
	ListAccount(accountCode string, params Params) (*Response, []CreditPayment, error)
	Get(uuid string) (*Response, *CreditPayment, error)
}
