package callback

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gustavodamazio/go-test/models"
)

// CallbackService é a estrutura do serviço
type CallbackService struct {
	Client *http.Client
}

// NewCallbackService cria uma nova instância do CallbackService
func NewCallbackService() *CallbackService {
	return &CallbackService{
		Client: &http.Client{},
	}
}

// SendCallback executa a requisição de callback
func (cs *CallbackService) SendCallback(requestBody models.RequestBody) error {
	// Converte CALLBACK_DATA para um byte slice
	data := []byte(requestBody.Data.CALLBACK_DATA)

	// Cria uma nova requisição HTTP
	req, err := http.NewRequest(requestBody.Data.CALLBACK_METHOD, requestBody.Data.CALLBACK_URL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("erro ao criar a requisição: %w", err)
	}

	// Define os headers da requisição
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", requestBody.Data.CALLBACK_HEADERS.Authorization)

	// Envia a requisição HTTP
	resp, err := cs.Client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar a requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler a resposta: %w", err)
	}

	// Exibe a resposta
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", body)

	return nil
}
