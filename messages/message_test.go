package messages

import (
	"net"
	"testing"
)

func TestOiCorrectMessage(t *testing.T) {
	expected := "nome_teste"
	data := make([]byte, 261)
	data[0] = byte(uint8(1))
	data[1] = byte(uint8(2))
	data[2] = byte(uint8(3))
	data[3] = byte(uint8(4))
	data[4] = byte(uint8(0))
	copy(data[4:], []byte(expected))
	msg := OiMessage(data)
	ip := net.ParseIP("1.2.3.4")
	if !ip.Equal(msg.Addr()) {
		t.Errorf("Erro na interpretação de endereço IP")
	}

	if name := msg.Name(); name != expected {
		t.Errorf("Erro na intepretação de nome de usuário: %s -> %s", expected, name)
	}
}
