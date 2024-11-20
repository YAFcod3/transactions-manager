
# **Transactions App**

<!-- ## **Levantar la Aplicación** -->

### **Requisitos Previos**
- **Docker** y **Docker Compose** instalados.

### **Clonar el Repositorio**
Clona el repositorio en tu máquina local:
```bash
git clone https://github.com/YAFcod3/transactions-manager.git
cd transactions-manager
```
### **Configurar Variables de Entorno**
```env
PORT=9000
URL_API_EXTERNAL_GET_RATE=https://concurso.dofleini.com/exchange-rate/api/
SUPPORTED_CURRENCIES=USD,EUR,GBP,JPY,CAD,AUD
MONGO_USERNAME=root
MONGO_PASSWORD=
MONGO_DB_NAME=
MONGO_PORT_EXTERNAL=27017
REDIS_PORT_EXTERNAL=6379
REDIS_PASSWORD=
ENV=development
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
- Gestión de tipos de transacciones con eliminación lógica (*soft delete*).

---

## **Características Principales**

- **Conversión de monedas** con tasas actualizadas cada hora desde un proveedor externo.
- **Generación de códigos únicos**, asegurando unicidad y consistencia incluso bajo alta concurrencia.
- **Gestión de tipos de transacciones**, permitiendo creación, actualización y eliminación lógica.

---

## **Diseño Arquitectónico**

### **Arquitectura General**
- **Patrón**: Modular con separación de responsabilidades.
- **Componentes Clave**:
  - **Presentación**: Manejo de rutas y middlewares usando Fiber.
  - **Negocio**: Servicios para lógica y validaciones.
  - **Persistencia**: Uso de MongoDB para datos permanentes y Redis para datos volátiles.
  - **Infraestructura**: Contenedores Docker para despliegue eficiente.

## **Seguridad**

### **Autenticación**
- Middleware basado en JWT para proteger rutas sensibles.

### **Validaciones**
- Validaciones estrictas de entrada para evitar datos inconsistentes.

### **Configuraciones Seguras**
- Uso de variables de entorno para proteger datos sensibles como contraseñas y claves.

---


---
