package controllers
/*
import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

var (
	cameraActive      bool
	stopChan          chan bool
	mutex             sync.Mutex
	contadorVisitas   int
	lastDetectionTime time.Time
	lastDetectionRect image.Rectangle
)

func distanciaRects(a, b image.Rectangle) float64 {
	cx1, cy1 := float64(a.Min.X+a.Dx()/2), float64(a.Min.Y+a.Dy()/2)
	cx2, cy2 := float64(b.Min.X+b.Dx()/2), float64(b.Min.Y+b.Dy()/2)

	dx := cx2 - cx1
	dy := cy2 - cy1
	return math.Sqrt(dx*dx + dy*dy)
}

func StartWebcam(c *gin.Context) {
	contadorVisitas = 0

	mutex.Lock()
	defer mutex.Unlock()

	if cameraActive {
		c.JSON(http.StatusOK, gin.H{"message": "La cámara ya está activa"})
		return
	}

	stopChan = make(chan bool)
	cameraActive = true

	go func() {
		webcam, err := gocv.OpenVideoCapture(0)
		// Inicializa el descriptor HOG para detección de personas
		hog := gocv.NewHOGDescriptor()
		hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())
		defer hog.Close()

		if err != nil {
			log.Println("❌ Error al abrir la webcam:", err)
			cameraActive = false
			return
		}
		defer webcam.Close()

		img := gocv.NewMat()
		defer img.Close()

		log.Println("📷 Webcam iniciada")
		var personaPresente bool
		var tiempoUltimaAusencia time.Time

		for {
			select {
			case <-stopChan:
				log.Println("🛑 Webcam detenida")
				return
			default:
				if ok := webcam.Read(&img); !ok || img.Empty() {
					continue
				}

				// Detectar personas
				rects := hog.DetectMultiScale(img)

				log.Printf("🧍 Detecciones: %d\n", len(rects))

				// dentro del bucle
				if len(rects) > 0 {
					if !personaPresente {
						now := time.Now()
						if now.Sub(tiempoUltimaAusencia) > 3*time.Second {
							contadorVisitas++
							personaPresente = true
							log.Printf("👤 Persona detectada. Total: %d\n", contadorVisitas)
						}
					}
				} else {
					if personaPresente {
						tiempoUltimaAusencia = time.Now()
						personaPresente = false
					}
				}

				time.Sleep(100 * time.Millisecond)
			}
		}

	}()

	c.JSON(http.StatusOK, gin.H{"message": "Webcam iniciada"})
}

func StopWebcam(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if !cameraActive {
		c.JSON(http.StatusOK, gin.H{"message": "La cámara no está activa"})
		return
	}

	stopChan <- true
	close(stopChan)
	cameraActive = false

	// Obtener fecha y hora actual
	now := time.Now()
	fecha := now.Format("2006-01-02")
	hora := now.Format("15:04")

	// Crear estructura
	type Visita struct {
		Id         int    `json:"id"`
		Visitantes int    `json:"visitantes"`
		Hora       string `json:"hora"`
		Fecha      string `json:"fecha"`
		Enviado    bool   `json:"enviado"`
	}

	visita := Visita{
		Id:         0,
		Visitantes: contadorVisitas,
		Hora:       hora,
		Fecha:      fecha,
		Enviado:    false,
	}

	payload, err := json.Marshal([]Visita{visita})
	if err != nil {
		log.Println("❌ Error serializando visita:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar la visita"})
		return
	}

	// Enviar POST a /visitas
	resp, err := http.Post("http://localhost:8080/visitas", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("❌ Error haciendo POST a /visitas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":    "Webcam detenida",
			"visitantes": contadorVisitas,
			"error":      "No se pudo registrar la visita",
		})
		return
	}
	defer resp.Body.Close()

	log.Println("✅ Visita enviada a /visitas")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Webcam detenida",
		"visitantes": contadorVisitas,
	})
}
*/