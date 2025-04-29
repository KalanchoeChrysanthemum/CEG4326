package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"bufio"
	"time"
	"fmt"
	"log"

	db "ic-hardware/database"
	"github.com/tarm/serial"
	"gocv.io/x/gocv"
	"ic-hardware/vision"
)

// CONSTANTS
const key = "0e6670a2bee3b4b7f25d64d66582a776"

func computeHMAC(message []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(message)
	hash := h.Sum(nil)
	return hash[:16]
}

func verifyHMAC(message, messageMAC []byte) bool {
	expectedMAC := computeHMAC(message)
	return hmac.Equal(messageMAC, expectedMAC)
}

func VerifyUser(widHex string, hashHex string) (bool, error) {
    widBytes, err := hex.DecodeString(widHex)
    if err != nil {
        return false, fmt.Errorf("invalid WID hex: %w", err)
    }

    hashBytes, err := hex.DecodeString(hashHex)
    if err != nil {
        return false, fmt.Errorf("invalid HASH hex: %w", err)
    }

    if !verifyHMAC(widBytes, hashBytes) {
        return false, nil
    }

    var wid [16]byte
    copy(wid[:], widBytes)

    user, err := db.QueryUser(wid)
    if err != nil {
        return false, fmt.Errorf("user not found: %w", err)
    }

    if !hmac.Equal(user.HASH[:16], hashBytes) {
        return false, nil
    }

    return true, nil
}

func process(wid, hash string) bool {
    valid, err := VerifyUser(wid, hash)
    if err != nil {
	log.Fatal(err)
    }

    if valid {
	    return true
    } else {
	    return false
    }
}


func main() {
    webcam, err := gocv.OpenVideoCapture(0)
    if err != nil {
	    log.Fatalf("Error opening video capture: %v", err)
    }
    defer webcam.Close()

    time.Sleep(2 * time.Second)

    inputPath := "captured.jpg"
    outputImgPath := "output.jpg"

    defer db.DB.Close()

    config := &serial.Config{
	Name:		"/dev/cu.usbmodem1301",
	Baud:		9600,
	Size:		8,
	StopBits:	serial.Stop1,
	Parity:		serial.ParityNone,
    }

    port, err := serial.OpenPort(config)
    if err != nil {
	log.Fatal(err)
    }

    defer port.Close()

    reader := bufio.NewReader(port)
    for {
	message, err := reader.ReadString('\n')
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Println(message)

	parts := strings.Split(strings.TrimSpace(message), ",")
	if len(parts) != 2 {
	    log.Println("Invalid message format")
	    continue
	}

	wid := parts[0]
	hash := parts[1]

	fmt.Println("WID:", wid)
	fmt.Println("Hash:", hash)

	if (process(wid, hash)) {
		img := gocv.NewMat()
		defer img.Close()

		if ok := webcam.Read(&img); !ok {
			log.Fatalf("Cannot read image from camera")
		}

		if img.Empty() {
			log.Fatalf("Captured image empty")
		}

		if ok := gocv.IMWrite(inputPath, img); !ok {
			log.Fatalf("Failed to save captured image")
		}

		binary, err := vision.ProcessImage(inputPath, outputImgPath)
		if err != nil {
			log.Fatalf("Failed to process image: %v", err)
		}

		fmt.Println("Binary:\n", binary)
	}
    }
}

