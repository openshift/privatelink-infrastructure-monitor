# Privatelink-infrastructure-monitor

This projects monitors [Service
Quotas](https://docs.aws.amazon.com/servicequotas/latest/userguide/intro.html)
in AWS for the privatelink infrastructure.

This is required, so any resource that is close to it's quota can be proactively
be increased, before customers run into problems.

## Testing

### Local

### In-Cluster

Testing the application in an existing OSD cluster is a quick and easy way to
verify it's working, as it will have all the required quotas.

To test the monitoring in a cluster multiple manifests must be applied:

- The privatelink-monitoring application to perform a single gathering run.
- The [pushgateway](https://github.com/prometheus/pushgateway) to receive the
  metrics.

Both can be run in the same cluster, in which case the pushgateway must be
deployed first:

``` sh
oc apply -f deploy/pushgateway/*.yaml
```

This creates a `deployment`, as well as the required `service` and a
`servicemonitor`.

Once the `pushgateway` is up and running, start a monitoring run:

This requires the cronjob to be deployed with the correct *secrets* and
*configmaps*.

The following snippets should set everything up inside a staging cluster (so
make sure that the backplane tunnel is active). Also keep in mind the
credentials generated this way are only valid for about 2 hours, so if there are
authentication errors during the run this might be the case:


```sh
CLUSTER_ID=""
export AWS_DEFAULT_REGION=eu-central-1

if [ -z $CLUSTER_ID ]; then
  echo "Please set a cluster to use for this test"
fi

eval $(ocm backplane cloud credentials $CLUSTER_ID -o env)

oc process -f deploy/10-privatelink.secret.yaml \
    -p AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
    -p AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
    -p AWS_SESSION_TOKEN="$AWS_SESSION_TOKEN" \
    -p AWS_DEFAULT_REGION="$AWS_DEFAULT_REGION" --local | oc apply -f -

VPC=$(aws ec2 describe-vpcs | jq -r '.Vpcs[] | select(.Tags[].Key == "red-hat-clustertype").VpcId')
HOSTED_ZONE=$(aws route53 list-hosted-zones | jq -r '.HostedZones[] | select(.Config.Comment == "Managed by Terraform").Id' | cut -d "/" -f 3)
PUSHGATEWAY=http://prometheus-pushgateway.openshift-monitoring.svc:9091/

oc process -f deploy/11-privatelink.configmap.yaml \
    -p VPC_ID="${VPC}" \
    -p HOSTED_ZONE_ID="${HOSTED_ZONE}" \
    -p DEFAULT_REGION="${AWS_DEFAULT_REGION}" \
    -p PUSHGATEWAY_ADDRESS="${PUSHGATEWAY}" --local | oc apply -f -

oc create -f deploy/20-privatelink.cronjob.yaml
```

The cronjob should now run the application and the metrics can be checked using
a port-forward on the pushgateway:

```sh
oc port-forward prometheus-pushgateway-7c5dcf977d-4f4vv 9091
```

Check http://localhost:9091 for the metrics.
