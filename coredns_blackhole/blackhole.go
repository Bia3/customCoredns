package coredns_blackhole

import (
	"context"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

const BlockedRcode = dns.RcodeRefused

var log = clog.NewWithPlugin("blackhole")

type Blackhole struct {
	Next      plugin.Handler
	Blocklist *Blocklist
}

func (b Blackhole) Name() string { return "blackhole" }

func (b Blackhole) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	if b.Blocklist.Find(state.Name()) {
		log.Debugf("Blocked domain %s", state.Name())
		a := new(dns.Msg)
		//a.SetRcode(r, BlockedRcode)
		av4Rec, err := dns.NewRR(state.Name() + " 3600 IN A 127.0.0.1")
		if err != nil {
			log.Error(err)
		}
		av6Rec, err := dns.NewRR(state.Name() + " 3600 IN AAAA ::1")
		if err != nil {
			log.Error(err)
		}
		a.Answer = []dns.RR{
			av4Rec,
			av6Rec,
		}
		a.Authoritative = true
		//a.RecursionAvailable = true
		//a.RecursionDesired = r.RecursionDesired
		a.Id = state.Req.Id
		w.WriteMsg(a)
		return dns.RcodeSuccess, nil
	}
	return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
}
