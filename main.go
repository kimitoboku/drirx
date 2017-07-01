package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"os"
	"strings"
)

func genRevResolutionDomainName(ip string) string {
	l := strings.Split(ip, ".")
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	revIP := strings.Join(l, ".")
	return revIP + ".in-addr.arpa."
}

func extractA(rr dns.RR) string {
	rrary := strings.SplitN(rr.String(), "\t", 5)
	return rrary[4]
}

func main() {
	log.SetFlags(log.Lshortfile)
	dname := os.Args[1]
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	rdns := config.Servers[0] + ":" + config.Port
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dname), dns.TypeA)
	r, _, err := c.Exchange(m, rdns)
	if err != nil {
		log.Println(err)
	}
	if len(r.Answer) != 0 {
		for i := 0; i < len(r.Answer); i++ {
			if r.Answer[i].Header().Rrtype == dns.TypeA {
				a := extractA(r.Answer[i])
				revIP := genRevResolutionDomainName(a)
				mi := new(dns.Msg)
				mi.SetQuestion(dns.Fqdn(revIP), dns.TypePTR)
				ri, _, err := c.Exchange(mi, rdns)
				if err != nil {
					log.Println(err)
				}

				if len(ri.Answer) != 0 {
					for j := 0; j < len(ri.Answer); j++ {
						ptr := extractA(ri.Answer[j])
						fmt.Printf("%v -> %v -> %v\n", dname, a, ptr)
					}
				}
			}
		}
	}
}
