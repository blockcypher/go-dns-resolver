go-dns-resolver
===============

golang dns resolver library based on miekg/dns

# Usage
The usage is really simple. There are two methods to resolv a name, the only difference: `LookupString` takes a string as query type, `Lookup` takes `dns.Type`:


    func Lookup(qType uint16, name string) (*dns.Msg, error)
    func LookupString(qType string, name string) (*dns.Msg, error)

See the [lookup example](lookup.go)
