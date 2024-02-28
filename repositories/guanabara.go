package repositories

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type guanabaraClient struct {
	httpClient *resty.Client
	cache      *redis.Client
}

func NewGuanabaraRepository(cache *redis.Client) interfaces.GuanabaraClientRepository {
	client := resty.New().
		SetBaseURL(viper.GetString("GUANABARA_URL")).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+viper.GetString("GUANABARA_KEY"))

	return &guanabaraClient{
		httpClient: client,
		cache:      cache,
	}
}

func (g *guanabaraClient) GetLocations(ctx context.Context) ([]models.Location, []models.Group, error) {
	cachedLocations, cachedGroups, err := g.getCachedLocations(ctx)
	if err == nil {
		return cachedLocations, cachedGroups, nil
	}

	return g.fetchAndCacheLocations(ctx)
}

func (g *guanabaraClient) GetLocationsPair(ctx context.Context) ([]models.LocationPair, error) {
	cachedLocationsPair, err := g.getCachedLocationsPair(ctx)
	if err == nil {
		return cachedLocationsPair, nil
	}

	return g.fetchAndCacheLocationsPair(ctx)
}

func (g *guanabaraClient) GetTrips(ctx context.Context, d models.GetTrips) ([]models.TripInfoResponse, error) {
	res, err := g.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"departureLocation":      d.DepartureID,
			"arrivalLocation":        d.DestinationID,
			"departureDate":          d.DepartureDate,
			"departureLocationGroup": d.DepartureGroupID,
			"arrivalLocationGroup":   d.DestinationGroupID,
			"promoCode":              "",
			"passengerType":          0,
		}).
		SetDebug(true).
		Post("/externalsale/getTrips")

	if err != nil {
		return nil, fmt.Errorf("error fetching trips: %v", err)
	}

	if res.Error() != nil {
		return nil, fmt.Errorf("error in response: %v", res.Error())
	}

	var trips models.TripResponse
	err = json.Unmarshal(res.Body(), &trips)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling trips: %v", err)
	}

	return trips.Data, nil
}

func (g *guanabaraClient) fetchAndCacheLocationsPair(ctx context.Context) ([]models.LocationPair, error) {
	res, err := g.httpClient.R().
		SetContext(ctx).
		Get("/externalsale/getLocationsPair")

	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return nil, fmt.Errorf("error in response: %v", res.Error())
	}

	var locations models.LocationResponse
	err = json.Unmarshal(res.Body(), &locations)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling locations: %v", err)
	}

	// Cache the response body for future use
	err = g.cache.Set(ctx, "@chatbot:guanabara:locationsPair", res.Body(), 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("error caching locations: %v", err)
	}

	return locations.Data.LocationsPair, nil
}

func (g *guanabaraClient) getCachedLocationsPair(ctx context.Context) ([]models.LocationPair, error) {
	locationsAsBytes, err := g.cache.Get(ctx, "@chatbot:guanabara:locationspair").Bytes()
	if err != nil {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("empty locations")
	}

	var locations models.LocationResponse
	err = json.Unmarshal(locationsAsBytes, &locations)
	if err != nil {
		return nil, err
	}

	return locations.Data.LocationsPair, nil
}

func (g *guanabaraClient) getCachedLocations(ctx context.Context) ([]models.Location, []models.Group, error) {
	locationsAsBytes, err := g.cache.Get(ctx, "@chatbot:guanabara:locations").Bytes()
	if err != nil {
		return nil, nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil, fmt.Errorf("empty locations")
	}

	var locations models.LocationResponse
	err = json.Unmarshal(locationsAsBytes, &locations)
	if err != nil {
		return nil, nil, err
	}

	return locations.Data.Locations, locations.Data.Groups, nil
}

func (g *guanabaraClient) fetchAndCacheLocations(ctx context.Context) ([]models.Location, []models.Group, error) {
	res, err := g.httpClient.R().
		SetContext(ctx).
		Get("/externalsale/getLocations")

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return nil, nil, fmt.Errorf("error in response: %v", res.Error())
	}

	var locations models.LocationResponse
	err = json.Unmarshal(res.Body(), &locations)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling locations: %v", err)
	}

	// Cache the response body for future use
	err = g.cache.Set(ctx, "@chatbot:guanabara:locations", res.Body(), 24*time.Hour).Err()
	if err != nil {
		return nil, nil, fmt.Errorf("error caching locations: %v", err)
	}

	return locations.Data.Locations, locations.Data.Groups, nil
}

func (g *guanabaraClient) GetCity(ctx context.Context, city string) ([]string, error) {
	var response []string

	locationsArray, groupsArray, err := g.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(city)

	for _, location := range locationsArray {
		if strings.Contains(utils.RemoveDiacritics(strings.ToLower(location.CityName)), utils.RemoveDiacritics(strings.ToLower(city))) {
			fmt.Println(location)
			format := fmt.Sprintf(`%d - %s`, location.ID, location.Name)

			response = append(response, format)
		}

	}

	for _, location := range groupsArray {
		if strings.Contains(utils.RemoveDiacritics(strings.ToLower(location.Name)), utils.RemoveDiacritics(strings.ToLower(city))) {
			fmt.Println(location)
			format := fmt.Sprintf(`%d:group - %s`, location.Id, location.Name)

			response = append(response, format)
		}

	}

	fmt.Println("Pontos:", response)

	return response, nil
}

func (g *guanabaraClient) GetSchedules(ctx context.Context, departure, arrival, dataDeparture string) ([]models.TripInfoResponse, error) {
	fmt.Print("Entro no g.GetSchedules()")
	var numDepartureGroup int
	var numDeparture int
	var numDestinationGroup int
	var numDestination int
	var err error
	var getTripsRequest models.GetTrips
	var response []string

	if strings.Contains(departure, ":group") {
		numDepartureGroup, err = strconv.Atoi(strings.Split(departure, ":")[0])
		if err != nil {
			return nil, err
		}
	} else {
		numDeparture, err = strconv.Atoi(departure)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(arrival, ":group") {
		numDestinationGroup, err = strconv.Atoi(strings.Split(arrival, ":")[0])
		if err != nil {
			return nil, err
		}
	} else {
		numDestination, err = strconv.Atoi(arrival)
		if err != nil {
			return nil, err
		}
	}

	fmt.Print(numDeparture)
	fmt.Print(numDestination)
	fmt.Print(dataDeparture)
	getTripsRequest.DepartureID = numDeparture
	getTripsRequest.DestinationID = numDestination
	getTripsRequest.DepartureDate = dataDeparture
	getTripsRequest.DepartureGroupID = numDepartureGroup
	getTripsRequest.DestinationGroupID = numDestinationGroup

	trips, err := g.GetTrips(ctx, getTripsRequest)
	if err != nil {
		return nil, err
	}
	if len(trips) == 0 {
		return nil, nil
	}

	for _, trip := range trips {
		fmt.Println(trip)
		//format := fmt.Sprintf(`Horario de saida: %s - Empresa: %s - Valor: %f`, trip.DepartureTime, trip.CompanyName, trip.PriceValue)
		format := fmt.Sprintf(". %s \n 📍 %s - %s \n 🏁 %s - %s \n %s", trip.CompanyName, trip.DepartureTime, trip.DepartureLocation.Name, trip.ArrivalTime, trip.ArrivalLocation.Name, trip.ClassOfServiceName)
		response = append(response, format)
	}

	return trips, nil
}

func (g *guanabaraClient) RouteValidation(ctx context.Context, departure, arrival string) (bool, error) {
	var response bool

	fmt.Println(departure)
	fmt.Println(arrival)

	numDeparture, err := strconv.Atoi(departure)
	numArrival, err := strconv.Atoi(arrival)
	if err != nil {
		return false, err
	}

	locationsArray, err := g.GetLocationsPair(ctx)
	if err != nil {
		return false, err
	}

	for _, location := range locationsArray {
		if location.OriginID == numDeparture && location.DestinationID == numArrival {
			response = true
		}
	}

	return response, nil
}
