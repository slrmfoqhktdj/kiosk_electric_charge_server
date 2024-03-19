package main

import (
	"encoding/json"
	"errors"
	"fmt"

	v16 "github.com/aliml92/ocpp/v16"
)

const (
	MessageTypeIdCall       = 2
	MessageTypeIdCallResult = 3
	MessageTypeIdCallError  = 4
)

var errInvalidAction = errors.New("invalid action")

var reqmapv16, resmapv16 map[string]func(json.RawMessage) (Payload, error)

type ocppError struct {
	id    string
	code  string
	cause string
}

func (e *ocppError) Error() string {
	return e.code + ": " + e.cause
}

// Call represents OCPP Call
type Call struct {
	MessageTypeId uint8
	UniqueId      string
	Action        string
	Payload       Payload
}

func (c *Call) getID() string {
	return c.UniqueId
}

// Create CallResult from a received Call
func (call *Call) createCallResult(r Payload) []byte {
	out := [3]interface{}{
		3,
		call.UniqueId,
		r,
	}
	raw, _ := json.Marshal(out)
	return raw
}

// Creates a CallError from a received Call
// TODO: organize error codes
func (call *Call) createCallError(err error) []byte {
	var id, code, cause string
	var ocppErr *ocppError
	if errors.As(err, &ocppErr) {
		id = ocppErr.id
		code = ocppErr.code
		cause = ocppErr.cause
	}
	if id == "" {
		id = "-1"
	}
	callError := &CallError{
		UniqueId:     id,
		ErrorCode:    code,
		ErrorDetails: cause,
	}
	switch code {
	case "ProtocolError":
		callError.ErrorDescription = "Payload for Action is incomplete"
	case "PropertyConstraintViolation":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "NotImplemented":
		callError.ErrorDescription = "Requested Action is not known by receiver"
	case "TypeConstraintViolationError":
		callError.ErrorDescription = "Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12)"
	case "PropertyConstraintViolationError":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "MessageTypeNotSupported":
		callError.ErrorDescription = "A message with an Message Type Number received that is not supported by this implementation"
	default:
		callError.ErrorDescription = "Unknown error"
	}
	return callError.marshal()
}

// CallResult represents OCPP CallResult
type CallResult struct {
	MessageTypeId uint8
	UniqueId      string
	Payload       json.RawMessage
}

func (cr *CallResult) getID() string {
	return cr.UniqueId
}

// CallError represents OCPP CallError
type CallError struct {
	MessageTypeId    uint8
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     interface{}
}

func (ce *CallError) marshal() []byte {
	ed := ce.ErrorDetails.(string)
	out := [5]interface{}{
		4,
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"cause":` + ed + `}`,
	}
	raw, _ := json.Marshal(out)
	return raw
}

func (ce *CallError) Error() string {
	return fmt.Sprintf("CallError: UniqueId=%s, ErrorCode=%s, ErrorDescription=%s, ErrorDetails=%s",
		ce.UniqueId, ce.ErrorCode, ce.ErrorDescription, ce.ErrorDetails)
}

func (ce *CallError) getID() string {
	return ce.UniqueId
}

type OcppMessage interface {
	getID() string
}

// umpack converts json byte to one of the ocpp messages, if not successful
// returns an error
// unpack expects ocpp messages in the below forms:
//   - [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}] -> Call
//   - [<MessageTypeId>, "<UniqueId>", {<Payload>}]             -> CallResult
//   - [<MessageTypeId>, "<UniqueId>", "<ErrorCode>", "<ErrorDescription>" , {<ErrorDetails>}] -> CallError
//
// if json byte does not conform with these three formats, OcppError is created
// and is returned as error
//
// TODO: improve error construction depending on protocol
func unpack(b []byte, proto string) (OcppMessage, error) {
	var rm []json.RawMessage //  raw message
	var mti uint8            //  MessageTypeId
	var ui string            //  UniqueId
	var a string             //  Action
	var p Payload            //  Payload
	var ocppMsg OcppMessage  //  OcppMessage
	var e *ocppError

	// unmarshal []byte to []json.RawMessage
	err := json.Unmarshal(b, &rm)
	if err != nil {
		e = &ocppError{
			id:    "-1",
			code:  "ProtocolError",
			cause: "Invalid JSON format",
		}
		return nil, e
	}

	// unmarshal [0]json.RawMessage to MessageTypeId
	err1 := json.Unmarshal(rm[0], &mti)
	// unmarshal [1]json.RawMessage to UniqueId
	err2 := json.Unmarshal(rm[1], &ui)
	if err1 != nil {
		if err2 != nil {
			e.id = "-1"
			e.cause = e.cause + "," + fmt.Sprintf("UniqueId: %v is not valid", rm[1])
		} else {
			e.id = ui
		}
		return nil, e
	}
	if len(ui) > 36 {
		e = &ocppError{
			code:  "ProtocolError",
			cause: fmt.Sprintf("UniqueId: %v is too long", ui),
		}
	}
	switch mti {
	case MessageTypeIdCall:
		call := &Call{
			MessageTypeId: mti,
			UniqueId:      ui,
		}
		if e != nil {
			return call, e
		}
		err = json.Unmarshal(rm[2], &a)
		if err != nil {
			e := &ocppError{
				id:    ui,
				code:  "ProtocolError",
				cause: "Message does not contain Action",
			}
			return call, e
		}
		call.Action = a

		p, err = unmarshalRequestPayload(a, rm[3], proto)

		var ocppErr *ocppError
		if err != nil {
			if errors.As(err, &ocppErr) {
				ocppErr.id = ui
			}
			return call, err
		}
		call.Payload = p
		ocppMsg = call
	case MessageTypeIdCallResult:
		p := rm[2]
		ocppMsg = &CallResult{
			MessageTypeId: mti,
			UniqueId:      ui,
			Payload:       p,
		}
	case MessageTypeIdCallError:
		var me [5]interface{}
		_ = json.Unmarshal(b, &me)
		ocppMsg = &CallError{
			MessageTypeId:    mti,
			UniqueId:         ui,
			ErrorCode:        me[2].(string),
			ErrorDescription: me[3].(string),
			ErrorDetails:     me[4],
		}
	default:
		e := &ocppError{
			id:    ui,
			code:  "MessageTypeNotSupported",
			cause: fmt.Sprintf("A message with: %v is not supported by this implementation", mti),
		}
		return nil, e
	}
	return ocppMsg, nil
}

// unmarshalRequestPayload unmarshals raw bytes of request type payload
// to a corresponding struct depending on Action and ocpp protocol
func unmarshalRequestPayload(actionName string, rawPayload json.RawMessage, proto string) (Payload, error) {
	var uf func(json.RawMessage) (Payload, error) // uf unmarshal function for a specific action request
	var ok bool

	switch proto {
	case OcppV16:
		uf, ok = reqmapv16[actionName]
	}
	if !ok {
		e := &ocppError{
			code:  "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", actionName),
		}
		return nil, e
	}
	return uf(rawPayload)
}

// unmarshalRequestPayloadv16 unmarshals raw request type payload to a ***Req type struct of ocppv16
func unmarshalRequestPayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload

	err := json.Unmarshal(rawPayload, &p)

	// log.Print(rawPayload)

	if err != nil {
		e := &ocppError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}

		return nil, e
	}

	// log.Print(p)

	err = ValidateV16.Struct(p)
	if err != nil {
		// TODO: construct more detailed error
		e := &ocppError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}

		return nil, e
	}
	payload = &p

	return payload, nil
}

// unmarshalResponsePv16 unmarshals raw bytes of request type payload
// to a corresponding struct depending on actionName
func unmarshalResponsePv16(actionName string, rawPayload json.RawMessage) (Payload, error) {
	uf, ok := resmapv16[actionName] // uf unmarshal function for a specific action request
	if !ok {
		return nil, errInvalidAction
	}
	return uf(rawPayload)
}

// unmarshalResponsePayloadv16 unmarshals raw response type payload to a ***Conf type struct of ocppv16
func unmarshalResponsePayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = ValidateV16.Struct(p)
	if err != nil {
		return nil, err
	}
	payload = &p
	return payload, nil
}

type StopTransactionRequestTT struct {
	IdTag           string         `json:"idTag,omitempty" validate:"omitempty,max=20"`
	MeterStop       int            `json:"meterStop" validate:"required,gte=0"`
	Timestamp       string         `json:"timestamp" validate:"required,ISO8601date"`
	TransactionId   int            `json:"transactionId" validate:"required"`
	Reason          string         `json:"reason,omitempty" validate:"omitempty,Reason"`
	TransactionData []MeterValueTT `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

type StatusNotificationReqTT struct {
	ConnectorId     int    `json:"connectorId" validate:"required,gte=0"`
	ErrorCode       string `json:"errorCode"  validate:"required,ChargePointErrorCode"`
	Info            string `json:"info,omitempty"  validate:"max=50"`
	Status          string `json:"status"  validate:"required,ChargePointStatus"`
	Timestamp       string `json:"timestamp,omitempty" validate:"omitempty,ISO8601date"`
	VendorId        string `json:"vendorId,omitempty"  validate:"max=255"`
	VendorErrorCode string `json:"vendorErrorCode,omitempty"  validate:"max=50"`
}

type MeterValuesReqTT struct {
	ConnectorId   int            `json:"connectorId" validate:"required,gte=0"`
	TransactionId int            `json:"transactionId"`
	MeterValue    []MeterValueTT `json:"meterValue" validate:"required,dive"`
}

type MeterValueTT struct {
	Timestamp    string           `json:"timestamp" validate:"required,ISO8601date"`
	SampledValue []SampledValueTT `json:"sampledValue" validate:"required,dive"`
}

type SampledValueTT struct {
	Value     string `json:"value" validate:"required"`
	Context   string `json:"context,omitempty" validate:"omitempty,ReadingContext"`
	Format    string `json:"format,omitempty" validate:"omitempty,ValueFormat"`
	Measurand string `json:"measurand,omitempty" validate:"omitempty,Measurand"`
	Phase     string `json:"phase,omitempty" validate:"omitempty,Phase"`
	Location  string `json:"location,omitempty" validate:"omitempty,Location"`
	Unit      string `json:"unit,omitempty" validate:"omitempty,UnitOfMeasure"`
}

type ChangeConfigurationReqTT struct {
	Key   string `json:"key" validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=500"`
}

type RemoteStartTransactionReqTT struct {
	ConnectorId     int                  `json:"connectorId,omitempty" validate:"gte=0"`
	IdTag           string               `json:"idTag" validate:"required,max=20"`
	ChargingProfile *v16.ChargingProfile `json:"chargingProfile,omitempty" validate:"omitempty"`
}

type GetConfigurationConfTT struct {
	ConfigurationKey []ConfigurationKey `json:"configurationKey,omitempty" validate:"omitempty,dive"`
	UnknownKey       []string           `json:"unknownKey,omitempty" validate:"omitempty,dive,max=50"`
}

type ConfigurationKey struct {
	Key      string  `json:"key" validate:"required,max=50"`
	Readonly bool    `json:"readonly"`
	Value    *string `json:"value,omitempty" validate:"omitempty,max=500"`
}

func init() {

	reqmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              unmarshalRequestPayloadv16[v16.BootNotificationReq],
		"Authorize":                     unmarshalRequestPayloadv16[v16.AuthorizeReq],
		"DataTransfer":                  unmarshalRequestPayloadv16[v16.DataTransferReq],
		"DiagnosticsStatusNotification": unmarshalRequestPayloadv16[v16.DiagnosticsStatusNotificationReq],
		"FirmwareStatusNotification":    unmarshalRequestPayloadv16[v16.FirmwareStatusNotificationReq],
		"Heartbeat":                     unmarshalRequestPayloadv16[v16.HeartbeatReq],
		"MeterValues":                   unmarshalRequestPayloadv16[MeterValuesReqTT],
		"StartTransaction":              unmarshalRequestPayloadv16[v16.StartTransactionReq],
		"StatusNotification":            unmarshalRequestPayloadv16[StatusNotificationReqTT],
		"StopTransaction":               unmarshalRequestPayloadv16[StopTransactionRequestTT],
		"CanCelReservation":             unmarshalRequestPayloadv16[v16.CancelReservationReq],
		"ChangeAvailability":            unmarshalRequestPayloadv16[v16.ChangeAvailabilityReq],
		"ChangeConfiguration":           unmarshalRequestPayloadv16[ChangeConfigurationReqTT],
		"ClearCache":                    unmarshalRequestPayloadv16[v16.ClearCacheReq],
		"ClearChargingProfile":          unmarshalRequestPayloadv16[v16.ClearChargingProfileReq],
		"GetCompositeSchedule":          unmarshalRequestPayloadv16[v16.GetCompositeScheduleReq],
		"GetConfiguration":              unmarshalRequestPayloadv16[v16.GetConfigurationReq],
		"GetDiagnostics":                unmarshalRequestPayloadv16[v16.GetDiagnosticsReq],
		"GetLocalListVersion":           unmarshalRequestPayloadv16[v16.GetLocalListVersionReq],
		"RemoteStartTransaction":        unmarshalRequestPayloadv16[RemoteStartTransactionReqTT],
		"RemoteStopTransaction":         unmarshalRequestPayloadv16[v16.RemoteStopTransactionReq],
		"ReserveNow":                    unmarshalRequestPayloadv16[v16.ReserveNowReq],
		"Reset":                         unmarshalRequestPayloadv16[v16.ResetReq],
		"SendLocalList":                 unmarshalRequestPayloadv16[v16.SendLocalListReq],
		"SetChargingProfile":            unmarshalRequestPayloadv16[v16.SetChargingProfileReq],
		"TriggerMessage":                unmarshalRequestPayloadv16[v16.TriggerMessageReq],
		"UnlockConnector":               unmarshalRequestPayloadv16[v16.UnlockConnectorReq],
		"UpdateFirmware":                unmarshalRequestPayloadv16[v16.UpdateFirmwareReq],
	}

	resmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              unmarshalResponsePayloadv16[v16.BootNotificationConf],
		"Authorize":                     unmarshalResponsePayloadv16[v16.AuthorizeConf],
		"DataTransfer":                  unmarshalResponsePayloadv16[v16.DataTransferConf],
		"DiagnosticsStatusNotification": unmarshalResponsePayloadv16[v16.DiagnosticsStatusNotificationConf],
		"FirmwareStatusNotification":    unmarshalResponsePayloadv16[v16.FirmwareStatusNotificationConf],
		"Heartbeat":                     unmarshalResponsePayloadv16[v16.HeartbeatConf],
		"MeterValues":                   unmarshalResponsePayloadv16[v16.MeterValuesConf],
		"StartTransaction":              unmarshalResponsePayloadv16[v16.StartTransactionConf],
		"StatusNotification":            unmarshalResponsePayloadv16[v16.StatusNotificationConf],
		"StopTransaction":               unmarshalResponsePayloadv16[v16.StopTransactionConf],
		"CancelReservation":             unmarshalResponsePayloadv16[v16.CancelReservationConf],
		"ChangeAvailability":            unmarshalResponsePayloadv16[v16.ChangeAvailabilityConf],
		"ChangeConfiguration":           unmarshalResponsePayloadv16[v16.ChangeConfigurationConf],
		"ClearCache":                    unmarshalResponsePayloadv16[v16.ClearCacheConf],
		"ClearChargingProfile":          unmarshalResponsePayloadv16[v16.ClearChargingProfileConf],
		"GetCompositeSchedule":          unmarshalResponsePayloadv16[v16.GetCompositeScheduleConf],
		"GetConfiguration":              unmarshalResponsePayloadv16[GetConfigurationConfTT],
		"GetDiagnostics":                unmarshalResponsePayloadv16[v16.GetDiagnosticsConf],
		"GetLocalListVersion":           unmarshalResponsePayloadv16[v16.GetLocalListVersionConf],
		"RemoteStartTransaction":        unmarshalResponsePayloadv16[v16.RemoteStartTransactionConf],
		"RemoteStopTransaction":         unmarshalResponsePayloadv16[v16.RemoteStopTransactionConf],
		"ReserveNow":                    unmarshalResponsePayloadv16[v16.ReserveNowConf],
		"Reset":                         unmarshalResponsePayloadv16[v16.ResetConf],
		"SendLocalList":                 unmarshalResponsePayloadv16[v16.SendLocalListConf],
		"SetChargingProfile":            unmarshalResponsePayloadv16[v16.SetChargingProfileConf],
		"TriggerMessage":                unmarshalResponsePayloadv16[v16.TriggerMessageConf],
		"UnlockConnector":               unmarshalResponsePayloadv16[v16.UnlockConnectorConf],
		"UpdateFirmware":                unmarshalResponsePayloadv16[v16.UpdateFirmwareConf],
	}

}
