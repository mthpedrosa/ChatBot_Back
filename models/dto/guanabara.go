package dto

type TripInfoResponseDTO struct {
	ControlNumber         string `json:"controlNumber"`
	DepartureTime         string `json:"departureTime"`
	DepartureLocationName string `json:"departureName"`
	/*DepartureLocation struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"departureLocation"`*/
	ArrivalTime         string `json:"arrivalTime"`
	ArrivalLocationName string `json:"arrivalName"`
	/*ArrivalLocation struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"arrivalLocation"`*/
	ArrivalDay         int     `json:"arrivalDay"`
	IDCompany          int     `json:"idCompany"`
	CompanyName        string  `json:"companyName"`
	PriceValue         float64 `json:"priceValue"`
	Duration           string  `json:"duration"`
	ClassOfServiceName string  `json:"classOfServiceName"`
}
