package main

import "testing"

func TestIsNull(t *testing.T){
	test := ""
	expected := true
	actual := IsNull(&test)
	if actual != expected{
		t.Error("Failed with null String")
	}
	test = "Some Apples"
	expected = false
	actual = IsNull(&test)
	if actual != expected{
		t.Error("Failed with actual String")
	}
}

func TestCheckAction(t *testing.T) {
	tests := map[string]bool{
		"GET": true,
		"POST": true,
		"DELETE": true,
		"PUT": true,
		"PINEAPPLE": false,
	}
	for verb, expected := range tests{
		actual := CheckAction(&verb)
		if actual != expected {
			t.Error("Failed with ", verb)
		}

	}

}