package exception

import (
	"encoding/json"
	"fmt"
)

type UnprocessableEntityError struct {
	Msg map[string]string `json:"msg"`
}

func (e *UnprocessableEntityError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func (e *UnprocessableEntityError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]string{"msg": e.Msg})
}
