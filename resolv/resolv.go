package resolv

import (
	"fmt"
	"github.com/miekg/dns"
)

const resolvConf = "/etc/resolv.conf"

func LookupString(qType string, name string) (*dns.Msg, error) {
	t, ok := dns.StringToType[qType]
	if !ok {
		return nil, fmt.Errorf("Invalid type '%s'", qType)
	}

	return Lookup(t, name)
}

func Lookup(qType uint16, name string) (*dns.Msg, error) {
	name = dns.Fqdn(name)
	conf, err := dns.ClientConfigFromFile(resolvConf)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load resolv.conf: %s", err)
	}
	client := &dns.Client{}
	msg := &dns.Msg{}
	msg.SetQuestion(name, qType)

	response := &dns.Msg{}
	for _, server := range conf.Servers {
		server := fmt.Sprintf("%s:%s", server, conf.Port)
		response, err = lookup(msg, client, server, false)
		if err == nil {
			return response, nil
		}
	}
	return response, fmt.Errorf("Couldn't resolve %s: No server responded", name)
}

func lookup(msg *dns.Msg, client *dns.Client, server string, edns bool) (*dns.Msg, error) {
	if edns {
		opt := &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
		}
		opt.SetUDPSize(dns.DefaultMsgSize)
		msg.Extra = append(msg.Extra, opt)
	}

	response, _, err := client.Exchange(msg, server)
	if err != nil {
		return nil, err
	}

	if msg.Id != response.Id {
		return nil, fmt.Errorf("DNS ID mismatch, request: %d, response: %d", msg.Id, response.Id)
	}

	if response.MsgHdr.Truncated {
		if client.Net == "tcp" {
			return nil, fmt.Errorf("Got truncated message on tcp")
		}

		if edns { // Truncated even though EDNS is used
			client.Net = "tcp"
		}

		return lookup(msg, client, server, !edns)
	}

	return response, nil
}
