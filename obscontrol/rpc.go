package obscontrol

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

func (o *OBSClient) sendRequest(requestType string, requestId string, requestData map[string]interface{}) error {
	payload := map[string]interface{}{
		"op": 6,
		"d": map[string]interface{}{
			"requestType": requestType,
			"requestId":   requestId,
			"requestData": requestData,
		},
	}
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return o.Conn.WriteMessage(websocket.TextMessage, msg)
}

func (o *OBSClient) SwitchScene(sceneName string) error {
	o.idCounter++
	return o.sendRequest("SetCurrentProgramScene", "switchScene_"+strconv.Itoa(o.idCounter), map[string]interface{}{
		"sceneName": sceneName,
	})
}

func (o *OBSClient) ToggleSource(name string, visible bool) error {
	itemId, ok := o.SceneItems[name]
	if !ok {
		return ErrNotFound
	}
	o.idCounter++
	return o.sendRequest("SetSceneItemEnabled", "toggle_"+strconv.Itoa(o.idCounter), map[string]interface{}{
		"sceneName":       o.SceneName,
		"sceneItemId":     itemId,
		"sceneItemEnabled": visible,
	})
}
