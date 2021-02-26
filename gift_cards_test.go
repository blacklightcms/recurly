package recurly_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/fubotv/go-recurly"
	"github.com/google/go-cmp/cmp"
)

func TestGiftCards_Encoding(t *testing.T) {
	t.Parallel()

	moment, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	tests := []struct {
		v        recurly.GiftCard
		expected string
	}{
		{
			v: recurly.GiftCard{
				XMLName: xml.Name{Local: "gift_card"},
				ID:      2003020297591186183,
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
				</gift_card>
			`),
		},
		{
			v: recurly.GiftCard{
				XMLName:  xml.Name{Local: "gift_card"},
				ID:       2003020297591186183,
				Delivery: &recurly.Delivery{},
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
					<delivery></delivery>
				</gift_card>
			`),
		},
		{
			v: recurly.GiftCard{
				XMLName:           xml.Name{Local: "gift_card"},
				ID:                2003020297591186183,
				RedemptionCode:    "518822D87268C142",
				BalanceInCents:    2999,
				ProductCode:       "gift_card",
				UnitAmountInCents: 2999,
				Currency:          "USD",
				CreatedAt:         recurly.NewTime(moment),
				UpdatedAt:         recurly.NewTime(moment),
				DeliveredAt:       recurly.NewTime(moment),
				RedeemedAt:        recurly.NewTime(moment),
				CanceledAt:        recurly.NewTime(moment),
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
					<redemption_code>518822D87268C142</redemption_code>
					<balance_in_cents>2999</balance_in_cents>
					<product_code>gift_card</product_code>
					<unit_amount_in_cents>2999</unit_amount_in_cents>
					<currency>USD</currency>
					<created_at>2014-01-01T07:00:00Z</created_at>
					<updated_at>2014-01-01T07:00:00Z</updated_at>
					<delivered_at>2014-01-01T07:00:00Z</delivered_at>
					<redeemed_at>2014-01-01T07:00:00Z</redeemed_at>
					<canceled_at>2014-01-01T07:00:00Z</canceled_at>
				</gift_card>
			`),
		},
	}

	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("Encode/%d", i), func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if diff := cmp.Diff(buf.String(), tt.expected); diff != "" {
				t.Fatal(diff)
			}
		})

		t.Run(fmt.Sprintf("Decode/%d", i), func(t *testing.T) {
			t.Parallel()

			var g recurly.GiftCard
			if err := xml.Unmarshal([]byte(tt.expected), &g); err != nil {
				t.Fatal(err)
			} else if diff := cmp.Diff(tt.v, g); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
