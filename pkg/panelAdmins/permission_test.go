package panelAdmins

import (
	"testing"
)

type TestPermission struct {
	oldPermission Permission
	newPermission Permission
	expected      bool
}

func TestPermission_ConvertToJson(t *testing.T) {}

func TestUpdatePermission(t *testing.T) {
	table := []TestPermission{
		{
			expected: true,
			newPermission: Permission{
				Onboarding: ReadWrite{true, true},
				Role:       ReadWrite{true, true},
				Team:       ReadWrite{true, true},
				Tenant:     ReadWrite{true, true},
				Billing:    ReadWrite{true, true},
			},
			oldPermission: Permission{
				Onboarding: ReadWrite{false, true},
				Role:       ReadWrite{true, false},
				Team:       ReadWrite{false, true},
				Tenant:     ReadWrite{true, true},
				Billing:    ReadWrite{false, false},
			},
		},
	}

	for _, tt := range table {
		p := tt.oldPermission

		if err := p.UpdatePermission(tt.newPermission); err != nil {
			t.Fatal("could not update the permission")
		}

		if tt.newPermission == p {
			t.Logf("Expected to get %v, and got %v", tt.expected, tt.newPermission == p)
		}
	}
}
