package bringo

import (
	"github.com/go-resty/resty"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
)

const Url = "https://api.bringo247.ru/api"
const UrlTest = "https://lugh-demo.bringo247.ru/api"

type Bringo struct {
	url         string
	login       string
	password    string
	client      *resty.Client
	AccountInfo *LoginResponse
}

func (s *Bringo) Init(login string, password string, production bool) {
	s.login = login
	s.password = password
	if s.url = Url; !production {
		s.url = UrlTest
	}
	s.client = resty.New()
}

func (s *Bringo) request(method string, body interface{}, isPost bool) ([]byte, error) {
	request := s.client.R().SetHeader("Accept", "application/json")

	var response *resty.Response
	var err error

	if isPost {
		response, err = request.SetBody(body).Post(s.url + "/" + method)
	} else {
		response, err = request.Get(s.url + "/" + method)
	}

	if err != nil {
		log.Fatal("[Error] Bringo::PostMethod:", err.Error(), "Response:", string(response.Body()))
		return nil, err
	}

	var errResp ErrorResponse

	if err := json.Unmarshal(response.Body(), &errResp); err == nil {
		if errResp.Error.Message != "" {
			return nil, fmt.Errorf(errResp.Error.Message+" Code: %d", errResp.Error.Code)
		}
	}

	return response.Body(), nil
}

func (s *Bringo) Get(method string) ([]byte, error) {
	return s.request(method, nil, false)
}

func (s *Bringo) GetUnmarshal(method string, unmarshal interface{}) (error) {

	response, err := s.Get(method)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(response, unmarshal); err != nil {
		return err
	}

	return nil
}

func (s *Bringo) Post(method string, body interface{}) ([]byte, error) {
	return s.request(method, body, true)
}

func (s *Bringo) PostUnmarshal(method string, body interface{}, unmarshal interface{}) (error) {

	response, err := s.Post(method, body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(response, unmarshal); err != nil {
		return err
	}

	return nil
}

func (s *Bringo) Login() error {
	response, err := s.Post("login", &LoginRequest{Login: s.login, Password: s.password})
	if err != nil {
		return err
	}

	if err := json.Unmarshal(response, &s.AccountInfo); err != nil {
		return err
	}

	return nil
}


func (s *Bringo) Calculate(delivery *Delivery) (float64, error) {
	var result CalculateResponse

	err := s.PostUnmarshal("deliveries/price", delivery, &result)
	if err != nil {
		return -1, err
	}

	return result.Result, nil
}

func (s *Bringo) Create(delivery *Delivery) (*InfoResponse, error) {
	var result InfoResponse

	err := s.PostUnmarshal("deliveries/create", delivery, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Bringo) Cancel(id int) (*InfoResponse, error) {
	var result InfoResponse

	err := s.GetUnmarshal("deliveries/cancel/" + strconv.Itoa(id), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Bringo) Info(id int) (*InfoResponse, error) {
	var result InfoResponse

	err := s.GetUnmarshal("deliveries/info/" + strconv.Itoa(id), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func New() *Bringo  {
	instance := &Bringo{}
	return instance
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Result struct {
		SensitiveInfo struct {
			EMail            string `json:"eMail"`
			EMailConfirmed   bool   `json:"eMailConfirmed"`
			PhoneConfirmed   bool   `json:"phoneConfirmed"`
			RegistrationDate string `json:"registrationDate"`
		} `json:"sensitiveInfo"`
		Roles []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			Code  string `json:"code"`
		} `json:"roles"`
		ID            int    `json:"id"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		State         int    `json:"state"`
		HasCar        bool   `json:"hasCar"`
		IsDeaf        bool   `json:"isDeaf"`
		ProofOfAge    bool   `json:"proofOfAge"`
		OfferAccepted bool   `json:"offerAccepted"`
	} `json:"result"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

type InfoResponse struct {
	Result struct {
		ID         int            `json:"id"`
		Price      float64        `json:"price"`
		Deliveries []DeliveryInfo `json:"deliveries"`
	} `json:"result"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	GeoPoint    GeoPoint `json:"geoPoint"`
	AddressText string   `json:"addressText"`
	MetroName   string   `json:"metroName"`
	Contact     string   `json:"contact"`
	Phone       string   `json:"phone"`
	Comment     string   `json:"comment"`
	CityID      int      `json:"cityId"`
}

type TimeInterval struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Destination struct {
	Address      Address      `json:"address"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

type DeliverySegment struct {
	From      Destination `json:"from"`
	To        Destination `json:"to"`
	CargoCost int         `json:"cargoCost"`
	Height    int         `json:"height"`
	Length    int         `json:"length"`
	Width     int         `json:"width"`
	Weight    int         `json:"weight"`
	IsBuyout  bool        `json:"isBuyout"`
}

type Delivery struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	ExternalID       string            `json:"externalId"`
	DeliverySegments []DeliverySegment `json:"deliverySegments"`
}

type DeliveryInfo struct {
	ID           int     `json:"id"`
	Price        float64 `json:"price"`
	Name         string  `json:"name"`
	OldCloseCode string  `json:"oldCloseCode"`
	OldState     int     `json:"oldState"`
}
