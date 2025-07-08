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


## Testing 
Tener encuenta que todos los comandos se deben ejecuntar desde la raiz del proyecto 
1. configuarar la variable de entorno de testing para eso crear un archivo llamado .env.test teniendo como ejemplo
las variable que estan en el archivo .env.test.example
2. levantar la db de testing para eso ejecutar el siguiente comando 
```bash
docker compose --env-file .env.test -f docker-compose.test.yml up -d
```
3. ejecutar lost test 
```bash
go test ./src/internal/tests
```
4. si quieres ver todos los logs 
```bash
go test ./src/internal/tests -v
```

