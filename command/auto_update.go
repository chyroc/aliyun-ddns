package command

import (
	"fmt"
	"github.com/chyroc/aliyun-ddns/aliyun_dns"
	"github.com/chyroc/detect-ip/detect_ip"
	"github.com/urfave/cli/v2"
	"net"
	"time"
)

func parseAutoUpdateParam(c *cli.Context) (string, string, string, string, string, error) {
	domain := c.String("domain")
	rr := c.String("rr")
	ipType := c.String("ip-type")
	if ipType != "ipv4" && ipType != "ipv6" {
		return "", "", "", "", "", fmt.Errorf("ip-type must ipv4 or ipv6")
	}
	akid, aks, err := getAk(c)
	if err != nil {
		return "", "", "", "", "", err
	}
	return domain, rr, ipType, akid, aks, nil
}

func UpdateSet() *cli.Command {
	return &cli.Command{
		Name: "auto-update",
		Flags: []cli.Flag{
			domainFlag,
			rrFlag,
			ipTypeFlag,
			akidFlag,
			aksFlag,
		},
		Action: func(c *cli.Context) error {
			domain, rr, ipType, akid, aks, err := parseAutoUpdateParam(c)
			if err != nil {
				return err
			}

			dnsCli := aliyun_dns.New(akid, aks)

			// for loop 5 mis
			_ = autoUpdateIP(dnsCli, domain, rr, ipType)
			ticker := time.NewTicker(time.Minute * 5)
			for range ticker.C {
				_ = autoUpdateIP(dnsCli, domain, rr, ipType)
			}

			return nil
		},
	}
}

func detectIP(ipType string) net.IP {
	switch ipType {
	case "ipv4":
		return detect_ip.PublicIPV4(detect_ip.WithTimeout(time.Second * 8))
	case "ipv6":
		return detect_ip.PublicIPV6(detect_ip.WithTimeout(time.Second * 8))
	}
	return nil
}

func autoUpdateIP(dnsCli *aliyun_dns.Client, domain, rr, ipType string) net.IP {
	fmt.Printf("run at %s\n", time.Now())
	ip := detectIP(ipType)
	if ip == nil {
		// TODO:
		fmt.Printf("detect ip failed, at %s\n", time.Now())
		return nil
	}

	err := setDomainRR(dnsCli, domain, rr, ip.String())
	if err != nil {
		fmt.Printf("set ip failed, err %s, at %s\n", err, time.Now())
	}
	return ip
}
