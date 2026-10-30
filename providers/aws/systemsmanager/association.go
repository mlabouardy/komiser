package systemsmanager

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
)

func Associations(ctx context.Context, client providers.ProviderClient) ([]models.Resource, error) {
	resources := make([]models.Resource, 0)
	ssmClient := ssm.NewFromConfig(*client.AWSClient)

	input := &ssm.ListAssociationsInput{}

	for {
		associations, err := ssmClient.ListAssociations(ctx, input)
		if err != nil {
			return resources, err
		}

		for _, association := range associations.Associations {
			name := aws.ToString(association.AssociationName)
			if name == "" {
				name = aws.ToString(association.Name)
			}

			resources = append(resources, models.Resource{
				Provider:   "AWS",
				Account:    client.Name,
				Service:    "SSM Association",
				ResourceId: aws.ToString(association.AssociationId),
				Name:       name,
				Region:     client.AWSClient.Region,
				Cost:       0,
				FetchedAt:  time.Now(),
				Link: fmt.Sprintf(
					"https://%s.console.aws.amazon.com/systems-manager/state-manager/%s",
					client.AWSClient.Region, aws.ToString(association.AssociationId),
				),
				Metadata: map[string]string{
					"documentName":       aws.ToString(association.Name),
					"documentVersion":    aws.ToString(association.DocumentVersion),
					"associationVersion": aws.ToString(association.AssociationVersion),
				},
			})
		}

		if associations.NextToken == nil {
			break
		}

		input.NextToken = associations.NextToken
	}

	log.WithFields(log.Fields{
		"provider":  "AWS",
		"account":   client.Name,
		"region":    client.AWSClient.Region,
		"service":   "SSM Association",
		"resources": len(resources),
	}).Info("Fetched resources")

	return resources, nil
}
