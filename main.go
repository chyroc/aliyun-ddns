package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/urfave/cli/v2"
)

// https://help.aliyun.com/document_detail/124923.html

func parseParam(c *cli.Context) (string, string, string, string, string, error) {
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

func filterAliyunDnsRecord(dnsCli *dns.Client, domain, rr string) (*dns.RecordType, error) {
	size := 50
	page := 1
	total := 0
	for {
		res, err := dnsCli.DescribeDomainRecords(&dns.DescribeDomainRecordsArgs{
			DomainName: domain,
			Pagination: common.Pagination{PageNumber: page, PageSize: size},
			RRKeyWord:  rr,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range res.DomainRecords.Record {
			if v.RR == rr {
				return &v, nil
			}
		}
		total += len(res.DomainRecords.Record)
		if len(res.DomainRecords.Record) < size || total >= res.TotalCount {
			break
		}
		page++
	}
	return nil, nil
}

func addAliyunDnsRecord(dnsCli *dns.Client, domain, rr, recordType, ip string) error {
	_, err := dnsCli.AddDomainRecord(&dns.AddDomainRecordArgs{
		DomainName: domain,
		RR:         rr,
		Type:       recordType,
		Value:      ip,
		TTL:        600,
		Priority:   1,
		Line:       "default",
	})
	return err
}

func updateAliyunDnsRecord(dnsCli *dns.Client, recordID, rr, recordType, ip string) error {
	_, err := dnsCli.UpdateDomainRecord(&dns.UpdateDomainRecordArgs{
		RecordId: recordID,
		RR:       rr,
		Type:     recordType,
		Value:    ip,
		TTL:      600,
		Priority: 1,
		Line:     "default",
	})
	return err
}

func detectIPRecordType(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("'%s' is not a valid IP address", s)
	}
	if ip.To4() != nil {
		return dns.ARecord, nil
	}
	return dns.AAAARecord, nil
}

func main() {
	app := &cli.App{
		Name: "aliyun-ddns",
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
			domain, rr, ip, akid, aks, err := parseParam(c)
			if err != nil {
				return err
			}

			dnsCli := dns.NewClientNew(akid, aks)

			ipRecordType, err := detectIPRecordType(ip)
			if err != nil {
				return err
			}

			targetRrRecord, err := filterAliyunDnsRecord(dnsCli, domain, rr)
			if err != nil {
				return fmt.Errorf("list dns record failed: %w", err)
			}
			if targetRrRecord == nil {
				// create
				err = addAliyunDnsRecord(dnsCli, domain, rr, ipRecordType, ip)
				if err != nil {
					return fmt.Errorf("add dns record failed: %w", err)
				}
			} else if targetRrRecord.Value == ip {
				// none
				fmt.Printf("dns record '%s.%s' is already '%s', skip update\n", rr, domain, ip)
				return nil
			} else {
				// update
				err = updateAliyunDnsRecord(dnsCli, targetRrRecord.RecordId, rr, ipRecordType, ip)
				if err != nil {
					return fmt.Errorf("update dns record failed: %w", err)
				}
			}
			fmt.Printf("update dns record '%s.%s' to '%s' success\n", rr, domain, ip)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
