package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	domainFlag = &cli.StringFlag{
		Name:     "domain",
		Aliases:  []string{"d"},
		Usage:    "domain name",
		Required: true,
	}
	rrFlag = &cli.StringFlag{
		Name:     "rr",
		Aliases:  []string{"r"},
		Usage:    "resource record",
		Required: true,
	}
	ipFlag = &cli.StringFlag{
		Name:     "ip",
		Aliases:  []string{"i"},
		Usage:    "ip address",
		Required: true,
	}
	ipTypeFlag = &cli.StringFlag{
		Name:     "ip-type",
		Aliases:  []string{"t"},
		Usage:    "ip type(ipv4, ipv6)",
		Required: true,
	}
	akidFlag = &cli.StringFlag{
		Name:     "access-key-id",
		Aliases:  []string{"akid"},
		Usage:    "aliyun access key id(default: env ALIYUN_ACCESS_KEY_ID)",
		Required: false,
	}
	aksFlag = &cli.StringFlag{
		Name:     "access-key-secret",
		Aliases:  []string{"aks"},
		Usage:    "aliyun access key secret(default: env ALIYUN_ACCESS_KEY_SECRET)",
		Required: false,
	}
)

func getAk(c *cli.Context) (string, string, error) {
	akid := c.String("access-key-id")
	aks := c.String("access-key-secret")
	if akid == "" {
		akid = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if aks == "" {
		aks = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}
	if akid == "" || aks == "" {
		return "", "", fmt.Errorf("access-key-id or access-key-secret is empty")
	}

	return akid, aks, nil
}
