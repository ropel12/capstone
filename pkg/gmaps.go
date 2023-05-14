package pkg

import (
	"context"

	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

type Client struct {
	Client *maps.Client
	Log    *logrus.Logger
}
type Result struct {
	SchoolName string `json:"school_name"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Village    string `json:"village"`
	Detail     string `json:"detail"`
}

func NewClientGmaps(apikey string, log *logrus.Logger) *Client {
	client, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Errorf("[ERROR]When creating a Google Maps client object, Err: %v", err)
	}
	return &Client{
		Client: client,
		Log:    log,
	}
}

func (c *Client) Search(searchval string) any {
	r := &maps.PlaceAutocompleteRequest{
		Input:    searchval,
		Language: "id",
	}

	resp, err := c.Client.PlaceAutocomplete(context.Background(), r)
	if err != nil {
		c.Log.Errorf("[ERROR]WHEN GETTING DATA FROM GOOGLEMAPS, Err: %v", err)
	}
	res := []Result{}
	if len(resp.Predictions) == 0 {
		return res
	}
	for i, val := range resp.Predictions {
		if i > 4 {
			break
		}
		data := Result{}
		data.Detail = val.Description
		if len(val.Terms) == 6 {
			data.SchoolName = val.Terms[0].Value
			data.Village = val.Terms[2].Value
			data.District = val.Terms[2].Value
			data.City = val.Terms[3].Value
			data.Province = val.Terms[4].Value
			data.SchoolName = val.Terms[0].Value
		} else {
			data.SchoolName = val.Terms[0].Value
			data.Village = val.Terms[1].Value
			data.District = val.Terms[1].Value
			data.City = val.Terms[2].Value
			data.Province = val.Terms[3].Value
			data.SchoolName = val.Terms[0].Value
		}
		res = append(res, data)

	}
	return res
}
