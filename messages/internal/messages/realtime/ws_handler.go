package realtime

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SimpleAuthClient — очень лёгкий клиент для валидации JWT через auth service
type SimpleAuthClient struct {
	AuthURL string
}

func NewSimpleAuthClient(authURL string) *SimpleAuthClient {
	return &SimpleAuthClient{AuthURL: authURL}
}

// ValidateToken вызывает auth service и возвращает userID (string) или ошибку
func (a *SimpleAuthClient) ValidateToken(token string) (string, error) {
	req, err := http.NewRequest("GET", a.AuthURL+"/api/v1/auth/validate", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth validate failed: %s", string(body))
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.UserID, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWS — HTTP handler для /ws
// query param: token=<jwt>
func ServeWS(hub *Hub, authClient *SimpleAuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		userIDStr, err := authClient.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// валидируем uuid
		if _, err := uuid.Parse(userIDStr); err != nil {
			http.Error(w, "invalid user id from auth", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		client := &Client{
			UserID: userIDStr,
			Conn:   conn,
			send:   make(chan []byte, 256),
			Hub:    hub,
		}

		hub.register <- client
		go client.WriteLoop()
		go client.ReadLoop()

		welcome := SendMessageToClient{
			Type:   "welcome",
			ChatID: "",
			Data:   map[string]string{"user_id": userIDStr},
		}
		b, _ := json.Marshal(welcome)
		client.send <- b
	}
}
