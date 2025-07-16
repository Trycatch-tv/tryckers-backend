# tryckers-backend

Backend del directorio de Tryckers



## ⚙️ Instalación y ejecución del proyecto

1. Instalación de dependencias
   Para instalar las dependencias del proyecto, ejecuta el siguiente comando: 
```bash
go mod tidy
```

📌 Si tu terminal no reconoce el comando go, debes instalar Go desde: https://golang.org/dl/

2. copiar y pegar el archivo .env.example en la raiz del proyecto y luego renombrarlo como .env 
   y configurar las varibles

## Ejecución del proyecto

1. ejecutar el comando  en la raiz del proyecto (deben abrir primero docker desktop)
```bash
docker compose up -d
```

### Tienes dos maneras de ejecutar este proyecto:

#### 1. 🔹 Opción 1: normal
Desde la raíz del proyecto, ejecuta:
```bash
go run src/cmd/main.go
```

#### 2. 🔹 Opción 2: Modo desarrollo (dev watch) con "air"
Esta opción es totalmente opcional, pero mejora la experiencia de desarrollo. air reinicia automáticamente la aplicación cuando detecta cambios en los archivos, evitando tener que detener y reiniciar manualmente el servicio.

   ⚠️ La siguiente configuración es específica para Windows. Si estás en Linux o macOS, consulta cómo hacerlo en tu sistema operativo.

🪟 Configuración de Air en Windows
Instala Air con el siguiente comando:
   ```bash
go install github.com/air-verse/air@latest
```
Agrega la carpeta go/bin al PATH de tus variables de entorno para que el sistema reconozca el comando air.

La ruta suele estar en una ubicación como:
#"C:\Users\tu_usuario\go\bin"

Para agregar esta ruta al PATH:
Abre el menú de inicio y busca "Editar las variables de entorno del sistema".
Haz clic en "Variables de entorno".
En la sección Variables del sistema o Variables de usuario, busca la variable llamada Path.
Haz clic en Editar, luego en Nuevo, y pega la ruta anterior.
Guarda los cambios y cierra.
Abre una nueva terminal y, desde la raíz del proyecto, ejecuta en la raiz del proyecto:
   ```bash
      air
```

## API Documentation (Swagger)

This project uses [swaggo/swag](https://github.com/swaggo/swag) to auto-generate OpenAPI (Swagger) documentation from code comments.

### How to generate/update the docs

1. Install swag CLI (only once per machine):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate the documentation:
   ```bash
   swag init -g src/cmd/main.go -o docs
   ```

3. Start the application:
   ```bash
   go run src/cmd/main.go
   ```

4. Access the Swagger UI at:
   ```
   http://localhost:8080/swagger/index.html
   ```

> **Note:**  
> Do not commit the generated files in the `docs/` folder (`docs.go`, `swagger.json`, `swagger.yaml`).  
> These files are auto-generated and should be ignored via `.gitignore`.

### New dependencies

- `github.com/swaggo/swag/cmd/swag` (dev tool, not required in production)
- `github.com/swaggo/gin-swagger`
- `github.com/swaggo/files`

## 📄 Generar documentación Swagger

Para generar la documentación Swagger a partir de las anotaciones en el código, ejecuta desde la raíz del proyecto:

```
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g src/internal/api/routes/router.go
```

Esto generará la carpeta `docs/` con la documentación. Puedes consultar la documentación en el endpoint `/swagger/index.html` si tienes integrado Swagger UI en tu servidor.


## Testing 
Tener encuenta que todos los comandos se deben ejecuntar desde la raiz del proyecto 
1. configuarar las variables de entorno de testing para eso crear un archivo llamado .env.test teniendo como ejemplo
las variables que estan en el archivo .env.test.example
2. levantar la db de testing para eso ejecutar el siguiente comando 
```bash
docker compose --env-file .env.test -f docker-compose.test.yml up -d
```
3. ejecutar los test 
```bash
go test ./src/internal/tests
```
4. si quieres ver todos los logs 
```bash
go test ./src/internal/tests -v
```
