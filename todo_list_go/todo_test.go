package todo

import (
	"os"
	"testing"
	"time"
)

func newTasks() *Tareas {
	return &Tareas{
		{Tarea: "Tarea 1", Estado: Pendiente, Creado_en: time.Now(), Fecha_entrega: time.Now().Add(24 * time.Hour)},
		{Tarea: "Tarea 2", Estado: Pendiente, Creado_en: time.Now(), Fecha_entrega: time.Now().Add(-24 * time.Hour)}, // tarea atrasada
		{Tarea: "Tarea 3", Estado: Progreso, Creado_en: time.Now(), Fecha_entrega: time.Now().Add(48 * time.Hour)},
	}
}

func TestAgregar(t *testing.T) {
	tareas := Tareas{}
	fechaEntrega := time.Now().Add(24 * time.Hour)
	tareas.Agregar("nueva tarea", fechaEntrega)

	if len(tareas) != 1 {
		t.Errorf("se esperaba 1 tarea pero entraron %d: ", len(tareas))
	}
	if tareas[0].Tarea != "nueva tarea" {
		t.Errorf("el nombre de la tarea no es 'nueva tarea' si no: %s", tareas[0].Tarea)
	}
	if tareas[0].Estado != Pendiente {
		t.Errorf("el estado deberia ser pendiente pero se puso: %d", tareas[0].Estado)
	}
}

func TestCompletar(t *testing.T) {
	tareas := newTasks()

	// si funciona
	err := tareas.Completar(1)
	if err != nil {
		t.Errorf("no se esperaba un error al completar la tarea: %s", err)
	}
	if (*tareas)[0].Estado != Completada {
		t.Errorf("el estado de la tarea deberia ser completada, pero es %d", (*tareas)[0].Estado)
	}

	// si no funciona
	err = tareas.Completar(99)
	if err == nil {
		t.Error("se esperaba obtener error por indice invalido, pero no hay")
	}
}

func TestProgreso(t *testing.T) {
	tareas := newTasks()

	// si funciona
	err := tareas.Progreso(1)
	if err != nil {
		t.Errorf("no se esperaba error al marcar tarea en progreso: %s", err)
	}
	if (*tareas)[0].Estado != Progreso {
		t.Errorf("el estado de la tarea no es progreso ('empezo'), es: %d", (*tareas)[0].Estado)
	}

	// si no funciona
	err = tareas.Progreso(0)
	if err == nil {
		t.Error("se esperaba obtener error por indice invalido, pero no hay")
	}
}

func TestModificar(t *testing.T) {
	tareas := newTasks()
	nuevaFecha := time.Now().Add(72 * time.Hour)

	err := tareas.Modificar(1, "titulo de la tarea modificado")
	if err != nil {
		t.Errorf("no se esperaba un error al modificar el titulo: %s", err)
	}
	if (*tareas)[0].Tarea != "titulo de la tarea modificado" {
		t.Errorf("el titulo no se modifico bien")
	}

	err = tareas.Modificar(1, "", nuevaFecha)
	if err != nil {
		t.Errorf("no se esperaba un error al modificar la fecha: %s", err)
	}
	if !(*tareas)[0].Fecha_entrega.Equal(nuevaFecha) {
		t.Errorf("La fecha de entrega no se modifico bien")
	}

	err = tareas.Modificar(2, "tarea con ambos modificados", nuevaFecha)
	if err != nil {
		t.Errorf("no se esperaba un error al modificar ambos campos: %s", err)
	}
	if (*tareas)[1].Tarea != "tarea con ambos modificados" {
		t.Errorf("el titulo no se modifico bien")
	}
	if !(*tareas)[1].Fecha_entrega.Equal(nuevaFecha) {
		t.Errorf("La fecha de entrega no se modifico bien")
	}

	// si no funciona
	err = tareas.Modificar(100, "indice que no existe")
	if err == nil {
		t.Error("Se esperaba un error por indice invalido, pero no hay")
	}
}

func TestEliminar(t *testing.T) {
	tareas := newTasks()
	originalLen := len(*tareas)

	// si funciona
	err := tareas.Eliminar(1)
	if err != nil {
		t.Errorf("no se esperaba un error al eliminar la tarea: %s", err)
	}
	if len(*tareas) != originalLen-1 {
		t.Errorf("se esperaba que la lista de tareas tuviera %d elementos, pero tiene %d", originalLen-1, len(*tareas))
	}

	// si no funciona
	err = tareas.Eliminar(99)
	if err == nil {
		t.Error("se esperaba obtener error por indice invalido, pero no hay")
	}
}

func TestAlmacenarYCargar(t *testing.T) {

	tareas := newTasks()
	nombreArchivoTemporal := "test_todos.json"

	defer os.Remove(nombreArchivoTemporal)

	err := tareas.Almacenar(nombreArchivoTemporal)
	if err != nil {
		t.Fatalf("error al almacenar las tareas: %s", err)
	}

	tareasCargadas := &Tareas{}
	err = tareasCargadas.Cargar(nombreArchivoTemporal)
	if err != nil {
		t.Fatalf("error al cargar las tareas: %s", err)
	}

	if len(*tareas) != len(*tareasCargadas) {
		t.Errorf("el numero de tareas almacenadas y cargadas no es el mismo. Esperado: %d, Obtenido: %d", len(*tareas), len(*tareasCargadas))
	}

	if (*tareas)[0].Tarea != (*tareasCargadas)[0].Tarea {
		t.Errorf("la tarea cargada es incorrecta. Esperado: '%s', Obtenido: '%s'", (*tareas)[0].Tarea, (*tareasCargadas)[0].Tarea)
	}
}

func TestContarPendientes(t *testing.T) {
	tareas := newTasks()

	if conteo := tareas.Contar_pendientes(); conteo != 2 {
		t.Errorf("conteo de pendientes incorrecto. Esperado: 2, Obtenido: %d", conteo)
	}

	tareas.Completar(1)
	if conteo := tareas.Contar_pendientes(); conteo != 1 {
		t.Errorf("conteo de pendientes incorrecto tras completar. Esperado: 1, Obtenido: %d", conteo)
	}
}

func TestContarAtrasadas(t *testing.T) {
	tareas := newTasks()

	if conteo := tareas.Contar_atrasadas(); conteo != 1 {
		t.Errorf("conteo de atrasadas incorrecto. Esperado: 1, Obtenido: %d", conteo)
	}

	tareas.Completar(2)
	if conteo := tareas.Contar_atrasadas(); conteo != 0 {
		t.Errorf("conteo de atrasadas incorrecto tras completar. Esperado: 0, Obtenido: %d", conteo)
	}
}
