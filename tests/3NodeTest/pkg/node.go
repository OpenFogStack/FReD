package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Node struct {
	url string
}

func NewNode(url string) (node *Node) {
	node = &Node{url: url}
	return
}

func (n *Node) CreateKeygroup(kgname string, expectedStatusCode int) (responseBody map[string]string) {
	log.Debug().Msgf("Sending a Create Keygroup for group %s; expecting %d", kgname, expectedStatusCode)
	return n.sendPost("keygroup/"+kgname, nil, expectedStatusCode)
}

func (n *Node) DeleteKeygroup(kgname string, expectedStatusCode int) (responseBody map[string]string) {
	log.Debug().Msgf("Sending a Delete Keygroup for group %s; expecting %d", kgname, expectedStatusCode)
	return n.sendDelete("keygroup/"+kgname, nil, expectedStatusCode)
}

func (n *Node) PutItem(kgname, item string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	log.Debug().Msgf("Sending a Put for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)
	return n.sendPut(fmt.Sprintf("/keygroup/%s/item/%s", kgname, item), data, expectedStatusCode)
}

func (n *Node) sendGet(path string, expectedStatusCode int) (responseBody map[string]string) {
	resp, err := http.Get(n.url + path)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal().Err(err).Msg("SendGet got HTTP error")
		return nil
	}
	if resp.StatusCode != expectedStatusCode {
		log.Error().Msgf("SendGet got wrong HTTP Status Code Response. Expected: %d, Got: %d", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), responseBody)

	if err != nil {
		log.Fatal().Msg("sendGet got a response with invalid json")
	}

	return
}

func (n *Node) sendPut(path string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	client := &http.Client{}
	jsonBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPut, n.url+path, bytes.NewBuffer(jsonBytes))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal().Err(err).Msg("sendPut got HTTP error")
		return nil
	}
	if resp.StatusCode != expectedStatusCode {
		log.Error().Msgf("SendPut got wrong HTTP Status Code Response. Expected: %d, Got: %d", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), responseBody)
	if err != nil {
		log.Fatal().Msg("sendPut got a response with invalid json")
	}
	return
}

func (n *Node) sendPost(path string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	var jsonBytes []byte
	if data != nil {
		var err error
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			log.Fatal().Msgf("Cannot marshal JSON: %v", data)
		}
	}
	var resp *http.Response
	var err error
	if jsonBytes != nil {
		resp, err = http.Post(n.url+path, "application/json", bytes.NewBuffer(jsonBytes))
	} else {
		resp, err = http.Post(n.url+path, "", nil)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("sendPost got HTTP error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatusCode {
		log.Error().Msgf("SendPost got wrong HTTP Status Code Response. Expected: %d, Got: %d.", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), responseBody)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		log.Fatal().Msg("sendGet got a response with invalid json")
	}
	return
}

func (n *Node) sendDelete(path string, data map[string]string, expectedStatusCode int) (responseBody map[string]string) {
	var jsonBytes []byte
	if data != nil {
		var err error
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			log.Fatal().Msgf("Cannot marshal JSON: %v", data)
		}
	}
	var resp *http.Response
	var err error
	client := &http.Client{}
	if jsonBytes != nil {
		jsonBytes, _ := json.Marshal(data)
		req, _ := http.NewRequest(http.MethodDelete, n.url+path, bytes.NewBuffer(jsonBytes))
		resp, err = client.Do(req)
	} else {
		req, _ := http.NewRequest(http.MethodDelete, n.url+path, nil)
		resp, err = client.Do(req)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("sendPost got HTTP error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatusCode {
		log.Error().Msgf("SendDelete got wrong HTTP Status Code Response. Expected: %d, Got: %d.", expectedStatusCode, resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// Load buf into responseBody
	err = json.Unmarshal(buf.Bytes(), responseBody)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		log.Fatal().Msg("sendGet got a response with invalid json")
	}
	return
}
