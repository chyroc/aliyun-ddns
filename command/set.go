package command

import (
	"fmt"
	"os"

	"github.com/chyroc/aliyun-ddns/aliyun_dns"
	"github.com/urfave/cli/v2"
)

func parseSetParam(c *cli.Context) (string, string, string, string, string, error) {
	domain := c.String("domain")
	rr := c.String("rr")
	ip := c.String("ip")
	akid := c.String("access-key-id")
	aks := c.String("access-key-secret")
	if akid == "" {
		akid = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if aks == "" {
		aks = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}
	if akid == "" || aks == "" {
		return "", "", "", "", "", fmt.Errorf("access-key-id or access-key-secret is empty")
	}
	return domain, rr, ip, akid, aks, nil
}

func Set() *cli.Command {
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
				Name:     "ip",
				Aliases:  []string{"i"},
				Usage:    "ip address",
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
			domain := c.String("domain")
			rr := c.String("rr")
			ip := c.String("ip")
			domain, rr, ip, akid, aks, err := parseSetParam(c)
			if err != nil {
				return err
			}

			dnsCli := aliyun_dns.New(akid, aks)

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
		},
	}
}
