package collectors

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

type QuotaCollector interface {
	Quota() (float64, error)
	Usage() (float64, error)
	Name() string
	Id() string
	MetricName() string
}

// GetQuotaValue returns the value of an AWS  Quota as identified by the given service and quota code
func GetQuotaValue(client *servicequotas.Client, serviceCode string, quotaCode string) (float64, error) {
	sqOutput, err := client.GetServiceQuota(context.TODO(), &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	})

	if err != nil {
		return 0, err
	}

	return *sqOutput.Quota.Value, nil
}
