package portfolio

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
	"github.com/awe76/saga-coordinator/workflow"
)

type CreateBuildingMapPortfolio struct {
	UUID                     string   `json:"uuid"`
	CompanyUUID              string   `json:"company_uuid"`
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	LogoUrl                  string   `json:"logo_url"`
	PortfolioLogoBuildingIds []string `json:"portfolio_building_ids"`
}

type CreatePortfolio struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	LogoUrl     string   `json:"logo_url"`
	BuildingIds []string `json:"building_ids"`
}

type BuildingResponse struct {
	UUID string `json:"uuid"`
}

type PortfolioRo struct {
	ID          int                `json:"id"`
	UUID        string             `json:"uuid"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	LogoUrl     string             `json:"logo_url"`
	Buildings   []BuildingResponse `json:"buildings"`
}

const level = 0.8

func HandleCreateBuildingMapPotrfolio(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", op)

	if rand.Float32() < level {
		op.Payload = &PortfolioRo{
			ID:          1,
			UUID:        "test-uuid",
			Name:        "test-name",
			Description: "test-description",
			LogoUrl:     "test-logo-url",
			Buildings: []BuildingResponse{
				{
					UUID: "uuid-1",
				},
				{
					UUID: "uuid-2",
				},
			},
		}

		return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
	} else {
		return p.SendMessage(workflow.WORKFLOW_OPERATION_FAILED, op)
	}
}

func HandleCreateBuildingMapPotrfolioRollback(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
}

func HandleCreatePotfolio(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", op)

	if rand.Float32() < level {
		op.Payload = &CreatePortfolio{
			UUID:        "test-uuid",
			Name:        "test-name",
			Description: "test-description",
			LogoUrl:     "test-logo-url",
			BuildingIds: []string{"uuid-1", "uuid-2"},
		}
		return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
	} else {
		return p.SendMessage(workflow.WORKFLOW_OPERATION_FAILED, op)
	}
}

func HandleCreatePotfolioRollback(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
}
