package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2" // Framework web Fiber
	"gonum.org/v1/gonum/mat"      // Biblioteca para álgebra lineal
)

// MatrixRequest representa la estructura JSON de entrada
type MatrixRequest struct {
	Data [][]float64 `json:"data"`
}

// Estructura de la respuesta: matrices Q y R
type QRResponse struct {
	Q [][]float64 `json:"Q"`
	R [][]float64 `json:"R"`
}

// rotaterMatriz rota la matriz original 90° en sentido antihorario
func rotaterMatriz(matriz [][]float64) [][]float64 {
	filas := len(matriz)
	columnas := len(matriz[0])

	rotada := make([][]float64, columnas)
	for i := range rotada {
		rotada[i] = make([]float64, filas)
		for j := range rotada[i] {
			rotada[i][j] = matriz[filas-1-j][i]
		}
	}
	return rotada
}

// convertirAMatriz convierte una matriz de Gonum a una matriz bidimensional de float64
func convertirAMatriz(m mat.Matrix) [][]float64 {
	r, c := m.Dims()
	resultado := make([][]float64, r)
	for i := range resultado {
		resultado[i] = make([]float64, c)
		for j := range resultado[i] {
			resultado[i][j] = m.At(i, j)
		}
	}
	return resultado
}

// obtenerQR descompone una matriz usando factorización QR
func obtenerQR(matriz [][]float64) ([][]float64, [][]float64) {
	filas, columnas := len(matriz), len(matriz[0])
	elementos := make([]float64, 0, filas*columnas)
	for _, fila := range matriz {
		elementos = append(elementos, fila...)
	}

	original := mat.NewDense(filas, columnas, elementos)

	var qr mat.QR
	qr.Factorize(original)

	var Q mat.Dense
	var R mat.Dense
	qr.QTo(&Q)
	qr.RTo(&R)

	return convertirAMatriz(&Q), convertirAMatriz(&R)
}

func main() {
	app := fiber.New()
	// Endpoint POST que recibe una matriz, la rota y aplica factorización QR.
	app.Post("/qr", func(c *fiber.Ctx) error {
		var req MatrixRequest
		// Validación de entrada
		if err := c.BodyParser(&req); err != nil {
			log.Println("Error al parsear la solicitud:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El formato de la solicitud no es válido.",
			})
		}

		if len(req.Data) == 0 || len(req.Data[0]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "La matriz no puede estar vacía.",
			})
		}

		// Realiza la rotación y descomposición QR
		rotada := rotaterMatriz(req.Data)
		Q, R := obtenerQR(rotada)
		qr := QRResponse{Q: Q, R: R}

		jsonQR, err := json.Marshal(qr)
		if err != nil {
			log.Println("Error serializando QR:", err)
			return c.Status(500).SendString("Error interno al preparar los datos.")
		}
		// Envío del resultado al servicio Node
		resp, err := http.Post("http://node-api:3000/stats", "application/json", bytes.NewBuffer(jsonQR))
		if err != nil {
			log.Printf("Fallo la conexión con Node.js: %v", err)
			return c.Status(500).SendString("No se pudo contactar al servicio de estadísticas.")
		}
		defer resp.Body.Close()

		//Reponse
		var resultado map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
			log.Println("Error decodificando respuesta de Node.js:", err)
			return c.Status(500).SendString("Respuesta del servicio de estadísticas no es válida.")
		}

		return c.Status(resp.StatusCode).JSON(resultado)
	})

	log.Println("Servidor Go corriendo en el puerto 3001")
	log.Fatal(app.Listen(":3001"))
}
