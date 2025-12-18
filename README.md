Este proyecto implementa un sistema distribuido de dos servicios:

- Una API en **Go** (usando [Fiber](https://gofiber.io/)) que recibe matriz original como entrada, realizará la rotación de la matriz y luego enviará los datos
resultantes a la segunda API en Node.js...
- Una API en **Node.js** (usando [Express.js](https://expressjs.com/)) que recibe matrices `Q` y `R`, calcula estadísticas básicas, y devuelve un resumen.

Ambas APIs están contenizados con docker y se comunican vía HTTP usando `docker-compose`.



# Tecnologías Usadas
- [Go](https://golang.org/) 1.21 + Fiber
- [Node.js](https://nodejs.org/) 18 + Express
- [Gonum](https://www.gonum.org/) (Go) para factorización QR
- Docker + Docker Compose
- Comunicación entre microservicios vía HTTP


## Requisitos

- Docker
- Docker Compose

## Comandos

Clona el repositorio
- git clone ...
- docker compose up --build


## Modo de uso
##  POST /qr (API Go) por http post

Hacer una peticion post hacia el endpoint de go que es http://localhost:3001/qr y enviar una matriz en el cuerpo de la solicitud como JSON, como el ejemplo mostrado

{ "data": [ [1, 2], [3, 4] ] }

En caso de que el formato no sea el indicado la respuesta que vera será 

{
    "error": "El formato de la solicitud no es válido."
}

En caso de que el formato sea correcto el response code seria 200 y la respuesta seria dependiendo de la matriz que se mande, un ejemplo seria

{
    "average": -1.3416407864998738,
    "isDiagonal": false,
    "max": 0.4472135954999581,
    "min": -6.708203932499369,
    "sum": -10.73312629199899
}


## Consideraciones Técnicas

- La matriz enviada es rotada 90° (antihorario) antes de aplicar la factorización QR.
- La factorización QR se realiza con la biblioteca Gonum (Go).
- El resultado (Q y R) se envía como JSON al servicio Node.js.
- Node.js calcula:
  - Máximo y mínimo valor
  - Suma total y promedio
  - Verifica si alguna matriz es diagonal
- Se sigue un patrón de arquitectura desacoplada (comunicación vía HTTP).
- La solución es portable y ejecutable en cualquier entorno gracias a Docker.


## Notas 


- El `docker-compose.yml` expone los puertos 3000 (Node) y 3001 (Go).
- Ambos servicios están construidos desde cero con su `Dockerfile`.
