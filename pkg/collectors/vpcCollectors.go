package collectors

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

const (
	QUOTA_VPCS_PER_REGION                 string = "L-F678F1CE"
	QUOTA_SUBNETS_PER_VPC                 string = "L-407747CB"
	QUOTA_INTERFACE_VPC_ENDPOINTS_PER_VPC string = "L-29B6F2EB"
	QUOTA_ROUTE_TABLES_PER_VPC            string = "L-589F43AA"
	QUOTA_ROUTES_PER_ROUTE_TABLE          string = "L-93826ACB"
	QUOTA_CODE_IPV4_BLOCKS_PER_VPC        string = "L-83CA0A9D"
	SERVICE_CODE_VPC                      string = "vpc"
)

type VpcsPerRegion struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	Region             string
}

func (c VpcsPerRegion) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_VPCS_PER_REGION)
}

func (c VpcsPerRegion) Usage() (float64, error) {
	describeVpcsOutput, err := c.Ec2Client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return 0, err
	}
	return float64(len(describeVpcsOutput.Vpcs)), nil
}

func (c VpcsPerRegion) Id() string {
	return "region_" + c.Region
}

func (c VpcsPerRegion) MetricName() string {
	return "vpcs_per_region"
}

func (c VpcsPerRegion) Name() string {
	return "vpcs_per_region_" + c.Region
}

type SubnetsPerVpc struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	VpcID              string
}

func (c SubnetsPerVpc) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_SUBNETS_PER_VPC)
}

func (c SubnetsPerVpc) Usage() (float64, error) {
	describeSubnetsOutput, err := c.Ec2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{{
			Name:   aws.String("vpc-id"),
			Values: []string{c.VpcID},
		}},
	})
	if err != nil {
		return 0, err
	}
	return float64(len(describeSubnetsOutput.Subnets)), nil
}

func (c SubnetsPerVpc) Id() string {
	return "vpc_" + c.VpcID
}

func (c SubnetsPerVpc) MetricName() string {
	return "subnets_per_vpc"
}

func (c SubnetsPerVpc) Name() string {
	return "subnets_per_vpc_" + c.VpcID
}

type InterfaceVpcEndpointsPerVpc struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	VpcID              string
}

func (c InterfaceVpcEndpointsPerVpc) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_INTERFACE_VPC_ENDPOINTS_PER_VPC)
}

func (c InterfaceVpcEndpointsPerVpc) Usage() (float64, error) {
	describeVpcEndpointsOutput, err := c.Ec2Client.DescribeVpcEndpoints(context.TODO(), &ec2.DescribeVpcEndpointsInput{
		Filters: []types.Filter{{
			Name:   aws.String("vpc-id"),
			Values: []string{c.VpcID},
		}},
	})
	if err != nil {
		return 0, err
	}
	numEndpoints := len(describeVpcEndpointsOutput.VpcEndpoints)
	return float64(numEndpoints), nil
}

func (c InterfaceVpcEndpointsPerVpc) Id() string {
	return "vpc_" + c.VpcID
}

func (c InterfaceVpcEndpointsPerVpc) MetricName() string {
	return "interface_vpc_endpoints_per_vpc"
}

func (c InterfaceVpcEndpointsPerVpc) Name() string {
	return "interface_vpc_endpoints_per_vpc_" + c.VpcID
}

type RoutesPerRouteTableCollector struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	RouteTableID       string
}

func (c RoutesPerRouteTableCollector) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_ROUTES_PER_ROUTE_TABLE)
}

func (c RoutesPerRouteTableCollector) Usage() (float64, error) {
	descRouteTableOutput, err := c.Ec2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
		RouteTableIds: []string{c.RouteTableID},
	})
	if err != nil {
		return 0, err
	}
	numRoutes := len(descRouteTableOutput.RouteTables[0].Routes)
	return float64(numRoutes), nil
}

func (c RoutesPerRouteTableCollector) Id() string {
	return "route_table_" + c.RouteTableID
}

func (c RoutesPerRouteTableCollector) MetricName() string {
	return "routes_per_route_table_collector"
}

func (c RoutesPerRouteTableCollector) Name() string {
	return "routes_per_route_table_" + c.RouteTableID
}

type RouteTablesPerVPCCollector struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	VpcID              string
}

func (c RouteTablesPerVPCCollector) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_ROUTE_TABLES_PER_VPC)
}

func (c RouteTablesPerVPCCollector) Usage() (float64, error) {
	descRouteTableOutput, err := c.Ec2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{{
			Name:   aws.String("vpc-id"),
			Values: []string{c.VpcID},
		}},
	})

	if err != nil {
		return 0, err
	}
	numRTB := len(descRouteTableOutput.RouteTables)
	return float64(numRTB), nil
}

func (c RouteTablesPerVPCCollector) Id() string {
	return "vpc_" + c.VpcID
}

func (c RouteTablesPerVPCCollector) MetricName() string {
	return "route_tables_per_vpc"
}

func (c RouteTablesPerVPCCollector) Name() string {
	return "route_tables_per_vpc_" + c.VpcID
}

type Ipv4BlocksPerVPCCollector struct {
	ServiceQuotaClient *servicequotas.Client
	Ec2Client          *ec2.Client
	VpcID              string
}

func (c Ipv4BlocksPerVPCCollector) Quota() (float64, error) {
	return GetQuotaValue(c.ServiceQuotaClient, SERVICE_CODE_VPC, QUOTA_CODE_IPV4_BLOCKS_PER_VPC)
}

func (c Ipv4BlocksPerVPCCollector) Usage() (float64, error) {
	descVpcOut, err := c.Ec2Client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
		DryRun: aws.Bool(false),
		VpcIds: []string{c.VpcID},
	})

	if err != nil {
		return 0, err
	}
	if len(descVpcOut.Vpcs) != 1 {
		return 0, errors.New("Unexcpected number of VPCs returned")
	}

	return float64(len(descVpcOut.Vpcs[0].CidrBlockAssociationSet)), nil
}

func (c Ipv4BlocksPerVPCCollector) Id() string {
	return "vpc_" + c.VpcID
}

func (c Ipv4BlocksPerVPCCollector) MetricName() string {
	return "ipv4_blocks_per_vpc"
}

func (c Ipv4BlocksPerVPCCollector) Name() string {
	return "ipv4_blocks_per_vpc_" + c.VpcID
}
