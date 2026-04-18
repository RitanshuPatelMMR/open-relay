package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type wsMessage struct {
	Type      string          `json:"type"`
	EventID   string          `json:"event_id"`
	Status    string          `json:"status"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
	Headers   json.RawMessage `json:"headers"`
	Method    string          `json:"method"`
	Path      string          `json:"path"`
}

var localPort int

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Forward incoming webhooks to localhost",
	Long:  `Connects to OpenRelay via WebSocket and forwards every incoming webhook to your local server.`,
	RunE:  runListen,
}

func init() {
	listenCmd.Flags().IntVarP(&localPort, "port", "p", 3000, "Local port to forward webhooks to")
	rootCmd.AddCommand(listenCmd)
}

func runListen(cmd *cobra.Command, args []string) error {
	if apiKey == "" {
		return fmt.Errorf("--api-key is required")
	}

	// build websocket URL
	wsBase := strings.Replace(serverURL, "http://", "ws://", 1)
	wsBase = strings.Replace(wsBase, "https://", "wss://", 1)
	wsURL := fmt.Sprintf("%s/ws?api_key=%s", wsBase, url.QueryEscape(apiKey))

	fmt.Printf("» Connecting to %s\n", serverURL)
	fmt.Printf("» Forwarding webhooks → http://localhost:%d\n", localPort)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket connect failed: %w", err)
	}
	defer conn.Close()

	fmt.Println("» Connected. Waiting for events...\n")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	msgCh := make(chan wsMessage)
	errCh := make(chan error)

	go func() {
		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			var msg wsMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}
			msgCh <- msg
		}
	}()

	for {
		select {
		case <-quit:
			fmt.Println("\n» Disconnected.")
			return nil

		case err := <-errCh:
			return fmt.Errorf("websocket error: %w", err)

	case msg := <-msgCh:
        if msg.EventID == "" {
        continue
    }
    go forwardToLocal(msg, localPort)
		}
	}
}

func forwardToLocal(msg wsMessage, port int) {
	targetURL := fmt.Sprintf("http://localhost:%d%s", port, msg.Path)
	if msg.Path == "" {
		targetURL = fmt.Sprintf("http://localhost:%d/", port)
	}

	method := msg.Method
	if method == "" {
		method = "POST"
	}

	var bodyReader *strings.Reader
	if msg.Payload != nil {
		bodyReader = strings.NewReader(string(msg.Payload))
	} else {
		bodyReader = strings.NewReader("{}")
	}

	req, err := http.NewRequest(method, targetURL, bodyReader)
	if err != nil {
		fmt.Printf("✗ [%s] build request failed: %v\n", msg.EventID[:8], err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OpenRelay-Event-ID", msg.EventID)

	// forward original headers if present
	if msg.Headers != nil {
		var headers map[string]interface{}
		if json.Unmarshal(msg.Headers, &headers) == nil {
			for k, v := range headers {
				if s, ok := v.(string); ok {
					req.Header.Set(k, s)
				}
			}
		}
	}

	start := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		fmt.Printf("✗ [%s] %s %s — error: %v\n", msg.EventID[:8], method, msg.Path, err)
		return
	}
	defer res.Body.Close()

	fmt.Printf("✓ [%s] %s %s → %d  %dms\n", msg.EventID[:8], method, targetURL, res.StatusCode, elapsed)
}