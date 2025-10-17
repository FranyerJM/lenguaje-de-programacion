package todo

import "fmt"

const (
	Color_por_defecto = "\x1b[39m"

	Color_rojo     = "\x1b[91m"
	Color_verde    = "\x1b[32m"
	Color_azul     = "\x1b[94m"
	Color_gris     = "\x1b[90m"
	Color_amarillo = "\x1b[93m"
)

func rojo(s string) string {
	return fmt.Sprintf("%s%s%s", Color_rojo, s, Color_por_defecto)
}

func verde(s string) string {
	return fmt.Sprintf("%s%s%s", Color_verde, s, Color_por_defecto)
}

func azul(s string) string {
	return fmt.Sprintf("%s%s%s", Color_azul, s, Color_por_defecto)
}

func gris(s string) string {
	return fmt.Sprintf("%s%s%s", Color_gris, s, Color_por_defecto)
}

func amarillo(s string) string {
	return fmt.Sprintf("%s%s%s", Color_amarillo, s, Color_por_defecto)
}
