package dto

type EventRequest struct {
	Type        string `json:"type"`
	Origin      string `json:"origin,omitempty"`
	Destination string `json:"destination,omitempty"`
	Amount      int64  `json:"amount"`
}

func NewEventRequest(t, origin, destination string, amount int64) EventRequest {
	return EventRequest{
		Type:        t,
		Origin:      origin,
		Destination: destination,
		Amount:      amount,
	}
}

func (eventType *EventRequest) GetType() string {
	return eventType.Type
}

func (eventOrigin *EventRequest) GetOrigin() string {
	return eventOrigin.Origin
}

func (eventDest *EventRequest) GetDestination() string {
	return eventDest.Destination
}

func (eventAmount *EventRequest) GetAmount() int64 {
	return eventAmount.Amount
}

func (eventType *EventRequest) SetType(t string) {
	eventType.Type = t
}

func (eventOrigin *EventRequest) SetOrigin(origin string) {
	eventOrigin.Origin = origin
}

func (eventDestination *EventRequest) SetDestination(destination string) {
	eventDestination.Destination = destination
}

func (eventAmount *EventRequest) SetAmount(amount int64) {
	eventAmount.Amount = amount
}
