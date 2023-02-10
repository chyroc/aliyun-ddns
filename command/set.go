package command

import (
	"fmt"
	"github.com/chyroc/aliyun-ddns/aliyun_dns"
	"github.com/urfave/cli/v2"
)

func parseSetParam(c *cli.Context) (string, string, string, string, string, error) {
	domain := c.String("domain")
	rr := c.String("rr")
	ip := c.String("ip")
	akid, aks, err := getAk(c)
	if err != nil {
		return "", "", "", "", "", err
	}
	return domain, rr, ip, akid, aks, nil
}

func Set() *cli.Command {
	return &cli.Command{
		Name: "set",
		Flags: []cli.Flag{
			domainFlag,
			rrFlag,
			ipFlag,
			akidFlag,
			aksFlag,
		},
		Action: func(c *cli.Context) error {
			domain, rr, ip, akid, aks, err := parseSetParam(c)
			if err != nil {
				return err
			}

			dnsCli := aliyun_dns.New(akid, aks)

			return setDomainRR(dnsCli, domain, rr, ip)
		},
	}
}

func setDomainRR(dnsCli *aliyun_dns.Client, domain, rr, ip string) error {
	fmt.Printf("start set '%s.%s' dns to '%s'\n", rr, domain, ip)

	ipRecordType, err := aliyun_dns.DetectIPRecordType(ip)
	if err != nil {
		return err
	}

	targetRrRecord, err := dnsCli.FilterDomainRRRecord(domain, rr)
	if err != nil {
		return fmt.Errorf("list dns record failed: %w", err)
	}
	if targetRrRecord == nil {
		// create
		_, err = dnsCli.Add(domain, rr, ipRecordType, ip)
		if err != nil {
			return fmt.Errorf("add dns record failed: %w", err)
		}
	} else if targetRrRecord.Value == ip {
		// none
		fmt.Printf("dns record '%s.%s' is already '%s', skip update\n", rr, domain, ip)
		return nil
	} else {
		// update
		err = dnsCli.Update(targetRrRecord.RecordId, rr, ipRecordType, ip)
		if err != nil {
			return fmt.Errorf("update dns record failed: %w", err)
		}
	}
	fmt.Printf("update dns record '%s.%s' to '%s' success\n", rr, domain, ip)
	return nil
}
