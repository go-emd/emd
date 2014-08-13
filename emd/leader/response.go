package leader

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}

	s = string(b)
	return
}

func Respond(rw http.ResponseWriter, success bool, message interface{}) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{"success": success, "message": message})
}
