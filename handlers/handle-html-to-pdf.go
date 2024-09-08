package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gustavodamazio/go-test/models"
)

func HandleHtmlToPDF(w http.ResponseWriter, r *http.Request) {
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

	html := requestBody.Data.HTML
	if html == "" {
		http.Error(w, "HTML content is empty", http.StatusBadRequest)
		return
	}

	// Create a new PDF generator instance
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		http.Error(w, "Failed to create PDF generator", http.StatusInternalServerError)
		return
	}

	// Add HTML content to the PDF generator
	pdfg.AddPage(pdf.NewPageReader(strings.NewReader(html)))
	pdfg.Orientation.Set(pdf.OrientationLandscape)

	// Generate the PDF
	err = pdfg.Create()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfg.Bytes())))

	// Write PDF bytes to the response
	_, err = w.Write(pdfg.Bytes())
	if err != nil {
		http.Error(w, "Failed to send PDF", http.StatusInternalServerError)
		return
	}
}
