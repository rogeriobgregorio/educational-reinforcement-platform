package pkg

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

// GenerateUUIDv7 cria um UUID v7 baseado em timestamp e aleatório.
func GenerateUUIDv7() (string, error) {
	var uuid [16]byte

	// Adiciona o timestamp (48 bits de tempo em milissegundos desde a época)
	timestamp := uint64(time.Now().UnixMilli())

	// Preencher os primeiros 6 bytes com o timestamp (48 bits)
	uuid[0] = byte(timestamp >> 40)
	uuid[1] = byte(timestamp >> 32)
	uuid[2] = byte(timestamp >> 24)
	uuid[3] = byte(timestamp >> 16)
	uuid[4] = byte(timestamp >> 8)
	uuid[5] = byte(timestamp)

	// Preencher os bytes restantes com aleatoriedade
	if _, err := io.ReadFull(rand.Reader, uuid[6:]); err != nil {
		return "", fmt.Errorf("[GenerateUUIDv7] ERROR: %w", err)
	}

	// Setando a versão do UUID (v7) em 6 bits no byte 6
	uuid[6] = (uuid[6] & 0x0f) | 0x70 // 0111 0000 (versão 7)

	// Setando os bits da variante no byte 8 (primeiros 2 bits 10)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 1000 0000 (variante 1)

	// Retornar o UUID no formato adequado (string)
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		uint32(uuid[0])<<24|uint32(uuid[1])<<16|uint32(uuid[2])<<8|uint32(uuid[3]),
		uint16(uuid[4])<<8|uint16(uuid[5]),
		uint16(uuid[6])<<8|uint16(uuid[7]),
		uint16(uuid[8])<<8|uint16(uuid[9]),
		uint64(uuid[10])<<40|uint64(uuid[11])<<32|uint64(uuid[12])<<24|uint64(uuid[13])<<16|uint64(uuid[14])<<8|uint64(uuid[15]),
	), nil
}
