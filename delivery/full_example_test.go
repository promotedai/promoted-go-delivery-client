package delivery

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/promotedai/schema/generated/go/proto/common"
	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/promotedai/schema/generated/go/proto/event"
)

type Product struct {
	ID    int
	Name  string
	Price int
}

type MockDeliveryAPI struct {
	mock *gomock.Controller
	resp *delivery.Response
}

type MockMetricsAPI struct {
	mock *gomock.Controller
}

func (m *MockDeliveryAPI) RunDelivery(req *DeliveryRequest) (*delivery.Response, error) {
	return m.resp, nil
}

func (m *MockMetricsAPI) RunMetricsLogging(req *event.LogRequest) error {
	return nil
}

func TestFullExample(t *testing.T) {
	// Set up mocks for Promoted APIs.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeliveryAPI := &MockDeliveryAPI{mock: ctrl}
	mockMetricsAPI := &MockMetricsAPI{mock: ctrl}

	client := &PromotedDeliveryClient{
		deliveryAPI: mockDeliveryAPI,
		metricsAPI:  mockMetricsAPI,
	}

	// Mock a response from the delivery API where the products are re-ranked (reverse order).
	mockDeliveryAPI.resp = &delivery.Response{
		Insertion: []*delivery.Insertion{{ContentId: "2"}, {ContentId: "1"}},
	}

	// Retrieve products
	products := getProducts()

	// Create a map of products to reorder after Promoted ranking.
	productsMap := make(map[int]Product, len(products))

	// Create insertions for each product for the Promtoed delivery request.
	insertions := make([]*delivery.Insertion, 0, len(products))
	for _, product := range products {
		insertions = append(insertions, &delivery.Insertion{ContentId: strconv.Itoa(product.ID)})
		productsMap[product.ID] = product
	}

	// Create a Promoted delivery request.
	req := &DeliveryRequest{
		Request: &delivery.Request{
			PlatformId: 0,
			UserInfo:   &common.UserInfo{AnonUserId: "12355"},
			Paging: &delivery.Paging{
				Size:     100,
				Starting: &delivery.Paging_Offset{Offset: 0},
			},
			Insertion: insertions,
		},
	}

	// Call the Promoted delivery API.
	response, err := client.Deliver(req)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Apply Promoted's re-ranking to the products.
	rerankedProducts := make([]Product, 0, len(products))
	for _, insertion := range response.Response.Insertion {
		id, _ := strconv.Atoi(insertion.ContentId)
		rerankedProducts = append(rerankedProducts, productsMap[id])
	}

	// Verify the re-ranked products.
	assert.Equal(t, 2, rerankedProducts[0].ID)
	assert.Equal(t, 1, rerankedProducts[1].ID)

	assert.Equal(t, "Product 2", rerankedProducts[0].Name)
	assert.Equal(t, "Product 1", rerankedProducts[1].Name)

	// Go ahead and return results to the caller of your API.
}

func getProducts() []Product {
	return []Product{
		{
			ID:    1,
			Name:  "Product 1",
			Price: 100,
		},
		{
			ID:    2,
			Name:  "Product 2",
			Price: 200,
		},
	}
}
