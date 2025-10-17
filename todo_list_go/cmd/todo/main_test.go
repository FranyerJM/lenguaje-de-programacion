package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName      = "todo"
	testFileName = ".test_todos.json"
)

func TestMain(m *testing.M) {

	os.Setenv("TODO_FILENAME", testFileName)

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		os.Stderr.WriteString("no se pudo construir el binario de prueba\n")
		os.Exit(1)
	}

	exitCode := m.Run()

	// Limpieza
	os.Remove(binName)
	os.Remove(testFileName)

	os.Exit(exitCode)
}

func setupTestFile(t *testing.T) {
	content := []byte(`[{"Tarea":"Tarea de prueba 1","Estado":0,"Creado_en":"2025-10-10T21:13:47.4328525-04:00","Fecha_entrega":"2025-10-11T23:59:59-04:00","Completado_en":"0001-01-01T00:00:00Z"}]`)
	if err := ioutil.WriteFile(testFileName, content, 0644); err != nil {
		t.Fatalf("no se pudo crear el archivo de prueba: %v", err)
	}
}

func TestComandos(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)

	t.Run("AgregarNuevaTarea", func(t *testing.T) {
		os.Remove(testFileName)
		cmd := exec.Command(cmdPath, "-agregar", "nueva tarea desde prueba")
		cmd.Stdin = strings.NewReader("\n")
		if err := cmd.Run(); err != nil {
			t.Fatalf("el comando 'agregar' no funciono: %v", err)
		}
	})

	t.Run("ListarTareas", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-listar")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("el comando 'listar' no funciono: %v", err)
		}

		if !strings.Contains(string(out), "nueva tarea desde prueba") {
			t.Errorf("la salida de 'listar' no contiene la tarea esperada. Salida:\n%s", string(out))
		}
	})

	t.Run("CompletarTarea", func(t *testing.T) {
		setupTestFile(t)
		cmd := exec.Command(cmdPath, "-completar", "1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("el comando 'completar' no funciono: %v", err)
		}

		// Verificar que el estado cambió
		cmdList := exec.Command(cmdPath, "-listar")
		out, _ := cmdList.CombinedOutput()
		if !strings.Contains(string(out), "shi") {
			t.Error("la tarea no se puso como completada ('shi') en la salida de 'listar'.")
		}
	})

	t.Run("EliminarTarea", func(t *testing.T) {
		setupTestFile(t)
		cmd := exec.Command(cmdPath, "-eliminar", "1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("el comando 'eliminar' no funciono: %v", err)
		}

		// Verificar que la tarea ya no existe
		content, err := ioutil.ReadFile(testFileName)
		if err != nil {
			t.Fatalf("no se pudo leer el archivo de prueba: %v", err)
		}
		if string(content) != "[]" {
			t.Errorf("el archivo de tareas no esta vacio despues de eliminar. Contenido: %s", string(content))
		}
	})

	t.Run("SinArgumentosMuestraAyuda", func(t *testing.T) {
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != 0 {
					t.Fatalf("el comando sin argumentosno funciono: %v", err)
				}
			} else {
				t.Fatalf("error al ejecutar el comando sin argumentos: %v", err)
			}
		}

		if !strings.Contains(string(out), "vamos a ver como usar este programa") {
			t.Error("el comandoo sin argumentos no muestra la ayuda esperada")
		}
	})

	t.Run("ModificarTarea", func(t *testing.T) {
		setupTestFile(t)
		cmd := exec.Command(cmdPath, "-modificar", "1", "Tarea modificada")

		var stdin bytes.Buffer

		stdin.WriteString("\n")
		cmd.Stdin = &stdin

		if err := cmd.Run(); err != nil {
			t.Fatalf("el comando 'modificar' no funciono: %v", err)
		}

		cmdList := exec.Command(cmdPath, "-listar")
		out, _ := cmdList.CombinedOutput()
		if !strings.Contains(string(out), "Tarea modificada") {
			t.Errorf("la tarea no fue modificada bien, Salida: \n%s", out)
		}
	})
}
