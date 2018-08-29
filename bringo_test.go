package bringo

import (
	"testing"
	"time"
	"os"
)

func initValid(t *testing.T) *Bringo {

	login :=  os.Getenv("BRINGO_LOGIN")
	password :=  os.Getenv("BRINGO_PASSWORD")

	api := &Bringo{}
	api.Init(login, password, false)
	if err := api.Login(); err == nil {
		return api
	} else {
		t.Fatal("LogIn correct fail")
	}
	return nil
}

func getDelivery() *Delivery {
	ct := time.Now()
	return &Delivery{
		Name:        "Some kind of test delivery",   //название доставки
		Description: "This is demo delivery for AV", // описание, курьер увидит его в комментарии к доставке,
		ExternalID:  "string up to 128 characters",
		DeliverySegments: []DeliverySegment{ // сегменты. Массив. Пока из одного элемента.
			{
				CargoCost: 100.00, // стоимость груза
				Height: 1.00, // габариты, в метрах
				Length: 1.00,
				Width: 1.00,
				Weight: 1.00, // вес, в килограммах
				IsBuyout: false, // true, если курьер должен выкупить груз, а потом получить деньги от покупателя
				From: Destination{
					Address: Address{
						GeoPoint: GeoPoint{ //координаты
							Lat: 55.738784790039062,
							Lng: 37.548374176025391,
						},
						AddressText: "Студенческая, Москва", //текстовый адрес
						MetroName: "Студенческая", // название метро
						Contact: "Jane", // имя отправителя
						Phone: "+79525741559", // телефон отправителя
						Comment: "Some comment for address", // комментарий к адресу
						CityID: 1, //id города. 1 - Москва, 2 - Питер
					},
					TimeInterval: TimeInterval{
						From: time.Date(2018,06,04,23,23,0,0,ct.Location()), // "2018-06-04T23:23:00.7133073+03:00" начало интервала получения груза
						To: time.Date(2018,06,04,23,53,0,0,ct.Location()), // "2018-06-04T23:53:00.7143079+03:00" конец интервала получения груза
					},
				},
				To: Destination{ // куда
					Address: Address{
						GeoPoint: GeoPoint{
							Lat: 55.729648590087891,
							Lng: 37.470893859863281,
						},
						AddressText: "Славянский бульвар, Москва",
						MetroName: "Славянский бульвар",
						Contact: "John",
						Phone: "+79525741559",
						Comment: "",
						CityID: 1,
					},
					TimeInterval: TimeInterval{
						From: time.Date(2018,06,04,23,53,0,0,ct.Location()), // "2018-06-04T23:53:00.7143079+03:00",
						To: time.Date(2018,06,5,0,53,0,0,ct.Location()), //"2018-06-05T00:53:00.7143079+03:00"
					},
				},
			},
		},
	}
}

func TestBringo_Init(t *testing.T) {
	api := &Bringo{}
	api.Init("test", "test", false)
}

func TestBringo_Login(t *testing.T) {
	initValid(t)
}

func TestBringo_Calculate(t *testing.T) {
	api := initValid(t)

	price, err := api.Calculate(getDelivery())

	if err != nil {
		t.Fatal("Calculation error", err)
	}

	if price <= 0 {
		t.Fatal("Calculation result error price", price)
	}
}

