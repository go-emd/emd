package leader

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// The type that all REST endpoint requests get 
// transferred to.  This then get serialized to 
// json very easily.
type Response map[string]interface{}

// Allows the leader.Response type to be 
// easily serialized to json when sent back 
// to the clients REST request.
func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}

	s = string(b)
	return
}

// Sets the REST headers for the json response and sends the 
// json serialized data with it.
func Respond(rw http.ResponseWriter, success bool, message interface{}) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{"success": success, "message": message})
}
