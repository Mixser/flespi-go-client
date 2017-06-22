package flesapi

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
)

const baseApiUrl = "https://flespi.io"

type client struct {
	Token string
}

type Result2 interface {
}

type Result struct {
	Ident 	  string `json:"ident"`
	Timestamp int64 `json:"timestamp"`

	Longitude float64 `json:"position.longitude"`
	Latitude  float64 `json:"position.latitude"`

	Speed     float64 `json:"position.speed"`
	Course    float64 `json:"position.direction"`
}


func (r Result) ToString() string{
	return fmt.Sprintf("<ident: %s pos: (%0.6f, %0.6f, %0.2f, %0.2f) tm: %d>", r.Ident,
		r.Longitude,
		r.Latitude,
		r.Speed,
		r.Course,
		r.Timestamp)
}

func (r Result) String() string {
	return r.ToString()
}

type Error struct{
	Id int `json:"id"`
	Code int `json:"code"`
	Reason string `json:"reason"`
}

type Response struct {
	Result []Result `json:"result"`
	Errors []Error `json:"errors"`

	Next_key int `json:"next_key"`
}


type Args interface {
	encodeArgs() ([]byte, error)
}


type MessageArgs struct {
	Curr_key int `json:"curr_key"`
	Limit_count int `json:"limit_count"`
	//Limit_size int `json:"limit_size"`
	Timeout int `json:"timeout"`
	Delete bool `json:"delete"`
}


func (m MessageArgs) encodeArgs() ([]byte, error) {
	return json.Marshal(m)
}


func (c client) doRequest(url string, args Args) (*Response, error) {
	httpClient := &http.Client{}

	log.Println(url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Request create error %v", err)
		return nil, err
	}

	request.Header.Add("Authorization", "FlespiToken " + c.Token)
	httpResponse, err := httpClient.Do(request)

	if err != nil {
		log.Fatalf("Fetch error %v", err)
		return nil, err
	}

	defer httpResponse.Body.Close()

	bytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		log.Fatalf("Read error %v", err)
		return nil, err
	}

	response := Response{}
	err = json.Unmarshal(bytes, &response)

	if err != nil {
		log.Fatalf("Unmarshal error %v", err)
		return nil, err
	}

	return &response, nil
}


func (c client) GetChannelMessages(channel int, args MessageArgs) (*Response, error) {
	data, err := args.encodeArgs()
	if err != nil {
		log.Fatalf("Encode args error %v", err)
		return nil, err
	}
	url := fmt.Sprintf("%s/gw/channels/%d/messages?data=%s", baseApiUrl, channel, data)
	return c.doRequest(url, args)
}

func (c client) GetAbqueMessages(abque int, args MessageArgs) (*Response, error) {
	data, err := args.encodeArgs()

	if err != nil {
		log.Fatalf("Encode args error %v", err)
	}

	url := fmt.Sprintf("%s/abques/%d/messages?data=%s", baseApiUrl, abque, data)
	return c.doRequest(url, args)
}

func NewClient(token string) *client {
	apiClient := new(client)
	apiClient.Token = token
	return apiClient
}
