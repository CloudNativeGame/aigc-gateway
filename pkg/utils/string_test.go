package utils

import "testing"

func TestGetDomainFromEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		domain   string
	}{
		{
			endpoint: "https://logto.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com/",
			domain:   ".c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com",
		},
		{
			endpoint: "https://logto.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com",
			domain:   ".c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com",
		},
	}

	for _, test := range tests {
		actual := GetDomainFromEndpoint(test.endpoint)
		expect := test.domain
		if actual != expect {
			t.Errorf("expect domain: %s, but actual got %v", expect, actual)
		}
	}
}
