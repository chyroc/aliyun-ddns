package aliyun_dns

import (
	"fmt"
	"net"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
)

type Client struct {
	cli *dns.Client
}

func New(akid, aks string) *Client {
	return &Client{cli: dns.NewClientNew(akid, aks)}
}

func (r *Client) FilterDomainRRRecord(domain, rr string) (*dns.RecordType, error) {
	size := 50
	page := 1
	total := 0
	for {
		res, err := r.cli.DescribeDomainRecords(&dns.DescribeDomainRecordsArgs{
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

func (r *Client) Add(domain, rr, recordType, ip string) (string, error) {
	res, err := r.cli.AddDomainRecord(&dns.AddDomainRecordArgs{
		DomainName: domain,
		RR:         rr,
		Type:       recordType,
		Value:      ip,
		TTL:        600,
		Priority:   1,
		Line:       "default",
	})
	if err != nil {
		return "", err
	}
	return res.RecordId, nil
}

func (r *Client) Update(recordID, rr, recordType, ip string) error {
	_, err := r.cli.UpdateDomainRecord(&dns.UpdateDomainRecordArgs{
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

func DetectIPRecordType(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("'%s' is not a valid IP address", s)
	}
	if ip.To4() != nil {
		return dns.ARecord, nil
	}
	return dns.AAAARecord, nil
}
