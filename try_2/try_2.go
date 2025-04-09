package try_2

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

// Estado representa los posibles estados de un proceso.
type Estado int

// Enumeración de los estados posibles de un proceso.
const (
	Nuevo Estado = iota
	Listo
	Ejecutando
	Bloqueado
	Saliente
)

// Constantes para simular instrucciones.
const (
	InstruccionNormal = "INSTR"
	InstruccionES     = "ES"
	InstruccionFin    = "FIN"
)

// Proceso define la estructura de un proceso en el sistema.
type Proceso struct {
	ID            int      // Identificador del proceso.
	Estado        Estado   // Estado actual del proceso.
	ContadorPC    int      // Contador del Programa (Program Counter).
	EstadoES      string   // Información del estado de E/S.
	Instrucciones []string // Lista de instrucciones del proceso.
	TiempoBloqueo int      // Tiempo restante en estado Bloqueado.
}

// BCP (Bloque de Control de Programa) almacena la información de control de un proceso.
type BCP struct {
	Proceso *Proceso // Referencia al proceso.
}

// Dispatcher maneja las colas de procesos y la asignación de la CPU.
type Dispatcher struct {
	Listos     *list.List // Cola para procesos en estado Listo.
	Bloqueados *list.List // Cola para procesos en estado Bloqueado.
}

// NuevoDispatcher crea una nueva instancia de Dispatcher.
func NuevoDispatcher() *Dispatcher {
	return &Dispatcher{
		Listos:     list.New(),
		Bloqueados: list.New(),
	}
}

// SimularInstruccion simula la ejecución de una instrucción de un proceso.
func SimularInstruccion(proceso *Proceso, instruccion string) {
	switch {
	case instruccion == InstruccionFin:
		// Finalizar el proceso.
		fmt.Printf("Proceso %d finaliza.\n", proceso.ID)
		proceso.Estado = Saliente
	case strings.HasPrefix(instruccion, InstruccionES):
		// Manejar una instrucción de E/S.
		fmt.Printf("Proceso %d bloqueado por E/S.\n", proceso.ID)
		proceso.Estado = Bloqueado
		// Extraer el tiempo de bloqueo de la instrucción. Ejemplo: "ES3" -> 3 ciclos.
		tiempoBloqueo, err := strconv.Atoi(instruccion[2:])
		if err != nil {
			// Manejar error, por ejemplo, establecer un tiempo de bloqueo predeterminado.
			tiempoBloqueo = 1
		}
		proceso.TiempoBloqueo = tiempoBloqueo
	default:
		// Ejecutar una instrucción normal.
		fmt.Printf("Proceso %d ejecuta instrucción: %s\n", proceso.ID, instruccion)
	}
}

// Metodo para guardar el estado de un proceso en el BCP (ST).
func (d *Dispatcher) GuardarEstado(proceso *Proceso) {
	fmt.Printf("Guardando estado del proceso %d\n", proceso.ID)
	// Aquí se simularía guardar el estado actual del proceso en el BCP.
}

// Metodo para poner un proceso en la cola correspondiente (PUSH).
func (d *Dispatcher) PonerEnCola(proceso *Proceso) {
	switch proceso.Estado {
	case Listo:
		d.Listos.PushBack(proceso)
	case Bloqueado:
		d.Bloqueados.PushBack(proceso)
	default:
		fmt.Printf("El proceso %d no está en un estado que requiera encolamiento.\n", proceso.ID)
	}
	fmt.Printf("Proceso %d puesto en cola %v\n", proceso.ID, proceso.Estado)
}

// Metodo para sacar el siguiente proceso de la cola de listos (PULL).
func (d *Dispatcher) SacarDeColaListos() *Proceso {
	if d.Listos.Len() == 0 {
		return nil
	}
	e := d.Listos.Front()
	d.Listos.Remove(e)
	proceso := e.Value.(*Proceso)
	fmt.Printf("Proceso %d sacado de la cola de listos\n", proceso.ID)
	return proceso
}

// Metodo para cargar el proceso en la CPU (LOAD).
func (d *Dispatcher) CargarProceso(proceso *Proceso) {
	fmt.Printf("Cargando proceso %d en la CPU\n", proceso.ID)
	proceso.Estado = Ejecutando
}

// Metodo para ejecutar el proceso (EXEC).
func (d *Dispatcher) EjecutarProceso(proceso *Proceso, m int) {
	fmt.Printf("Ejecutando proceso %d\n", proceso.ID)
	// Ejecuta hasta 'm' instrucciones o hasta que el proceso esté bloqueado o finalice.
	for i := 0; i < m && proceso.Estado == Ejecutando; i++ {
		if proceso.ContadorPC >= len(proceso.Instrucciones) {
			proceso.Estado = Saliente
			break
		}
		instruccion := proceso.Instrucciones[proceso.ContadorPC]
		proceso.ContadorPC++
		SimularInstruccion(proceso, instruccion)
		// Si el proceso se bloquea o finaliza, salir del ciclo.
		if proceso.Estado != Ejecutando {
			break
		}
	}
	// Cambiar el estado del proceso según su condición actual.
	if proceso.Estado == Ejecutando {
		proceso.Estado = Listo
	}
}

func main() {
	// Crear un Dispatcher.
	disp := NuevoDispatcher()

	// Crear algunos procesos de prueba.
	procesos := []*Proceso{
		{ID: 1, Estado: Nuevo, ContadorPC: 0, EstadoES: "", Instrucciones: []string{"INSTR1", "INSTR2", "ES3", "FIN"}},
		{ID: 2, Estado: Nuevo, ContadorPC: 0, EstadoES: "", Instrucciones: []string{"INSTR1", "ES2", "INSTR3", "FIN"}},
	}

	// Añadir procesos a la cola de listos.
	for _, p := range procesos {
		p.Estado = Listo
		disp.PonerEnCola(p)
	}

	// Simulación de ejecución de procesos.
	for disp.Listos.Len() > 0 || disp.Bloqueados.Len() > 0 {
		// Sacar un proceso de la cola de listos y ejecutarlo.
		proceso := disp.SacarDeColaListos()
		if proceso != nil {
			disp.CargarProceso(proceso)
			// Supongamos que 'm' es 2 para este ejemplo.
			disp.EjecutarProceso(proceso, 2)
			// Si el proceso no ha finalizado, vuelve a ponerse en la cola correspondiente.
			if proceso.Estado != Saliente {
				disp.PonerEnCola(proceso)
			}
		}
		for {
			// ... código existente ...

			// Manejar procesos en la cola de Bloqueados.
			for e := disp.Bloqueados.Front(); e != nil; e = e.Next() {
				proceso := e.Value.(*Proceso)
				proceso.TiempoBloqueo--
				if proceso.TiempoBloqueo <= 0 {
					// Mover el proceso de vuelta a la cola de Listos.
					disp.Listos.PushBack(proceso)
					// Eliminar de la cola de Bloqueados.
					disp.Bloqueados.Remove(e)
				}
			}

			// Verificar si la simulación debe terminar.
			if disp.Listos.Len() == 0 && disp.Bloqueados.Len() == 0 {
				// Verificar si todos los procesos han finalizado.
				// Si es así, romper el ciclo.
			}
		}
	}

	// La lógica para finalizar la simulación y manejar la salida iría aquí.
}
