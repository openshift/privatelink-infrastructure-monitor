package collectors

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

const (
	QUOTA_CODE_TRANSIT_GATEWAY_PER_ACCT string = "L-A2478D36"
	SERVICE_CODE_EC2                    string = "ec2"
)

type TransitGatewaysPerAcctCollector struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
}

func (tg TransitGatewaysPerAcctCollector) Quota() (float64, error) {
	return GetQuotaValue(tg.ServiceQuotaClient, SERVICE_CODE_EC2, QUOTA_CODE_TRANSIT_GATEWAY_PER_ACCT)
}

func (tg TransitGatewaysPerAcctCollector) Usage() (float64, error) {
	transitGatewayOut, err := tg.Ec2Client.DescribeTransitGateways(context.TODO(), &ec2.DescribeTransitGatewaysInput{
		DryRun:     aws.Bool(false),
		MaxResults: aws.Int32(100),
		NextToken:  nil,
	})
	if err != nil {
		return 0, err
	}
	return float64(len(transitGatewayOut.TransitGateways)), nil
}

func (tg TransitGatewaysPerAcctCollector) Id() string {
	return "all"
}

func (tg TransitGatewaysPerAcctCollector) MetricName() string {
	return "transit_gateways_per_account"
}

func (tg TransitGatewaysPerAcctCollector) Name() string {
	return "transit_gateways_per_account"
}
