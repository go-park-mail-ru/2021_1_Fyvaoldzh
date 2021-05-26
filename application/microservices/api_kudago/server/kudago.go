package kudago_server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"kudago/application/microservices/api_kudago/kudago"
	"kudago/application/microservices/api_kudago/kudago_proto"
	"kudago/application/models"
	"kudago/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type KudagoServer struct {
	usecase kudago.Usecase
	logger  *logger.Logger
}

func NewKudagoServer(usecase kudago.Usecase, lg *logger.Logger) *KudagoServer {
	return &KudagoServer{usecase: usecase, logger: lg}
}

func (k *KudagoServer) AddBasic(_ context.Context, input *kudago_proto.Input) (*kudago_proto.Empty, error) {

	str := "https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,place,body_text,categories,tags,images&order_by=-publication_date&categories=cinema,education,entertainment,exhibition,festival,tour,concert"
	for i := uint64(0); i < input.Num; {
		answer, err := k.GetEvents(str)
		if err != nil {
			return nil, err
		}
		str = answer.Next
		for _, elem := range answer.Results {

			if elem.Place.Id != 0 && len(elem.Tags) > 0 && len(elem.Title) <= 60 {
				url := "https://kudago.com/public-api/v1.4/places/" + strconv.FormatUint(elem.Place.Id, 10) + "/?fields=title,coords,subway,address"
				place, err := k.GetPlace(url)
				if err != nil {
					return nil, err
				}
				if len(place.Title) <= 60 && len(place.Subway) <= 60 && len(place.Subway) <= 60 && place.Map.Longitude != 0 {
					flag, err := k.usecase.AddEvent(elem, place)
					if err != nil {
						return nil, err
					}
					if flag {
						i++
						if i == input.Num {
							return &kudago_proto.Empty{}, nil
						}
					}
				}

			}
			if answer.Next == "" {
				return &kudago_proto.Empty{}, nil
			}
		}
	}

	return &kudago_proto.Empty{}, nil
}

func (k *KudagoServer) AddToday(_ context.Context, _ *kudago_proto.Empty) (*kudago_proto.Empty, error) {
	k.logger.Debug("start pushing today's events")
	str := "https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,place,body_text,categories,tags,images&order_by=-publication_date&categories=cinema,education,entertainment,exhibition,festival,tour,concert"
	timeNow := time.Now().Unix()
	for {
		answer, err := k.GetEvents(str)
		if err != nil {
			return nil, err
		}
		str = answer.Next
		for _, elem := range answer.Results {
			if elem.PublicationDate < timeNow-86400 {
				return &kudago_proto.Empty{}, nil
			}
			if elem.Place.Id != 0 && len(elem.Tags) > 0 && len(elem.Title) <= 60 {
				url := "https://kudago.com/public-api/v1.4/places/" + strconv.FormatUint(elem.Place.Id, 10) + "/?fields=title,coords,subway,address"
				place, err := k.GetPlace(url)
				if err != nil {
					return nil, err
				}
				if len(place.Title) <= 60 && len(place.Subway) <= 60 && len(place.Subway) <= 60 && place.Map.Longitude != 0 {
					flag, err := k.usecase.AddEvent(elem, place)
					if err != nil {
						return nil, err
					}
					if flag {
						k.logger.Debug("successfully added 1 event")
					}
				}

			}
			if answer.Next == "" {
				return &kudago_proto.Empty{}, nil
			}
		}
	}

}

func (k *KudagoServer) GetPlace(url string) (models.Place, error) {
	placeResp, err := http.Get(url)
	if err != nil {
		k.logger.Warn(err)
		return models.Place{}, err
	}
	placeBody, err := ioutil.ReadAll(placeResp.Body)
	if err != nil {
		k.logger.Warn(err)
		return models.Place{}, err
	}
	var place models.Place
	err = json.Unmarshal(placeBody, &place)
	if err != nil {
		k.logger.Warn(err)
		return models.Place{}, err
	}

	placeResp.Body.Close()

	return place, nil
}

func (k *KudagoServer) GetEvents(url string) (models.Answer, error) {
	resp, err := http.Get(url)
	if err != nil {
		k.logger.Warn(err)
		return models.Answer{}, err
	}

	var answer models.Answer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		k.logger.Warn(err)
		return models.Answer{}, err
	}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		k.logger.Warn(err)
		return models.Answer{}, err
	}

	resp.Body.Close()

	return answer, nil
}
