package command

import (
	"fmt"
	"os"

	"github.com/chyroc/aliyun-ddns/aliyun_dns"
	"github.com/urfave/cli/v2"
)

func parseGetParam(c *cli.Context) (string, string, string, string, error) {
	domain := c.String("domain")
	rr := c.String("rr")
	akid := c.String("access-key-id")
	aks := c.String("access-key-secret")
	if akid == "" {
		akid = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if aks == "" {
		aks = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}
	if akid == "" || aks == "" {
		return "", "", "", "", fmt.Errorf("access-key-id or access-key-secret is empty")
	}
	return domain, rr, akid, aks, nil
}

func Get() *cli.Command {
	return &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "domain name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "rr",
				Aliases:  []string{"r"},
				Usage:    "resource record",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "access-key-id",
				Aliases:  []string{"akid"},
				Usage:    "aliyun access key id(default: env ALIYUN_ACCESS_KEY_ID)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "access-key-secret",
				Aliases:  []string{"aks"},
				Usage:    "aliyun access key secret(default: env ALIYUN_ACCESS_KEY_SECRET)",
				Required: false,
			},
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
