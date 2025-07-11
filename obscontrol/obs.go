package obscontrol 

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type OBSClient struct {
	Conn      *websocket.Conn
	SceneName string
	SceneItems map[string]int
	password  string
	idCounter int
}

func NewOBSClient(password string) *OBSClient {
	return &OBSClient{
		password:  password,
		SceneItems: make(map[string]int),
	}
}

func (o *OBSClient) Connect(url string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	o.Conn = conn

	// Step 1: wait for Hello
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	var hello map[string]interface{}
	json.Unmarshal(msg, &hello)
	authData := hello["d"].(map[string]interface{})["authentication"].(map[string]interface{})
	challenge := authData["challenge"].(string)
	salt := authData["salt"].(string)
	rpcVersion := int(hello["d"].(map[string]interface{})["rpcVersion"].(float64))

	// Step 2: send Identify
	auth := computeAuth(o.password, salt, challenge)
	identify := map[string]interface{}{
		"op": 1,
		"d": map[string]interface{}{
			"rpcVersion":   rpcVersion,
			"authentication": auth,
		},
	}
	payload, _ := json.Marshal(identify)
	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}

	// Step 3: wait for Identified (op 2)
	_, msg, err = conn.ReadMessage()
	if err != nil {
		return err
	}
	log.Println("OBS identified successfully")

	return nil
}

func (o *OBSClient) Close() {
	o.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(time.Second)
	o.Conn.Close()
}

func computeAuth(password, salt, challenge string) string {
	firstHash := sha256.Sum256([]byte(password + salt))
	secret := base64.StdEncoding.EncodeToString(firstHash[:])
	finalHash := sha256.Sum256([]byte(secret + challenge))
	return base64.StdEncoding.EncodeToString(finalHash[:])
}
