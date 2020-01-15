package gmicro_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

func TestMessageHandler(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "Test"}
	gm.CreateGroup(&g)

	m1 := member.Member{ID: uuid.New(), Name: "Test1"}
	gm.AddMember(g.ID, &m1)

	m2 := member.Member{ID: uuid.New(), Name: "Test2"}
	gm.AddMember(g.ID, &m2)

	e := expense.Expense{
		ID:         uuid.New(),
		GroupID:    g.ID,
		Amount:     25.4,
		Payer:      m1.ID,
		Recipients: m1.ID.String() + ";" + m2.ID.String(),
	}
	ebody1, _ := json.Marshal(&e)

	e.GroupID = uuid.New()
	ebody2, _ := json.Marshal(&e)

	p := payment.Payment{
		ID:        uuid.New(),
		GroupID:   g.ID,
		Amount:    14.6,
		Payer:     m1.ID,
		Recipient: m2.ID,
	}
	pbody1, _ := json.Marshal(&p)

	p.GroupID = uuid.New()
	pbody2, _ := json.Marshal(&p)

	cases := []struct {
		op   string
		body []byte
		fail bool
	}{
		{"add-expense", ebody1, false},
		{"add-expense", []byte(`{"id":"test"}`), true},
		{"add-expense", ebody2, true},

		{"delete-expense", ebody1, false},
		{"delete-expense", []byte(`{"id":"test"}`), true},
		{"delete-expense", ebody2, true},

		{"add-payment", pbody1, false},
		{"add-payment", []byte(`{"id":"test"}`), true},
		{"add-payment", pbody2, true},

		{"delete-payment", pbody1, false},
		{"delete-payment", []byte(`{"id":"test"}`), true},
		{"delete-payment", pbody2, true},
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
