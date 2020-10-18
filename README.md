# Tarea 1 Sistemas Distribuidos

## Integrantes:
* Diego Moraga Araya 201773035-8 [CC]
* Cristóbal Zambrano Lagos 201773005-6 [CC]

### Distribución de Entidades en las VM's
- Máquina 1: 10.10.28.41 
  - **Entidad Finanzas**
- Máquina 2: 10.10.28.42 
  - **Entidad Logística**
- Máquina 3: 10.10.28.43 
  - **Entidad Clientes**
  - El sistema consultará mediante consola el tipo de cliente que es (`"pyme", "retail"`).
- Máquina 4: 10.10.28.44 
  - **Entidad Camiones**

### Orden de ejecución
* Primero se debe ejecutar  **Logística**, en segundo lugar **Finanzas**, luego **Clientes** y una vez que estén todos corriendo, ejecutar **Camiones**.

### Instrucciones de ejecución
1. Dirigirse a la carpeta correspondiente a la entidad en la máquina que corresponda
   1. Para *Finanzas*:
        ~~~
        $ cd finanzas
        ~~~
    2. Para *Logística*:
        ~~~
        $ cd logistica
        ~~~
    3. Para *Clientes*:
        ~~~
        $ cd clientes
        ~~~
    4. Para *Camiones*:
        ~~~
        $ cd camiones
        ~~~
2. Ejecutar en el orden antes especificado para cada VM el siguiente comando:
   ~~~
   $ make run
   ~~~

*Nota: En caso de cambiar los archivos .csv para distintos pedidos, basta con reemplazarlos y mantener los mismos nombres.*