package node

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	client "gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/vendor/go-client"
)

func getRecvBody(recvErr error, recv *http.Response, recvData interface{}) string {
	// check if we got a Data object that contains our data (if we don't have it already)
	if recvData != nil {

		// i honestly have no idea how this works i just copied it from here
		// https://stackoverflow.com/questions/42511940/golang-deepequal-and-reflect-zero
		field := reflect.ValueOf(recvData).Field(0)

		// but first let's check if we didn't actually receive an empty struct
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return fmt.Sprintf("%#v", recvData)
		}
	}

	// if we received an error that is a *client.GenericSwaggerError, that error has a Body attribute that has the response body
	if recvErr != nil {
		if recvErr, ok := recvErr.(client.GenericSwaggerError); ok {
			m := recvErr.Model()
			if m, ok := m.(client.ModelError); ok {
				return m.Error_
			}

			return string(recvErr.Body())
		}
	}

	// check if the data is still in the buffer from the response (if we don't have it already)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(recv.Body); err != nil {
		return buf.String()
	}

	return ""
}

func checkResponse(recv *http.Response, recvErr error, recvData interface{}, expCode int, expEmpty bool) (err error) {
	// if there is no response, something went wrong with the request
	if recv == nil {
		return fmt.Errorf("response failed")
	}

	recvCode := recv.StatusCode

	// check if the status code is not what we expected
	if recvCode != expCode {
		return fmt.Errorf("expected response code %#v but got %#v", expCode, recvCode)
	}

	recvBody := getRecvBody(recvErr, recv, recvData)

	// if we expect an empty response but we got something else, report that
	if expEmpty && (recvBody != "") {
		return fmt.Errorf("expected an empty response but got %#v", recvBody)
	}

	// if we didn't expect an empty response but we did get one, report that as well
	if !expEmpty && (recvBody == "" || len(recvBody) == 0) {
		return fmt.Errorf("expected a response but got nothing")
	}

	return nil

}
