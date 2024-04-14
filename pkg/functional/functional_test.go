package functional_test

import (
	"fmt"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/functional"
)

type addr struct {
	address string
	port    int16
}

func Test_functional_int(t *testing.T) {
	tc := []struct {
		name     string
		list     []int
		expected string
	}{
		{"add all integers to string", []int{1, 2}, "012"},
	}

	for _, d := range tc {
		t.Run(d.name, func(t *testing.T) {
			result := functional.FoldMap(
				d.list,
				0,
				func(a int) string {
					return fmt.Sprint(a)
				},
				func(curr, next string) string {
					return curr + next
				})
			if result != d.expected {
				t.Errorf("Expected %s, got %s", d.expected, result)
			}
		})
	}
}

func Test_functional_address(t *testing.T) {
	tc := []struct {
		name     string
		list     []addr
		expected string
	}{
		{"add all addresses to string", []addr{{address: "http://foo", port: 1234}}, "http://foo:1234"},
	}

	for _, d := range tc {
		t.Run(d.name, func(t *testing.T) {
			result := functional.FoldMap(
				d.list,
				addr{},
				func(a addr) string {
					if a.address == "" && a.port == 0 {
						return ""
					}

					return fmt.Sprintf("%s:%d", a.address, a.port)
				},
				func(curr, next string) string {
					return curr + next
				})
			if result != d.expected {
				t.Errorf("Expected %s, got %s", d.expected, result)
			}
		})
	}
}
