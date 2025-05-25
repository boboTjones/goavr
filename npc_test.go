package main

import (
	"testing"
)

func TestNPCInitialization(t *testing.T) {
	port := "8080"
	npc := NewNPC(port)

	if npc.port != port {
		t.Errorf("Expected port %s, got %s", port, npc.port)
	}

	if npc.rx == nil || npc.tx == nil {
		t.Error("Expected non-nil channels for rx and tx")
	}
}
