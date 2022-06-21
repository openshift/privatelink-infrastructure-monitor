package collectors

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

type Route53RecordsPerHostedZoneCollector struct {
	ServiceQuotaClient *servicequotas.Client
	R53Client          *route53.Client
	HostedZoneID       string
}

func (c Route53RecordsPerHostedZoneCollector) Quota() (float64, error) {
	getLimitOut, err := c.R53Client.GetHostedZoneLimit(context.TODO(), &route53.GetHostedZoneLimitInput{
		HostedZoneId: aws.String(c.HostedZoneID),
		Type:         route53types.HostedZoneLimitTypeMaxRrsetsByZone,
	})
	if err != nil {
		return 0, err
	}
	return float64(getLimitOut.Limit.Value), nil
}

func (c Route53RecordsPerHostedZoneCollector) Usage() (float64, error) {
	getLimitOut, err := c.R53Client.GetHostedZoneLimit(context.TODO(), &route53.GetHostedZoneLimitInput{
		HostedZoneId: aws.String(c.HostedZoneID),
		Type:         route53types.HostedZoneLimitTypeMaxRrsetsByZone,
	})
	if err != nil {
		return 0, err
	}
	return float64(getLimitOut.Count), nil
}

func (c Route53RecordsPerHostedZoneCollector) Id() string {
	return "hosted_zone_" + c.HostedZoneID
}

func (c Route53RecordsPerHostedZoneCollector) MetricName() string {
	return "resource_records_per_hosted_zone"
}

func (c Route53RecordsPerHostedZoneCollector) Name() string {
	return "resource_records_per_hosted_zone_" + c.HostedZoneID
}
