package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gustavodamazio/go-test/models"
	"github.com/gustavodamazio/go-test/services/callback"
	"github.com/gustavodamazio/go-test/services/storage"
)

func HandleHtmlToPDF(storageService *storage.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleHtmlToPDF(w, r, storageService)
	}
}

func handleHtmlToPDF(w http.ResponseWriter, r *http.Request, storageService *storage.StorageService) {
	// Read JSON string from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON to extract HTML string
	var requestBody models.RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	html_storage_path := requestBody.Data.HTML_STORAGE_PATH
	if html_storage_path == "" {
		http.Error(w, "HTML storage path is empty", http.StatusBadRequest)
		return
	}

	html_file, err := storageService.ReadFile(html_storage_path)
	if err != nil {
		http.Error(w, "Failed to read HTML file", http.StatusInternalServerError)
		return
	}

	// Create a new PDF generator instance
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		http.Error(w, "Failed to create PDF generator", http.StatusInternalServerError)
		return
	}

	// Add HTML content to the PDF generator
	pdfg.AddPage(pdf.NewPageReader(bytes.NewReader(html_file)))
	pdfg.Orientation.Set(pdf.OrientationLandscape)

	// Generate the PDF
	err = pdfg.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	html_storage_path_output := html_storage_path + "-output.pdf"
	err = storageService.WriteFile(html_storage_path_output, pdfg.Bytes())
	if err != nil {
		http.Error(w, "Failed to write PDF file", http.StatusInternalServerError)
		return
	}

	// Send callback
	if requestBody.Data.CALLBACK_URL != "" {
		callbackService := callback.NewCallbackService()
		err = callbackService.SendCallback(requestBody)
		if err != nil {
			fmt.Printf("Failed to send callback: %v\n", err)
			return
		}
	}

	// Set response headers
	response := fmt.Sprintf(`{"data":{"message":"PDF generated successfully","pdf_storage_path": "%s"}}`, html_storage_path_output)
	w.Header().Set("Content-Type", "application/json")
	// Write reponse to the client with a success message and output file path
	_, err = w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Failed to send PDF", http.StatusInternalServerError)
		return
	}
}
