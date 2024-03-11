package delivery

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/promotedai/schema/generated/go/proto/common"
	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/promotedai/schema/generated/go/proto/event"
)

const defaultDeliveryTimeoutMillis = 250
const defaultMetricsTimeoutMillis = 3000
const defaultMaxRequestInsertions = 1000

const serverVersion = "go.1.0.0"

// PromotedDeliveryClient is a client for interacting with the Promoted.ai Delivery API.
type PromotedDeliveryClient struct {
	deliveryAPI               DeliveryAPI
	metricsAPI                MetricsAPI
	sdkDelivery               DeliveryAPI
	deliveryEndpoint          string
	deliveryAPIKey            string
	deliveryTimeoutMillis     int64
	metricsEndpoint           string
	metricsAPIKey             string
	metricsTimeoutMillis      int64
	maxRequestInsertions      int
	applyTreatmentChecker     ApplyTreatmentChecker
	shadowTrafficDeliveryRate float32
	sampler                   Sampler
	performChecks             bool
	blockingShadowTraffic     bool
}

// Deliver sends a delivery request and returns the response.
func (client *PromotedDeliveryClient) Deliver(deliveryRequest *DeliveryRequest) (*DeliveryResponse, error) {
	plan := client.Plan(deliveryRequest.OnlyLog, deliveryRequest.Experiment)
	client.PrepareRequest(deliveryRequest, plan)

	var apiResponse *delivery.Response
	var err error
	if plan.UseAPIResponse {
		apiResponse, err = client.CallDeliveryAPI(apiResponse, err, deliveryRequest)
		if err != nil {
			log.Printf("Error calling Delivery API, falling back: %v\n", err)
		}
	}

	// Note this returns a delivery response based on this apiResponse if it's set, and creates
	// an SDK response otherwise.
	return client.HandleSDKAndLog(deliveryRequest, plan, apiResponse)
}

func (client *PromotedDeliveryClient) CallDeliveryAPI(apiResponse *delivery.Response, err error, deliveryRequest *DeliveryRequest) (*delivery.Response, error) {
	return client.deliveryAPI.RunDelivery(deliveryRequest)
}

// Plan returns a DeliveryPlan that determines SDK execution, always using SDK if we are
// only logging, and otherwise checking the experiment to decide.
func (client *PromotedDeliveryClient) Plan(onlyLog bool, experiment *event.CohortMembership) *DeliveryPlan {
	useApiResponse := !onlyLog && client.shouldApplyTreatment(experiment)
	return NewDeliveryPlan(client.generateClientID(), useApiResponse)
}

// PrepareRequest prepares the delivery request using the plan.
func (client *PromotedDeliveryClient) PrepareRequest(deliveryRequest *DeliveryRequest, plan *DeliveryPlan) {
	if client.performChecks {
		validationErrors := deliveryRequest.Validate()
		for _, validationError := range validationErrors {
			log.Printf("Delivery Request Validation Error: %s\n", validationError)
		}
	}
	client.ensureClientRequestID(deliveryRequest.Request, plan.ClientRequestID)
	client.fillInRequestFields(deliveryRequest.Request)
}

// HandleSDKAndLog handles SDK delivery, logs, and shadow traffic.
func (client *PromotedDeliveryClient) HandleSDKAndLog(deliveryRequest *DeliveryRequest, plan *DeliveryPlan, apiResponse *delivery.Response) (*DeliveryResponse, error) {
	cohortMembership := client.cloneCohortMembership(deliveryRequest.Experiment)

	var response *delivery.Response
	var execSrv delivery.ExecutionServer

	if apiResponse != nil {
		response = apiResponse
		execSrv = delivery.ExecutionServer_API
	} else {
		var err error
		response, err = client.sdkDelivery.RunDelivery(deliveryRequest)
		if err != nil {
			return nil, err
		}
		execSrv = delivery.ExecutionServer_SDK
	}

	// Log SDK DeliveryLog to Metrics API.
	if execSrv != delivery.ExecutionServer_API || cohortMembership != nil {
		client.logToMetrics(deliveryRequest, response, cohortMembership, execSrv)
	}

	// Send shadow traffic if needed.
	if !plan.UseAPIResponse && client.shouldSendShadowTraffic() {
		client.deliverShadowTraffic(deliveryRequest)
	}

	return &DeliveryResponse{
		Response:        response,
		ClientRequestID: plan.ClientRequestID,
		ExecutionServer: execSrv,
	}, nil
}

// deliverShadowTraffic sends shadow traffic, optionally asynchronously depending on client config.
func (client *PromotedDeliveryClient) deliverShadowTraffic(deliveryRequest *DeliveryRequest) {
	if client.blockingShadowTraffic {
		client.doDeliverShadowTraffic(deliveryRequest)
	} else {
		go client.doDeliverShadowTraffic(deliveryRequest)
	}
}

// doDeliverShadowTraffic actually sends shadow traffic.
func (client *PromotedDeliveryClient) doDeliverShadowTraffic(deliveryRequest *DeliveryRequest) {
	// Clone the request for safe modification.
	requestToSend := deliveryRequest.Clone(NoMaxRequestInsertions)

	// Ensure client info is filled in.
	if requestToSend.Request.ClientInfo == nil {
		requestToSend.Request.ClientInfo = &common.ClientInfo{}
	}
	requestToSend.Request.ClientInfo.ClientType = common.ClientInfo_PLATFORM_SERVER
	requestToSend.Request.ClientInfo.TrafficType = common.ClientInfo_SHADOW

	_, err := client.deliveryAPI.RunDelivery(requestToSend)
	if err != nil {
		log.Printf("Error calling Delivery API for shadow traffic: %v\n", err)
	}
}

// shouldApplyTreatment checks whether treatment should be applied to a cohort membership.
func (client *PromotedDeliveryClient) shouldApplyTreatment(cohortMembership *event.CohortMembership) bool {
	if client.applyTreatmentChecker != nil {
		return client.applyTreatmentChecker.ShouldApplyTreatment(cohortMembership)
	}
	if cohortMembership == nil {
		return true
	}
	return cohortMembership.Arm != event.CohortArm_CONTROL
}

// shouldSendShadowTraffic checks whether shadow traffic should be sent.
func (client *PromotedDeliveryClient) shouldSendShadowTraffic() bool {
	return client.shadowTrafficDeliveryRate > 0 && client.sampler.SampleRandom(client.shadowTrafficDeliveryRate)
}

// cloneCohortMembership clones a cohort membership.
func (client *PromotedDeliveryClient) cloneCohortMembership(cohortMembership *event.CohortMembership) *event.CohortMembership {
	if cohortMembership == nil {
		return nil
	}
	return &event.CohortMembership{
		Arm:      cohortMembership.Arm,
		CohortId: cohortMembership.CohortId,
	}
}

// logToMetrics logs to the Metrics API.
func (client *PromotedDeliveryClient) logToMetrics(deliveryRequest *DeliveryRequest, deliveryResponse *delivery.Response, cohortMembership *event.CohortMembership, execSrv delivery.ExecutionServer) {
	go func() {
		logRequest := client.createLogRequest(deliveryRequest, deliveryResponse, cohortMembership, execSrv)
		err := client.metricsAPI.RunMetricsLogging(logRequest)
		if err != nil {
			log.Printf("Error calling Metrics API: %v\n", err)
		}
	}()
}

// createLogRequest creates a log request from a delivery request/response.
func (client *PromotedDeliveryClient) createLogRequest(deliveryRequest *DeliveryRequest, deliveryResponse *delivery.Response, cohortMembershipToLog *event.CohortMembership, execSrv delivery.ExecutionServer) *event.LogRequest {
	logReq := &event.LogRequest{
		UserInfo:   deliveryRequest.Request.UserInfo,
		ClientInfo: deliveryRequest.Request.ClientInfo,
		PlatformId: deliveryRequest.Request.PlatformId,
		Timing:     deliveryRequest.Request.Timing,
	}

	if execSrv != delivery.ExecutionServer_API {
		deliveryLog := &delivery.DeliveryLog{
			Execution: &delivery.DeliveryExecution{
				ExecutionServer: execSrv,
				ServerVersion:   serverVersion,
			},
			Request:  deliveryRequest.Request,
			Response: deliveryResponse,
		}
		logReq.DeliveryLog = append(logReq.DeliveryLog, deliveryLog)
	}

	if cohortMembershipToLog != nil {
		logReq.CohortMembership = append(logReq.CohortMembership, cohortMembershipToLog)
	}

	return logReq
}

// fillInRequestFields fills in required fields on the request.
func (client *PromotedDeliveryClient) fillInRequestFields(req *delivery.Request) {
	if req.ClientInfo == nil {
		req.ClientInfo = &common.ClientInfo{}
	}
	clientInfo := req.ClientInfo
	clientInfo.ClientType = common.ClientInfo_PLATFORM_SERVER
	clientInfo.TrafficType = common.ClientInfo_PRODUCTION

	// Fill in client timestamp if not set by the caller.
	client.ensureClientTimestamp(req)
}

// ensureClientTimestamp ensures client timestamp is set on the request.
func (client *PromotedDeliveryClient) ensureClientTimestamp(req *delivery.Request) {
	if req.Timing == nil {
		req.Timing = &common.Timing{}
	}
	timing := req.Timing
	if timing.ClientLogTimestamp == 0 {
		timing.ClientLogTimestamp = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	}
}

// ensureClientRequestID ensures client request ID is set on the request.
func (client *PromotedDeliveryClient) ensureClientRequestID(req *delivery.Request, clientRequestID string) {
	if req.ClientRequestId == "" {
		req.ClientRequestId = clientRequestID
	}
}

// generateClientID generates a client ID.
func (client *PromotedDeliveryClient) generateClientID() string {
	return uuid.New().String()
}
