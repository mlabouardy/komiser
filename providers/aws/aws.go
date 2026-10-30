package aws

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	log "github.com/sirupsen/logrus"
	"github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"

	"github.com/tailwarden/komiser/providers/aws/apigateway"
	"github.com/tailwarden/komiser/providers/aws/cloudfront"
	"github.com/tailwarden/komiser/providers/aws/cloudtrail"
	"github.com/tailwarden/komiser/providers/aws/cloudwatch"
	"github.com/tailwarden/komiser/providers/aws/codebuild"
	"github.com/tailwarden/komiser/providers/aws/codecommit"
	"github.com/tailwarden/komiser/providers/aws/codedeploy"
	"github.com/tailwarden/komiser/providers/aws/datasync"
	"github.com/tailwarden/komiser/providers/aws/dynamodb"
	"github.com/tailwarden/komiser/providers/aws/ec2"
	"github.com/tailwarden/komiser/providers/aws/ecr"
	"github.com/tailwarden/komiser/providers/aws/ecs"
	"github.com/tailwarden/komiser/providers/aws/efs"
	"github.com/tailwarden/komiser/providers/aws/eks"
	"github.com/tailwarden/komiser/providers/aws/elasticache"
	"github.com/tailwarden/komiser/providers/aws/elb"
	"github.com/tailwarden/komiser/providers/aws/firehose"
	"github.com/tailwarden/komiser/providers/aws/iam"
	"github.com/tailwarden/komiser/providers/aws/kafka"
	"github.com/tailwarden/komiser/providers/aws/kinesis"
	"github.com/tailwarden/komiser/providers/aws/kinesisanalytics"
	"github.com/tailwarden/komiser/providers/aws/kms"
	"github.com/tailwarden/komiser/providers/aws/lambda"
	"github.com/tailwarden/komiser/providers/aws/lightsail"
	"github.com/tailwarden/komiser/providers/aws/neptune"
	"github.com/tailwarden/komiser/providers/aws/opensearch"
	"github.com/tailwarden/komiser/providers/aws/rds"
	"github.com/tailwarden/komiser/providers/aws/redshift"
	"github.com/tailwarden/komiser/providers/aws/route53"
	"github.com/tailwarden/komiser/providers/aws/s3"
	"github.com/tailwarden/komiser/providers/aws/secretsmanager"
	"github.com/tailwarden/komiser/providers/aws/servicecatalog"
	"github.com/tailwarden/komiser/providers/aws/sns"
	"github.com/tailwarden/komiser/providers/aws/sqs"
	"github.com/tailwarden/komiser/providers/aws/systemsmanager"
	awsUtils "github.com/tailwarden/komiser/providers/aws/utils"
	"github.com/tailwarden/komiser/providers/aws/wafv2"
	"github.com/tailwarden/komiser/utils"

	"github.com/uptrace/bun"
)

func listOfSupportedServices() []providers.FetchDataFunction {
	return []providers.FetchDataFunction{
		ec2.Instances,
		ec2.ElasticIps,
		lambda.Functions,
		lambda.EventSourceMappings,
		ec2.Acls,
		ec2.Subnets,
		ec2.SecurityGroups,
		ec2.AutoScalingGroups,
		ec2.InternetGateways,
		iam.Roles,
		iam.InstanceProfiles,
		iam.OIDCProviders,
		iam.SamlProviders,
		iam.Groups,
		iam.Policies,
		iam.Users,
		sqs.Queues,
		s3.Buckets,
		ec2.Instances,
		eks.KubernetesClusters,
		cloudfront.Distributions,
		cloudfront.Functions,
		dynamodb.Tables,
		ecs.Clusters,
		ecs.TaskDefinitions,
		ecs.ContainerInstances,
		ecr.Repositories,
		sns.Topics,
		ec2.Vpcs,
		ec2.Volumes,
		kms.Keys,
		rds.Clusters,
		rds.Instances,
		rds.Proxies,
		rds.Snapshots,
		rds.ClusterSnapshots,
		rds.ProxyEndpoints,
		rds.AutoBackups,
		elb.LoadBalancers,
		elb.TargetGroups,
		efs.ElasticFileStorage,
		apigateway.Apis,
		elasticache.Clusters,
		cloudwatch.Alarms,
		ec2.NetworkInterfaces,
		cloudwatch.Dashboards,
		ec2.ElasticIps,
		cloudwatch.LogGroups,
		cloudwatch.MetricStreams,
		ec2.Snapshots,
		opensearch.ServiceDomains,
		servicecatalog.Products,
		ec2.SpotInstanceRequests,
		ec2.KeyPairs,
		ec2.PlacementGroups,
		ec2.NetworkInsightsAnalyses,
		systemsmanager.MaintenanceWindows,
		systemsmanager.GetManagedEc2,
		ec2.VpcEndpoints,
		ec2.VpcPeeringConnections,
		kinesis.Streams,
		redshift.EventSubscriptions,
		codecommit.Repositories,
		codebuild.BuildProjects,
		codedeploy.DeploymentGroups,
		lightsail.Containers,
		lightsail.Databases,
		lightsail.VPS,
		neptune.Clusters,
		route53.HostedZones,
		cloudtrail.Trails,
		datasync.Agents,
		secretsmanager.Secrets,
		ec2.TransitGatewayPeeringAttachments,
		ec2.TransitGatewayVpcAttachments,
		firehose.DeliveryStreams,
		kinesisanalytics.KinesisAnalytics,
		kafka.Kafka,
		ec2.NatGateways,
		sns.Subscriptions,
		wafv2.WebAcls,
	}
}

func FetchResources(ctx context.Context, client providers.ProviderClient, regions []string, db *bun.DB, telemetry bool, analytics utils.Analytics, wp *providers.WorkerPool) {
	listOfSupportedRegions := getRegions()
	if len(regions) > 0 {
		log.Infof("Komiser will fetch resources from the following regions: %s", strings.Join(regions, ","))
		listOfSupportedRegions = regions
	}

	var costexplorerOutputList []*costexplorer.GetCostAndUsageOutput
	if jsonData, err := readCostExplorerCache(); err == nil {
		err := json.Unmarshal(jsonData, &costexplorerOutputList)
		if err != nil {
			log.Warn("Failed to unmarshal cached cost explorer data:", err)
		}
	} else {
		costexplorerOutputList, err = getCostexplorerOutput(
			ctx, client, utils.BeginningMonthsAgo(time.Now(), 6).Format("2006-01-02"), utils.EndingOfLastMonth(time.Now()).Format("2006-01-02"),
		)
		if err != nil {
			log.Warn("Failed to get cost explorer output:", err)
		}
		if err := writeCostExplorerCache(costexplorerOutputList); err != nil {
			log.Warn("Failed to write cost explorer cache:", err)
		}
		resources, err := awsUtils.ExtractResources(costexplorerOutputList)
		if err != nil {
			log.Warn("Failed to extract resources from cost explorer output:", err)
		}
		for i, resource := range resources {
			resource.ResourceId = "OLD_RESOURCE_" + strconv.Itoa(i)
			resource.Account = client.Name
			_, err = db.NewInsert().Model(&resource).On("CONFLICT (resource_id) DO UPDATE").Set("cost = EXCLUDED.cost, relations=EXCLUDED.relations").Exec(context.Background())
			if err != nil {
				log.WithError(err).Errorf("db trigger failed")
			}
		}
	}

	costexplorerOutputList, err := getCostexplorerOutput(ctx, client, utils.BeginningOfMonth(time.Now()).Format("2006-01-02"), time.Now().Format("2006-01-02"))
	if err != nil {
		log.Warn("Failed to get cost explorer output:", err)
	}
	ctxWithCostexplorerOutput := context.WithValue(ctx, awsUtils.CostexplorerKey, costexplorerOutputList)
	for _, region := range listOfSupportedRegions {
		c := client.AWSClient.Copy()
		c.Region = region
		client = providers.ProviderClient{
			AWSClient: &c,
			Name:      client.Name,
		}
		for _, fetchResources := range listOfSupportedServices() {
			fetchResources := fetchResources
			wp.SubmitTask(func() {
				resources, err := fetchResources(ctxWithCostexplorerOutput, client)
				if err != nil {
					log.Warnf("[%s][AWS] %s", client.Name, err)
				} else {
					for _, resource := range resources {
						_, err = db.NewInsert().Model(&resource).On("CONFLICT (resource_id) DO UPDATE").Set("cost = EXCLUDED.cost, relations=EXCLUDED.relations").Exec(context.Background())
						if err != nil {
							log.WithError(err).Errorf("db trigger failed")
						}
					}
					if telemetry {
						analytics.TrackEvent("discovered_resources", map[string]interface{}{
							"provider":     "AWS",
							"resources":    len(resources),
							"dependencies": calculateDependencies(resources),
						})
					}
				}
			})
		}
	}
}

func calculateDependencies(resources []models.Resource) int {
	total := 0
	for _, resource := range resources {
		total += len(resource.Relations)
	}
	return total
}

func getRegions() []string {
	return []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"}
}
