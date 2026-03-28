package crypto

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

func Decrypt(inputPath, outputPath string, masterKey []byte, useAES bool) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	header := make([]byte, headerSize)
	if _, err := io.ReadFull(inFile, header); err != nil {
		return errors.New("unable to load the file header")
	}
	if string(header[:5]) != "GOFENC" {
		return errors.New("invalid file format")
	}
	if header[5] != 1 {
		return errors.New("unsupported file format")
	}

	aead, err := newAEAD(masterKey, useAES)
	if err != nil {
		return err
	}

	nonceSize := aead.NonceSize()
	chunkIndex := uint64(0)

	for {
		nonce := make([]byte, nonceSize)
		_, err := io.ReadFull(inFile, nonce)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunkLenBuf := make([]byte, 4)
		if _, err := io.ReadFull(inFile, chunkLenBuf); err != nil {
			return errors.New("error reading chunk length")
		}
		chunkLen := binary.LittleEndian.Uint32(chunkLenBuf)

		encrypted := make([]byte, chunkLen)
		if _, err := io.ReadFull(inFile, encrypted); err != nil {
			return errors.New("error reading chunk data")
		}

		ad := make([]byte, 8)
		binary.LittleEndian.PutUint64(ad, chunkIndex)

		plaintext, err := aead.Open(nil, nonce, encrypted, ad)
		if err != nil {
			return errors.New("error decrypting")
		}

		if _, err := outFile.Write(plaintext); err != nil {
			return err
		}

		chunkIndex++
	}

	return nil
}
