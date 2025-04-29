package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
	"gocv.io/x/gocv"
	db "ic-hardware/database"
	"ic-hardware/vision"
)

// CONSTANTS
const key = "0e6670a2bee3b4b7f25d64d66582a776"
const inputPath = "pictures/captured.jpg"
const outputPath = "pictures/output.jpg"

/*
*
computeHMAC - Calculates and returns a 16-byte HMAC hash of the providedinput

message []byte - The message to be hashed
*/
func computeHMAC(message []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(message)
	hash := h.Sum(nil)
	return hash[:16]
}

/*
*
verifyHMAC - Verifies the hash of a given message matches the supplied hash

message []byte - The message to be hashed and verified

messageMAC - The expected hash being verified against
*/
func verifyHMAC(message, messageMAC []byte) bool {
	expectedMAC := computeHMAC(message)
	return hmac.Equal(messageMAC, expectedMAC)
}

/*
*
 */
func VerifyUser(widHex string, hashHex string) (bool, error) {
	widBytes, err := hex.DecodeString(widHex)
	if err != nil {
		return false, fmt.Errorf("Invalid WID hex: %w\n", err)
	}

	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return false, fmt.Errorf("Invalid HASH hex: %w\n", err)
	}

	if !verifyHMAC(widBytes, hashBytes) {
		return false, nil
	}

	var wid [16]byte
	copy(wid[:], widBytes)

	user, err := db.QueryUser(wid)
	if err != nil {
		return false, fmt.Errorf("User not found: %w\n", err)
	}

	if !hmac.Equal(user.HASH[:16], hashBytes) {
		return false, nil
	}

	return true, nil
}

func VerifyBinary(widHex, binary string) (bool, error) {
    widBytes, err := hex.DecodeString(widHex)
    if err != nil {
	return false, fmt.Errorf("Invalid WID hex: %w\n", err)
    }

    var wid [16]byte
    copy(wid[:], widBytes)

    user, err := db.QueryUser(wid)
    if err != nil {
	return false, fmt.Errorf("User not found: %w\n", err)
    }

    binaryBytes, err := db.BinaryStringToBytes(binary)
    if err != nil {
        return false, fmt.Errorf("Invalid binary: %w\n", err)
    }

    if len(binaryBytes) != len(user.BINARY) {
        return false, fmt.Errorf("Binary length mismatch\n")
    }

    // Compare byte-by-byte
    for i := 0; i < len(binaryBytes); i++ {
        if user.BINARY[i] != binaryBytes[i] {
            return false, fmt.Errorf("Binary values did not match\n")
        }
    }

    return true, nil
}

func process(wid, hash string) (bool, error) {
	valid, err := VerifyUser(wid, hash)

	if valid {
		return true, nil
	} else {
		return false, err
	}
}

func main() {
	defer db.DB.Close()

	// Open raspberry pi camera
	//webcam, err := gocv.OpenVideoCapture(0)
	//if err != nil {
	//	log.Fatalf("[ERROR] Failed to start camera: %v", err)
	//}
	//defer webcam.Close()

	// Give time buffer to allow camera to fully open
	//time.Sleep(2 * time.Second)

	// Config for serial port
	config := &serial.Config{
		Name:     "/dev/ttyACM0",
		Baud:     9600,
		Size:     8,
		StopBits: serial.Stop1,
		Parity:   serial.ParityNone,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open serial port: %v", err)
	}

	defer port.Close()

	// Open buffered reader for serial communication
	reader := bufio.NewReader(port)

	/**
	  Main logic loop

	  Continuously wait for Arduino to send WID,HASH pairs
	  to validate

	  Once recieved, extract the values, validate the hash, process LEDs
	*/
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[ERROR] Failed to read message: %v\n", err)
		}

		hashStart := time.Now()

		// Split message into 'parts' to separate wid and hash
		parts := strings.Split(strings.TrimSpace(message), ",")
		if len(parts) != 2 {
			log.Println("[ERROR] Invalid message format")
			continue
		}

		wid := parts[0]
		hash := parts[1]

		// Print recieved WID and hash for debugging purposes
		fmt.Println("=========================")
		fmt.Printf("WID:      %-16s\n", wid)
		fmt.Printf("Hash:     %-16s\n", hash)
    
		// Validate hash is valid before capturing and processing an image
		valid, err := process(wid, hash)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
		}

		if valid {
			elapsed := time.Since(hashStart)
			fmt.Printf("[INFO] Validating hash took: %s\n", elapsed)
	
			imgStart := time.Now()
			// Initialize image
			img := gocv.IMRead(inputPath, gocv.IMReadColor)

			if img.Empty() {
			    log.Printf("[ERROR] Captured image is empty: %v\n", err)
			    continue
			}

			binary, err := vision.ProcessImage(inputPath, outputPath)
			img.Close()
			if err != nil {
			    log.Printf("[ERROR] Failed to extract binary from captured image: %v\n", err)
			    continue
			}

			fmt.Printf("Binary:   %-16s\n", binary)

			success, err := VerifyBinary(wid, binary)
			if err != nil {
			    fmt.Printf("[INVALID USER] %v\n", err)
			    continue
			}

			if success {
			    fmt.Println("[VALID USER]")
			}

			imgElapsed := time.Since(imgStart)
			fmt.Printf("[INFO] Image processing took: %s\n", imgElapsed)
		} else {
		    fmt.Println("[INVALID USER]")
		}

		fmt.Println("=========================")
	}
}
