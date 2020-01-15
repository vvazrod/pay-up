package tmicro_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

func TestMessageHandler(t *testing.T) {
	e := expense.Expense{
		ID:         uuid.New(),
		GroupID:    uuid.New(),
		Amount:     25.4,
		Payer:      uuid.New(),
		Recipients: uuid.New().String() + ";" + uuid.New().String(),
	}
	ebody, _ := json.Marshal(&e)

	p := payment.Payment{
		ID:        uuid.New(),
		GroupID:   uuid.New(),
		Amount:    14.6,
		Payer:     uuid.New(),
		Recipient: uuid.New(),
	}
	pbody, _ := json.Marshal(&p)

	cases := []struct {
		op   string
		body []byte
		fail bool
	}{
		{"add-expense", ebody, false},
		{"add-expense", []byte(`{"id":"test"}`), true},
		{"add-expense", []byte(`{"recipient":"test;"}`), true},

		{"delete-expense", []byte(`{"group_id":"` + e.GroupID.String() + `"}`), false},
		{"delete-expense", []byte(``), true},
		{"delete-expense", []byte(`{}`), true},
		{"delete-expense", []byte(`{"group_id":"test"}`), true},
		{"delete-expense", []byte(`{"group_id":"` + uuid.New().String() + `"}`), true},

		{"add-payment", pbody, false},
		{"add-payment", []byte(`{"id":"test"}`), true},

		{"delete-payment", []byte(`{"group_id":"` + p.GroupID.String() + `"}`), false},
		{"delete-payment", []byte(``), true},
		{"delete-payment", []byte(`{}`), true},
		{"delete-payment", []byte(`{"group_id":"test"}`), true},
		{"delete-payment", []byte(`{"group_id":"` + uuid.New().String() + `"}`), true},
	}

	for _, tc := range cases {
		if !tc.fail {
			t.Run(fmt.Sprintf("Correct %s", tc.op), func(t *testing.T) {
				err := h(tc.op, tc.body)

				if err != nil {
					t.Errorf("Operation can't finish correctly [Error]: %v", err)
				}
			})
		} else {
			t.Run(fmt.Sprintf("Error %s", tc.op), func(t *testing.T) {
				err := h(tc.op, tc.body)

				if err == nil {
					t.Error("Using wrong values didn't return an error")
				}
			})
		}
	}
}
