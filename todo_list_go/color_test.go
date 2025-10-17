package todo

import (
	"fmt"
	"testing"
)

func TestRojo(t *testing.T) {
	s := "texto de prueba"
	esperado := fmt.Sprintf("%s%s%s", Color_rojo, s, Color_por_defecto)
	resultado := rojo(s)
	if resultado != esperado {
		t.Errorf("la funcion del color rojo() no funciono. Esperado: %q, Obtenido: %q", esperado, resultado)
	}
}

func TestVerde(t *testing.T) {
	s := "texto de prueba"
	esperado := fmt.Sprintf("%s%s%s", Color_verde, s, Color_por_defecto)
	resultado := verde(s)
	if resultado != esperado {
		t.Errorf("la funcion del color verde() no funciono. Esperado: %q, Obtenido: %q", esperado, resultado)
	}
}

func TestAzul(t *testing.T) {
	s := "texto de prueba"
	esperado := fmt.Sprintf("%s%s%s", Color_azul, s, Color_por_defecto)
	resultado := azul(s)
	if resultado != esperado {
		t.Errorf("la funcion del color azul() no funciono. Esperado: %q, Obtenido: %q", esperado, resultado)
	}
}

func TestGris(t *testing.T) {
	s := "texto de prueba"
	esperado := fmt.Sprintf("%s%s%s", Color_gris, s, Color_por_defecto)
	resultado := gris(s)
	if resultado != esperado {
		t.Errorf("la funcion del color gris() no funciono. Esperado: %q, Obtenido: %q", esperado, resultado)
	}
}

func TestAmarillo(t *testing.T) {
	s := "texto de prueba"
	esperado := fmt.Sprintf("%s%s%s", Color_amarillo, s, Color_por_defecto)
	resultado := amarillo(s)
	if resultado != esperado {
		t.Errorf("la funcion del color amarillo() no funciono. Esperado: %q, Obtenido: %q", esperado, resultado)
	}
}
