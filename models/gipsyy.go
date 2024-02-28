package models

type LocationGipsyy struct {
	Name        string `json:"name"`
	ExternalID  int    `json:"external_id"`
	Description string `json:"description"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

type Seat struct {
	SeatIdentifier    string      `json:"seatIdentifier"`
	Row               int         `json:"row"`
	Cell              int         `json:"cell"`
	Level             int         `json:"level"`
	IsUnavailable     bool        `json:"isUnavailable"`
	PricingIdentifier interface{} `json:"pricingIdentifier"`
	Price             float64     `json:"price"`
}

type PendingConnection struct {
	ControlNumber              string         `json:"controlNumber"`
	DepartureTime              string         `json:"departureTime"`
	ArrivalTime                string         `json:"arrivalTime"`
	DepartureLocation          LocationGipsyy `json:"departureLocation"`
	ArrivalLocation            LocationGipsyy `json:"arrivalLocation"`
	ArrivalDay                 int            `json:"arrivalDay"`
	IDCompany                  int            `json:"idCompany"`
	CompanyName                string         `json:"companyName"`
	PriceValue                 float64        `json:"priceValue"`
	PromoCodeError             bool           `json:"promoCodeError"`
	PromoCodeErrorMessage      interface{}    `json:"promoCodeErrorMessage"`
	ClassOfServiceName         string         `json:"classOfServiceName"`
	ServiceType                interface{}    `json:"serviceType"`
	CurrencyCode               interface{}    `json:"currencyCode"`
	Duration                   string         `json:"duration"`
	OriginalPrice              interface{}    `json:"originalPrice"`
	Seats                      []Seat         `json:"seats"`
	Extras                     interface{}    `json:"extras"`
	Sections                   interface{}    `json:"sections"`
	IDDailySchedule            int            `json:"idDailySchedule"`
	IDSchedule                 int            `json:"idSchedule"`
	ScheduleControlNumber      string         `json:"scheduleControlNumber"`
	ControlNumberDailySchedule string         `json:"controlNumberDailySchedule"`
	ConnectionIndex            int            `json:"connection_index"`
}

type Payment struct {
	ID                   string `json:"id"`
	BookingID            string `json:"booking_id"`
	Price                string `json:"price"`
	UserID               string `json:"user_id"`
	Status               string `json:"status"`
	ReservationsExpireAt string `json:"reservations_expire_at"`
}

type Ticket struct {
	ID                      string         `json:"id"`
	BookingID               string         `json:"booking_id"`
	PassengerName           string         `json:"passenger_name"`
	PassengerDocumentNumber string         `json:"passenger_document_number"`
	PassengerDocumentType   string         `json:"passenger_document_type"`
	Price                   string         `json:"price"`
	TripPrice               string         `json:"trip_price"`
	Status                  string         `json:"status"`
	DepartureAt             string         `json:"departure_at"`
	ArrivalAt               string         `json:"arrival_at"`
	DepartureAtLocal        string         `json:"departure_at_local"`
	ArrivalAtLocal          string         `json:"arrival_at_local"`
	SeatIdentifier          string         `json:"seat_identifier"`
	ConnectionIndex         int            `json:"connection_index"`
	PassengerType           string         `json:"passenger_type"`
	ArrivalLocation         LocationGipsyy `json:"arrival_location"`
	DepartureLocation       LocationGipsyy `json:"departure_location"`
}

type TripInfo struct {
	ID                string            `json:"id"`
	Price             string            `json:"price"`
	Return            bool              `json:"return"`
	DepartureAt       string            `json:"departure_at"`
	ArrivalAt         string            `json:"arrival_at"`
	DepartureAtLocal  string            `json:"departure_at_local"`
	ArrivalAtLocal    string            `json:"arrival_at_local"`
	ControlNumber     string            `json:"control_number"`
	PendingConnection PendingConnection `json:"pending_connection"`
	ClassOfService    string            `json:"class_of_service"`
	BabyPrice         string            `json:"baby_price"`
	ChildPrice        string            `json:"child_price"`
	ArrivalLocation   LocationGipsyy    `json:"arrival_location"`
	DepartureLocation LocationGipsyy    `json:"departure_location"`
	Payment           Payment           `json:"payment"`
	Tickets           []Ticket          `json:"tickets"`
}
