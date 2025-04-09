package try_1

import (
	"container/list"
	"fmt"
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

// Proceso define la estructura de un proceso en el sistema.
type Proceso struct {
	ID            int      // Identificador del proceso.
	Estado        Estado   // Estado actual del proceso.
	ContadorPC    int      // Contador del Programa (Program Counter).
	EstadoES      string   // Información del estado de E/S.
	Instrucciones []string // Lista de instrucciones del proceso.
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

func NuevoDispatcher() *Dispatcher {
	return &Dispatcher{
		Listos:     list.New(),
		Bloqueados: list.New(),
	}
}

// Metodo para guardar el estado de un proceso en el BCP (ST).
func (d *Dispatcher) GuardarEstado(proceso *Proceso) {
	// Aquí se simularía guardar el estado actual del proceso.
	fmt.Printf("Guardando estado del proceso %d\n", proceso.ID)
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
	// La CPU es simulada, así que esta función es simbólica y su "carga" es simplemente lógica.
	proceso.Estado = Ejecutando
}

// Metodo para ejecutar el proceso (EXEC).
func (d *Dispatcher) EjecutarProceso(proceso *Proceso, m int) {
	fmt.Printf("Ejecutando proceso %d\n", proceso.ID)
	// Ejecuta hasta 'm' instrucciones o hasta que el proceso esté bloqueado o finalice.
	for i := 0; i < m && proceso.Estado == Ejecutando; i++ {
		instruccion := proceso.Instrucciones[proceso.ContadorPC]
		proceso.ContadorPC++
		fmt.Printf("Proceso %d ejecutando instrucción %d: %s\n", proceso.ID, proceso.ContadorPC, instruccion)
		// Aquí debería implementarse la lógica para manejar instrucciones de E/S y terminación.
		// SimularInstruccion(proceso, instruccion) // Esto sería una llamada a la función de simulación de instrucción.
	}
	// Simular el cambio de estado al finalizar el quantum.
	if proceso.Estado == Ejecutando {
		proceso.Estado = Listo
	}
}

func main() {
	// Crear un Dispatcher.
	disp := NuevoDispatcher()
	fmt.Println(disp)

	// Crear algunos procesos de prueba.
	// En una implementación real, estos serían creados y cargados desde los archivos de entrada.
	procesos := []*Proceso{
		{ID: 1, Estado: Nuevo, ContadorPC: 0, EstadoES: "", Instrucciones: []string{"INSTR1", "INSTR2", "ES3", "INSTR4"}},
		{ID: 2, Estado: Nuevo, ContadorPC: 0, EstadoES: "", Instrucciones: []string{"INSTR1", "INSTR2", "ES3", "INSTR4"}},
	}

	// Añadir procesos a la cola de listos.
	for _, p := range procesos {
		p.Estado = Listo
		disp.PonerEnCola(p)
	}

	// Simulación de ejecución de procesos.
	for disp.Listos.Len() > 0 {
		// Sacar un proceso de la cola de listos y ejecutarlo.
		proceso := disp.SacarDeColaListos()
		if proceso != nil {
			disp.CargarProceso(proceso)
			// Supongamos que 'm' es 2 para este ejemplo.
			disp.EjecutarProceso(proceso, 2)
			// Si el proceso no ha finalizado, vuelve a ponerse en la cola de listos.
			if proceso.Estado == Listo {
				disp.PonerEnCola(proceso)
			}
		}
	}

	// La lógica para finalizar la simulación y manejar la salida iría aquí.
}
