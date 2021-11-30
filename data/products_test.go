package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "name",
		Price: 1.99,
		SKU: "abcde-abcf-abd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
