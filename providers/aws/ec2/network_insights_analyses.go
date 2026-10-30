package ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
)

func NetworkInsightsAnalyses(ctx context.Context, client providers.ProviderClient) ([]models.Resource, error) {
	resources := make([]models.Resource, 0)

	ec2Client := ec2.NewFromConfig(*client.AWSClient)
	input := &ec2.DescribeNetworkInsightsAnalysesInput{}

	for {
		output, err := ec2Client.DescribeNetworkInsightsAnalyses(ctx, input)
		if err != nil {
			return resources, err
		}

		for _, analysis := range output.NetworkInsightsAnalyses {
			tags := make([]models.Tag, 0)
			for _, tag := range analysis.Tags {
				tags = append(tags, models.Tag{
					Key:   aws.ToString(tag.Key),
					Value: aws.ToString(tag.Value),
				})
			}

			pathFound := "unknown"
			if analysis.NetworkPathFound != nil {
				pathFound = fmt.Sprint(*analysis.NetworkPathFound)
			}

			resources = append(resources, models.Resource{
				Provider:   "AWS",
				Account:    client.Name,
				Service:    "EC2 Network Insights Analysis",
				ResourceId: aws.ToString(analysis.NetworkInsightsAnalysisId),
				Region:     client.AWSClient.Region,
				Name:       aws.ToString(analysis.NetworkInsightsAnalysisId),
				Cost:       0,
				Tags:       tags,
				FetchedAt:  time.Now(),
				Link: fmt.Sprintf(
					"https://%s.console.aws.amazon.com/vpc/home?region=%s#NetworkInsightsAnalysis:NetworkInsightsAnalysisId=%s",
					client.AWSClient.Region, client.AWSClient.Region, aws.ToString(analysis.NetworkInsightsAnalysisId),
				),
				Metadata: map[string]string{
					"status":                string(analysis.Status),
					"networkPathFound":      pathFound,
					"networkInsightsPathId": aws.ToString(analysis.NetworkInsightsPathId),
				},
			})
		}

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	log.WithFields(log.Fields{
		"provider":  "AWS",
		"service":   "EC2 Network Insights Analysis",
		"account":   client.Name,
		"region":    client.AWSClient.Region,
		"resources": len(resources),
	}).Info("Fetched resources")

	return resources, nil
}
