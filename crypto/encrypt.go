package crypto

import (
	"encoding/binary"
	"io"
	"os"
)

const (
	chunkSize  = 64 * 1024
	headerSize = 5 + 1
)

func Encrypt(inputPath, outputPath string, masterKey []byte, useAES bool) error {
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

	aead, err := newAEAD(masterKey, useAES)
	if err != nil {
		return err
	}

	if _, err := outFile.Write([]byte("GOFENC")); err != nil {
		return err
	}
	if _, err := outFile.Write([]byte{1}); err != nil {
		return err
	}

	buf := make([]byte, chunkSize)
	chunkIndex := uint64(0)

	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			nonce, nonceErr := GenerateRandom(aead.NonceSize())
			if nonceErr != nil {
				return nonceErr
			}

			ad := make([]byte, 8)
			binary.LittleEndian.PutUint64(ad, chunkIndex)

			encrypted := aead.Seal(nil, nonce, buf[:n], ad)

			if _, err := outFile.Write(nonce); err != nil {
				return err
			}

			chunkLen := make([]byte, 4)
			binary.LittleEndian.PutUint32(chunkLen, uint32(len(encrypted)))
			if _, err := outFile.Write(chunkLen); err != nil {
				return err
			}
			if _, err := outFile.Write(encrypted); err != nil {
				return err
			}
			chunkIndex++
		}
		if err != io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
