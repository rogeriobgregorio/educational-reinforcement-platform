package pkg

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateUUID cria um UUID v4 aleatório.
func GenerateUUID() (string, error) {
	var uuid [16]byte

	// Preencher o UUID com valores aleatórios
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", fmt.Errorf("[GenerateUUID] ERROR: %w", err)
	}

	// Setando a versão do UUID (v4) em 6 bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // 0100 0000 (versão 4)

	// Setando os bits da variante
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 1000 0000 (variante 1)

	// Retornar o UUID no formato adequado (string)
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:],
	), nil
}
