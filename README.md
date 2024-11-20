
# **Transactions App**


### **Requisitos Previos**
- **Docker** / **Docker Compose** instalados.

### **Clonar el Repositorio**
```bash
git clone https://github.com/YAFcod3/transactions-manager.git
```

```bash
cd transactions-manager
```

### **Configurar Variables de Entorno**
```env
PORT=
URL_API_EXTERNAL_GET_RATE=https://concurso.dofleini.com/exchange-rate/api/
SUPPORTED_CURRENCIES=USD,EUR,GBP,JPY,CAD,AUD
MONGO_USERNAME=
MONGO_PASSWORD=
MONGO_DB_NAME=
MONGO_PORT_EXTERNAL=27017
REDIS_PORT_EXTERNAL=6379
REDIS_PASSWORD=
APP_JWT_SECRET=
ALLOW_ORIGINS=   (opcional)
```

### **Levantar la Aplicación**
Utilizar los siguientes scripts para iniciar la aplicación con Docker Compose.

### **PowerShell (Windows)**
```powershell
# Iniciar
.\compose.ps1 up

# Detener y eliminar contenedores
.\compose.ps1 down

```

### **Bash (Linux/MacOS)**
```bash
# Iniciar
./compose.sh up

# Detener y eliminar contenedores
./compose.sh down

```
---

## **Descripción del Proyecto**

Esta aplicación es una solución backend diseñada para gestionar transacciones financieras de manera eficiente y segura. Sus funcionalidades incluyen:
- Conversión de monedas con tasas de cambio actualizadas en tiempo real.
- Generación de códigos únicos para cada transacción.
- Gestión de tipos de transacciones
- Estadísticas de transacciones.

### **Endpoints**

##### **- `/exchange/api/conversion`**
##### **- `/exchange/api/currencies`**
##### **- `/exchange/api/settings/transactions-types`**
##### **- `/exchange/api/statistics/`**
##### **- `/exchange/api/transactions`**

---


## **Arquitectura General**
- **Patrón**: Modular con separación de responsabilidades.
- **Componentes Clave**:
  - **Presentación**: Manejo de rutas y middlewares.
  - **Negocio**: Servicios para lógica y validaciones.Se centralizan las operaciones principales del dominio.La lógica está desacoplada de los controladores (handlers)
  - **Persistencia**: **MongoDB**: Se utiliza como base de datos principal para el almacenamiento de información permanente, como transacciones y tipos de transacciones. **Redis**:  Utilizado como almacenamiento en memoria para datos temporales y de acceso rápido, como tasas de cambio actualizadas, contadores para la generación de códigos únicos, y la verificación de transacciones duplicadas.
  - **Utilidades (Utils)**:
    La carpeta  `utils`contiene herramientas reutilizables y funciones auxiliares, organizadas para apoyar los procesos clave de la aplicación tales como :
  - Generación de códigos únicos (generate_transaction_code): Gestiona de forma eficiente los contadores necesarios para generar códigos únicos para cada transacción, almacenando el estado en Redis.
  -  Actualización de tasas de cambio (update_rate): Implementa un proceso automatizado que consulta un servicio externo cada hora para actualizar las tasas de cambio en Redis. Este proceso sigue un enfoque similar al patrón adaptador, integrando el servicio externo con la lógica de la aplicación.



### Esta aplicación se ha desarrollado utilizando Fiber (Framework escrito en Go).
Go es un lenguaje  conocido por su alto rendimiento y capacidades nativas de concurrencia, lo que lo convierte en una excelente opción para este tipo de aplicaciones que pueden manejar un gran volumen de solicitudes simultáneas.
Por otro lado, Fiber, es un framework  para Go que es extremadamente rápido y eficiente, maneja la concurrencia por sí misma gracias a cómo está diseñada.
 Esta aplicación se organiza en módulos claramente definidos, facilitando su mantenimiento y extensibilidad.Cada módulo tiene una única responsabilidad.


## Estructura de carpetas

1. **`app/database/`**
   Contiene toda la configuración relacionada con las bases de datos, incluyendo la inicialización de MongoDB y Redis.

2. **`app/middleware/`**
   Incluye middlewares  como el manejo de autenticación con JWT que protege las rutas  de la API y el middleware que verifica que las transacciones no se duplican en un lapso de tiempo de 20 segundos.

3. **`app/models/`**
   Define las estructuras utilizadas en la aplicación. Estas estructuras representan las colecciones de la base de datos.

4. **`app/routes/`**
   Aquí se definen y registran todas las rutas de la API. Mantener las rutas en un solo lugar facilita su administración y escalabilidad.

5. **`app/handlers/`**
   Contiene las funciones que actúan como intermediarias entre las solicitudes HTTP y la lógica de negocio implementada en los servicios.

6. **`app/services/`**
   Contiene la lógica de negocio de la aplicación, encapsulando  las operaciones  de la aplicación. Esta capa sigue el principio de responsabilidad única, asegurando que cada servicio se encargue de una funcionalidad específica.

7. **`app/utils/`**
    Contiene utilidades  como el actualizador de tasas de cambio , la generación de códigos de transacción únicos y otras funciones reutilizables.

8. **`app/server/`**
   Centraliza la configuración del servidor Fiber.

9. **Raíz del proyecto**
   - `main.go`: Punto de entrada principal que inicializa el servidor, servicios y conexiones a la base de datos.
   - `docker-compose.yml`: Facilita el levantamiento de servicios como Redis y MongoDB en entornos de desarrollo.
   - `Dockerfile`: Configuración para la construcción de  la imágen de Docker del servidor de Fiber.



---
## Explicaciones de algunos elementos clave

### **Punto de Entrada Principal**

El archivo `main.go` inicializa las dependencias de la aplicación, como MongoDB y Redis, configura el servidor con Fiber, y gestiona procesos clave como la generación de códigos únicos y la actualización automática de tasas de cambio.


### **Configuración del Servidor**

El archivo `server.go` define la configuración del servidor usando Fiber. Incluye middlewares globales, como CORS y autenticación JWT, y registra las rutas de la aplicación. Esto asegura una estructura clara y modular para manejar solicitudes HTTP.


### **Actualización de Tasas de Cambio**

El archivo `update_rate.go` en la carpeta `utils` implementa un proceso automatizado para mantener actualizadas las tasas de cambio. Este proceso consulta un servicio externo cada hora y almacena las tasas en Redis, utilizando un diseño concurrente para no bloquear la ejecución principal de la aplicación.

#### **Características Principales:**
1. **Consulta Periódica**:
   - Realiza solicitudes al endpoint externo definido en la variable de entorno `URL_API_EXTERNAL_GET_RATE`.
   - Actualiza las tasas de cambio en Redis.

2. **Almacenamiento en Redis**:
   - Almacena las tasas de cambio, la moneda base  en un hash llamado exchange_rates.

3. **Ejecución Concurrente:**:

    - Usa una goroutine para que el proceso de actualización se ejecute en paralelo con el resto de la aplicación, sin bloquear la lógica principal.

    ```go
     go func() {
         for {
        select {
        case <-ticker.C:
            fetchExchangeRates()
        case <-stopChan:
            fmt.Println("Exchange rate updater stopped.")
            return
           }
        }
        }()
   ```
4. **Beneficios de la Implementación**:
    - Concurrencia: La actualización se ejecuta en paralelo, asegurando que no interfiera con otras operaciones del sistema.
    - Eficiencia: Al usar Redis y pipelines, el sistema minimiza la latencia y mejora el rendimiento.
    - Sigue un enfoque adaptador para integrar el servicio externo con Redis, lo que facilita el mantenimiento y posibles cambios futuros.


### **Generador de Códigos Únicos**

El archivo `generate_transaction_code.go` implementa un sistema para generar códigos únicos en el formato `TYYMMDDHH{00000000}`, donde:
- `YYMMDDHH` representa la fecha y hora actual.
- `{00000000}` es un contador incremental, reiniciado cada hora.

---

### **Características Clave**

1. **Persistencia con Redis**:
   - El contador (`counter`) y la última hora procesada (`lastHour`) se almacenan en Redis para garantizar continuidad tras reinicios.

2. **Concurrencia Segura**:
   - Usa `sync.Mutex` para evitar **condiciones de carrera** al manejar el contador en entornos concurrentes.

3. **Reinicio por Hora**:
   - Al cambiar de hora, el contador se reinicia automáticamente. En el lapso d una hora si el contador supera su valor máximo, no se generará ningún código hasta al menos la siguiente hora.


### **Ventajas**
- **Fiabilidad**: Los datos en Redis garantizan la persistencia del estado.
- **Seguridad Concurrente**: Evita conflictos usando un mutex.




## **`docker-compose.yml`**

Este archivo define los servicios, redes y volúmenes necesarios para desplegar y ejecutar la aplicación de manera eficiente utilizando Docker Compose.

### **Servicios**

#### **1. `transactions-manager-app`**
- **Función**: Servicio principal que ejecuta la aplicación.
- **Configuración**:
  - Construye la imagen desde el `Dockerfile` del proyecto.
  - Mapea el puerto definido en `.env` (`PORT`) al puerto del contenedor.
  - Depende de los servicios `transactionsManagerDb` (MongoDB) y `redis` para garantizar que estén disponibles antes de iniciar.
  - Reinicia automáticamente en caso de fallos (`restart: on-failure`).
- **Red**: `transactions-manager-network`.

---

#### **2. `transactionsManagerDb`**
- **Función**: Servicio que ejecuta MongoDB como base de datos principal.
- **Persistencia**:
  - Utiliza el volumen `transactionsManagerDbData` para almacenar los datos permanentemente en el host, asegurando que los datos no se pierdan al reiniciar el contenedor.
- **Configuración**:
  - Usa las credenciales definidas en `.env` para configurar un usuario y una base de datos inicial.
  - Mapea el puerto 27017 al puerto externo definido en `.env`.
  - Incluye un comando para habilitar la autenticación en MongoDB (`--auth`).
  - Verifica su estado con un `healthcheck` para asegurarse de que el servicio está disponible.

---

#### **3. `redis`**
- **Función**: Servicio que ejecuta Redis para el almacenamiento en memoria y el manejo de datos temporales.
- **Persistencia**:
  - Configurado con la opción `--appendonly yes`, lo que garantiza que los datos almacenados en Redis se guarden en disco. Esto previene la pérdida de información (como tasas de cambio o contadores de transacciones) incluso si el contenedor se reinicia.
  - Utiliza el volumen `redis_data` para almacenar los datos en el host.
- **Seguridad**:
  - Configurado con una contraseña (`requirepass`), definida en el archivo `.env`.
- **Configuración**:
  - Mapea el puerto 6379 al puerto externo definido en `.env`.
  - Realiza un `healthcheck` para garantizar que el servicio está disponible y funcionando correctamente.
