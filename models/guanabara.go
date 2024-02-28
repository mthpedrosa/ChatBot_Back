package models

type Location struct {
	ID               int      `json:"id"`
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	CityName         string   `json:"cityName"`
	StateCode        string   `json:"stateCode"`
	CountryCode      string   `json:"countryCode"`
	Lat              string   `json:"lat"`
	Lon              string   `json:"lon"`
	PointsOfInterest []string `json:"pointsOfInterest"`
}

type LocationPair struct {
	OriginID      int `json:"origemId"`
	DestinationID int `json:"destinoId"`
}

type LocationResponse struct {
	Data struct {
		Locations     []Location     `json:"locations"`
		LocationsPair []LocationPair `json:"localidades"`
		Groups        []Group        `json:"groups"`
	} `json:"data"`
}

type Group struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Locations []struct {
		Id int `json:"id"`
	} `json:"locations"`
}

type TripRequest struct {
	Departure    string `json:"departure"`
	Arrival      string `json:"arrival"`
	Data         string `json:"dataTrip"`
	Localizador  string `json:"localizador"`
	City         string `json:"city"`
	DepartureId  string `json:"departureID"`
	ArrivalId    string `json:"arrivalID"`
	FirstMessage bool   `json:"first_message"`
}

type ConnectionRoute struct {
	ControlNumber      string   `json:"controlNumber"`
	DepartureTime      string   `json:"departureTime"`
	ArrivalTime        string   `json:"arrivalTime"`
	ArrivalDay         int      `json:"arrivalDay"`
	IDCompany          int      `json:"idCompany"`
	CompanyName        string   `json:"companyName"`
	ClassOfServiceName string   `json:"classOfServiceName"`
	BPE                bool     `json:"bpe"`
	DepartureLocation  Location `json:"departureLocation"`
	ArrivalLocation    Location `json:"arrivalLocation"`
}

type TripInfoResponse struct {
	ControlNumber     string `json:"controlNumber"`
	DepartureTime     string `json:"departureTime"`
	DepartureLocation struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"departureLocation"`
	ArrivalTime     string `json:"arrivalTime"`
	ArrivalLocation struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"arrivalLocation"`
	ArrivalDay                 int               `json:"arrivalDay"`
	IDCompany                  int               `json:"idCompany"`
	CompanyName                string            `json:"companyName"`
	PriceValue                 float64           `json:"priceValue"`
	Prices                     []interface{}     `json:"prices"`
	ClassOfServiceName         string            `json:"classOfServiceName"`
	BPE                        bool              `json:"bpe"`
	HasConnection              bool              `json:"hasConnection"`
	AvailableSeats             int               `json:"availableSeats"`
	ConnectionRoutes           []ConnectionRoute `json:"connectionRoutes"`
	CurrencyCode               string            `json:"currencyCode"`
	Duration                   string            `json:"duration"`
	IDDailySchedule            int               `json:"idDailySchedule"`
	IDSchedule                 int               `json:"idSchedule"`
	ScheduleControlNumber      string            `json:"scheduleControlNumber"`
	ControlNumberDailySchedule string            `json:"controlNumberDailySchedule"`
}

type TripResponse struct {
	Data []TripInfoResponse `json:"data"`
}

type TripRequestAPI struct {
	DepartureLocation      int    `json:"departureLocation"`
	ArrivalLocation        int    `json:"arrivalLocation"`
	DepartureDate          string `json:"departureDate"`
	DepartureLocationGroup int    `json:"departureLocationGroup"`
	ArrivalLocationGroup   int    `json:"arrivalLocationGroup"`
	PromoCode              string `json:"promoCode"`
	PassengerType          int    `json:"passengerType"`
}

type GetTrips struct {
	DepartureID        int
	DestinationID      int
	DepartureGroupID   int
	DestinationGroupID int
	DepartureDate      string
}
