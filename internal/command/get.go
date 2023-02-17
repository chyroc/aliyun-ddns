package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/aliyun-ddns/aliyun_dns"
)

func parseGetParam(c *cli.Context) (string, string, string, string, error) {
	domain := c.String("domain")
	rr := c.String("rr")
	akid := c.String("access-key-id")
	aks := c.String("access-key-secret")
	akid, aks, err := getAk(c)
	if err != nil {
		return "", "", "", "", err
	}
	return domain, rr, akid, aks, nil
}

func Get() *cli.Command {
	return &cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			domainFlag,
			rrFlag,
			akidFlag,
			aksFlag,
		},
		Action: func(c *cli.Context) error {
			domain, rr, akid, aks, err := parseGetParam(c)
			if err != nil {
				return err
			}

			dnsCli := aliyun_dns.New(akid, aks)

			targetRrRecord, err := dnsCli.FilterDomainRRRecord(domain, rr)
			if err != nil {
				return fmt.Errorf("list dns record failed: %w", err)
			}
			if targetRrRecord == nil {
				fmt.Printf("none dns record '%s.%s'\n", rr, domain)
			} else {
				fmt.Printf("none dns record '%s.%s' is '%s'\n", rr, domain, targetRrRecord.Value)
			}
			return nil
		},
	}
}
