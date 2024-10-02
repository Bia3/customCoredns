package coredns_blackhole

import (
	"context"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"strings"
	"testing"
)

func contains(s []dns.RR, e dns.RR) bool {
	for _, a := range s {
		if strings.Compare(a.String(), e.String()) == 0 {
			return true
		}
	}
	return false
}

func TestBlackhole_ServeDNS(t *testing.T) {
	// Create a new Example Plugin. Use the test.ErrorHandler as the next plugin.
	list := NewBlocklist()
	exampleNetARR, _ := dns.NewRR("example.net.\t3600\tIN\tA\t127.0.0.1")
	exampleNetAAARR, _ := dns.NewRR("example.net.\t3600\tIN\tAAAA\t::1")
	list.Add("example.net.")
	bl := Blackhole{Next: test.ErrorHandler(), Blocklist: list}

	ctx := context.TODO()

	tests := []struct {
		reqDomain      string
		shouldErr      bool
		aName          dns.RR
		aaaaName       dns.RR
		expectedRCode  int
		expectedNumRes int
	}{
		{
			"example.net",
			false,
			exampleNetARR,
			exampleNetAAARR,
			dns.RcodeSuccess,
			2,
		},

		{
			"example.com",
			false,
			nil,
			nil,
			dns.RcodeServerFailure,
			0,
		},
	}

	for i, testcase := range tests {
		req := new(dns.Msg)

		req.SetQuestion(dns.Fqdn(testcase.reqDomain), dns.TypeA)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})

		code, err := bl.ServeDNS(ctx, rec, req)

		if err == nil && testcase.shouldErr {
			t.Fatalf("Test %d: expected errors, but got no error", i)
		} else if err != nil && !testcase.shouldErr {
			t.Fatalf("Test %d: expected no errors, but got '%v'", i, err)
			continue
		}

		if code != testcase.expectedRCode {
			t.Errorf("Test %d: Expected status code %d, but got %d", i, testcase.expectedRCode, code)
		}

		if rec.Rcode != testcase.expectedRCode {
			t.Errorf("Test %d: Expected return code %d, but got %d", i, testcase.expectedRCode, rec.Rcode)

		}
		if len(rec.Msg.Answer) != testcase.expectedNumRes {
			t.Errorf("Test %d: Expected %d results, but got %d", i, testcase.expectedNumRes, len(rec.Msg.Answer))
		}

		//_, okAName := rec.Msg.Answer[testcase.aName]
		if len(rec.Msg.Answer) >= 1 {
			if !contains(rec.Msg.Answer, testcase.aName) {
				t.Errorf("Test %d: Expected %a in Answers but got %v", i, testcase.aName, rec.Msg.Answer)
			}
			if !contains(rec.Msg.Answer, testcase.aaaaName) {
				t.Errorf("Test %d: Expected %a in Answers but got %v", i, testcase.aaaaName, rec.Msg.Answer)
			}
		}
	}
}
