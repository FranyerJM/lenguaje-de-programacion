package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"todo"
)

var archivo_tareas = ".todos.json"

func imprimir_ayuda() {
	fmt.Println("-------------------------------------------------")
	fmt.Println("vamos a ver como usar este programa:")
	fmt.Println("\n  -agregar <tu tarea>\t\tagrega una nueva tarea a la lista.")
	fmt.Println("  -listar\t\t\tte muestra todas las tareas pendientes y completadas.")
	fmt.Println("  -completar <numero>\t\tmarca la tarea como completada.")
	fmt.Println("  -progreso <numero>\t\tmarca la tarea como empezada.")
	fmt.Println("  -eliminar <numero>\t\telimina la tarea de la lista.")
	fmt.Println("  -modificar <numero>\t\tmodifica el titulo y/o fecha de una tarea.")
	fmt.Println("  -help\t\t\t\tmuestra esta ayuda.")
	fmt.Println("\nej: go run ./cmd/todo -agregar \"Entrega del Programa\"")
	fmt.Println("-------------------------------------------------")
}

func main() {

	if nombre_desde_env := os.Getenv("TODO_FILENAME"); nombre_desde_env != "" {
		archivo_tareas = nombre_desde_env
	}

	agregar := flag.Bool("agregar", false, "agrega una nueva tarea")
	completar := flag.Int("completar", 0, "marca una tarea como completada")
	progreso := flag.Int("progreso", 0, "marca una tarea como en progreso")
	modificar := flag.Int("modificar", 0, "modifica el titulo y fecha")
	eliminar := flag.Int("eliminar", 0, "elimina una tarea")
	listar := flag.Bool("listar", false, "lista todas las tareas")
	ayuda := flag.Bool("ayuda", false, "muestra la ayuda")

	flag.Parse()

	tareas := &todo.Tareas{}

	if err := tareas.Cargar(archivo_tareas); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		imprimir_ayuda()
		os.Exit(0)
	}

	switch {
	case *agregar:

		tarea, err := obtener_entrada(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Print("introduce la fecha de entrega (dd/mm) o presiona Enter para no poner fecha: ")
		lector_fecha := bufio.NewReader(os.Stdin)
		fecha_str, _ := lector_fecha.ReadString('\n')
		fecha_str = strings.TrimSpace(fecha_str)
		var fecha_entrega time.Time

		if fecha_str != "" {
			fecha_parseada, err := time.Parse("02/01", fecha_str)
			if err != nil {
				fmt.Fprintln(os.Stderr, "el formato de fecha es dd/mm, por ejemplo: 24/08")
				fmt.Println("hazlo de nuevo")
				os.Exit(1)
			}
			ano_actual := time.Now().Year()
			fecha_entrega = time.Date(ano_actual, fecha_parseada.Month(), fecha_parseada.Day(), 23, 59, 59, 0, time.Local)
		}
		tareas.Agregar(tarea, fecha_entrega)
		err = tareas.Almacenar(archivo_tareas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *completar > 0:
		err := tareas.Completar(*completar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = tareas.Almacenar(archivo_tareas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *progreso > 0:
		err := tareas.Progreso(*progreso)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = tareas.Almacenar(archivo_tareas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *modificar > 0:
		fmt.Print("introduce el nuevo titulo: ")
		nuevoTitulo, err := obtener_entrada(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Print("introduce la nueva fecha de entrega (dd/mm), presiona Enter para dejarla sin cambios o escribir '-' para eliminar la fecha: ")
		lector_fecha := bufio.NewReader(os.Stdin)
		fecha_str, _ := lector_fecha.ReadString('\n')
		fecha_str = strings.TrimSpace(fecha_str)

		if fecha_str == "-" {
			err = tareas.Modificar(*modificar, nuevoTitulo, time.Time{})
		} else if fecha_str != "" {
			fecha_parseada, err2 := time.Parse("02/01", fecha_str)
			if err2 != nil {
				fmt.Fprintln(os.Stderr, "el formato de fecha es dd/mm, por ejemplo: 24/08")
				os.Exit(1)
			}
			ano_actual := time.Now().Year()
			fecha_entrega := time.Date(ano_actual, fecha_parseada.Month(), fecha_parseada.Day(), 23, 59, 59, 0, time.Local)
			err = tareas.Modificar(*modificar, nuevoTitulo, fecha_entrega)
		} else {
			err = tareas.Modificar(*modificar, nuevoTitulo)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = tareas.Almacenar(archivo_tareas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *eliminar > 0:
		err := tareas.Eliminar(*eliminar)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = tareas.Almacenar(archivo_tareas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *listar:
		tareas.Imprimir()

	case *ayuda:
		imprimir_ayuda()
		os.Exit(0)

	default:

		fmt.Fprintln(os.Stdout, "el comando no existe")
		imprimir_ayuda()
		os.Exit(1)
	}

}

func obtener_entrada(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	escaner := bufio.NewScanner(r)
	escaner.Scan()
	if err := escaner.Err(); err != nil {
		return "", err
	}

	texto := escaner.Text()

	if len(texto) == 0 {
		return "", errors.New("no pusiste nada")
	}

	return texto, nil
}
