package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

type EstadoTarea int

const (
	Pendiente EstadoTarea = iota
	Progreso
	Completada
)

type elemento struct {
	Tarea         string
	Estado        EstadoTarea
	Creado_en     time.Time
	Fecha_entrega time.Time
	Completado_en time.Time
}

type Tareas []elemento

func (t *Tareas) Agregar(tarea string, fecha_entrega ...time.Time) {
	todo := elemento{
		Tarea:         tarea,
		Estado:        Pendiente,
		Creado_en:     time.Now(),
		Fecha_entrega: fecha_entrega[0],
		Completado_en: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Tareas) Completar(indice int) error {
	ls := *t
	if indice <= 0 || indice > len(ls) {
		return errors.New("indice incorrecto")
	}

	ls[indice-1].Completado_en = time.Now()
	ls[indice-1].Estado = Completada

	return nil
}

func (t *Tareas) Progreso(indice int) error {
	ls := *t
	if indice <= 0 || indice > len(ls) {
		return errors.New("indice incorrecto")
	}

	ls[indice-1].Estado = Progreso
	return nil
}

func (t *Tareas) Modificar(indice int, titulo string, fecha_entrega ...time.Time) error {
	ls := *t
	if indice <= 0 || indice > len(ls) {
		return errors.New("indice incorrecto")
	}

	if strings.TrimSpace(titulo) != "" {
		ls[indice-1].Tarea = titulo
	}

	if len(fecha_entrega) > 0 {
		ls[indice-1].Fecha_entrega = fecha_entrega[0]
	}

	return nil
}

func (t *Tareas) Eliminar(indice int) error {
	ls := *t
	if indice <= 0 || indice > len(ls) {
		return errors.New("indice incorrecto")
	}

	*t = append(ls[:indice-1], ls[indice:]...)

	return nil
}

func (t *Tareas) Cargar(nombre_archivo string) error {
	archivo, err := ioutil.ReadFile(nombre_archivo)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(archivo) == 0 {
		return err
	}
	err = json.Unmarshal(archivo, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tareas) Almacenar(nombre_archivo string) error {

	datos, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(nombre_archivo, datos, 0644)
}

func (t *Tareas) Imprimir() {

	tabla := simpletable.New()

	tabla.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Tarea"},
			{Align: simpletable.AlignCenter, Text: "listo?"},
			{Align: simpletable.AlignRight, Text: "Fecha Creación"},
			{Align: simpletable.AlignRight, Text: "Fecha Entrega"},
			{Align: simpletable.AlignRight, Text: "Fecha Completado"},
		},
	}

	var celdas [][]*simpletable.Cell

	for idx, item := range *t {
		idx++
		tarea := azul(item.Tarea)
		hecho := azul("ño")
		if item.Estado == Completada {
			tarea = verde(fmt.Sprintf(item.Tarea))
			hecho = verde("shi")
		}

		if item.Estado == Progreso {
			tarea = amarillo(fmt.Sprintf(item.Tarea))
			hecho = amarillo("empezó")
		}

		if item.Fecha_entrega.Before(time.Now()) && item.Estado != Completada {
			tarea = rojo(fmt.Sprintf(item.Tarea))
		}

		fecha_entrega_str := strings.Join(strings.Split(item.Fecha_entrega.Format(time.RFC1123), " ")[:3], " ")
		if item.Fecha_entrega.IsZero() {
			fecha_entrega_str = gris("-")
		}

		celdas = append(celdas, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: tarea},
			{Text: hecho},
			{Text: strings.Join(strings.Split(item.Creado_en.Format(time.RFC1123), " ")[:3], " ")},
			{Text: fecha_entrega_str},
			{Text: strings.Join(strings.Split(item.Completado_en.Format(time.RFC1123), " ")[:3], " ")},
		})
	}

	tabla.Body = &simpletable.Body{Cells: celdas}

	tabla.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{
				Align: simpletable.AlignCenter,
				Span:  6,
				Text: rojo(fmt.Sprintf(
					"Tareas pendientes: %d                         Tareas atrasadas: %d ",
					t.Contar_pendientes(),
					t.Contar_atrasadas(),
				)),
			},
		},
	}

	tabla.SetStyle(simpletable.StyleUnicode)

	tabla.Println()
}

func (t *Tareas) Contar_atrasadas() int {
	total := 0
	ahora := time.Now()
	for _, item := range *t {
		if item.Estado != Completada && !item.Fecha_entrega.IsZero() && item.Fecha_entrega.Before(ahora) {
			total++
		}
	}
	return total
}

func (t *Tareas) Contar_pendientes() int {
	total := 0
	for _, item := range *t {
		if item.Estado == Pendiente {
			total++
		}
	}

	return total
}
