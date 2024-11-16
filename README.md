


# Transactions Manager App

Este proyecto incluye un entorno Dockerizado para desarrollar y ejecutar la aplicación de gestión de transacciones. Se pueden utilizar scripts para simplificar la administración del entorno en sistemas Linux/MacOS (Bash) y Windows (PowerShell).

---

## **Requisitos Previos**

1. **Docker**:


2. **docker-compose**

3. **Scripts disponibles**:
   - **compose.sh**: Para sistemas Linux/macOS.
   - **compose.ps1**: Para sistemas Windows con PowerShell.

4. **Archivo `.env`**:
   - Configura el archivo `.env` en la raíz del proyecto con las siguientes variables:
     ```plaintext
     PORT=
     MONGO_USERNAME=
     MONGO_PASSWORD=
     MONGO_DB_NAME=
     REDIS_HOST=redis
     REDIS_PORT=6379
     REDIS_PASSWORD=
     ENV=development
     ```

---

## **Uso del Script en Linux/macOS**

### 1. **Configurar el archivo `compose.sh`**
   Asegúrate de que el script tenga permisos de ejecución:
   ```bash
   chmod +x compose.sh
   ```

### 2. **Comandos disponibles**

- **Levantar el entorno**:
  ```bash
  ./compose.sh up
  ```
  Esto construye y levanta todos los servicios definidos en `docker-compose.yml`.

- **Detener los contenedores**:
  ```bash
  ./compose.sh down
  ```

- **Limpiar recursos de Docker**:
  ```bash
  ./compose.sh clean
  ```

---

## **Uso del Script en Windows**

### 1. **Ejecutar el script `compose.ps1`**
   Asegúrate de que tienes permisos para ejecutar scripts de PowerShell. Si no, habilítalos ejecutando este comando como administrador:
   ```powershell
   Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned
   ```

### 2. **Comandos disponibles**

- **Levantar el entorno**:
  ```powershell
  .\compose.ps1 up
  ```
  Esto construye y levanta todos los servicios definidos en `docker-compose.yml`.

- **Detener los contenedores**:
  ```powershell
  .\compose.ps1 down
  ```

- **Limpiar recursos de Docker**:
  ```powershell
  .\compose.ps1 clean
  ```

---

## **Flujo de Trabajo**

1. Configura el archivo `.env` con las variables de entorno necesarias.
2. Usa los scripts para manejar el entorno Dockerizado:
   - Usa `up` para iniciar.
   - Usa `down` para detener.
   - Usa `clean` para limpiar recursos no utilizados.
3. Verifica que la aplicación esté corriendo accediendo a [http://localhost:9000](http://localhost:9000) en tu navegador.

---

## **Notas**

- **Cambio de Entorno**:
  - Cambia la variable `ENV` en el archivo `.env` para alternar entre desarrollo (`development`) y producción (`production`).

- **Problemas comunes**:
  - **Error: `.env` not found**:
    - Asegúrate de que el archivo `.env` existe en la raíz del proyecto.
  - **Permiso denegado al ejecutar scripts**:
    - En Linux/macOS:
      ```bash
      chmod +x compose.sh
      ```
    - En Windows:
      Habilita la ejecución de scripts con:
      ```powershell
      Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned
      ```

---
