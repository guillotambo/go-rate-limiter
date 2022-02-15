# Golang Rate Limiter

Se monta un server que expone un endpoint `GET /message` a partir del cual se devuelva un mensaje del servicio de [Fuck Off as a Service](https://www.foaas.com/).

Un usuario puede consumir este endpoint X cantidad de veces dentro de un periodo de cantidad T de segundos.

## Casos de Uso - Ejemplos:

- Se consume una vez la API con un userId determinado y devuelve el mensaje del servicio
- Se consume la API 5 veces dentro de un periodo de 10 segundos y esta devuelve los 5 mensajes del servicio
- Se consume la API 6 veces dentro de un periodo de 10 segundos y el sexto llamado devuelve un error.
- Se consume la API 6 veces dentro de un periodo de 10 segundos, se hace un septimo llamado 10 segundos despues del primer llamado y este devuelve un mensaje del servicio.  

# Desarrollo
Desarrollado en Golang. 

Rate Limiter como middleware, el cual fue usado para el endpoint especificado, pero la idea es que pueda ser reutilizado para cualquier endpoint.  

## Solución
Se mantiene una colección de los timestamps de cada request por cada usuario.  
Cuando llega un request:
- Se identifica a que usuario pertenece
- Se limpian de la colección de timestamps de ese usuario, los que hayan quedados obsoletos para la ventana de tiempo configurada.
- Si el tamaño de la coleccion es menor a la cantidad permitida, se agrega el nuevo timestamp. Se habilita al usuario a utilizar el servicio.
- Si no, el usuario no tiene permitido el uso del servicio, hasta que se cumpla la condición configurada.


## Decisiones de diseño
### Estructura de datos
Para la colección, se utiliza una cola (linked list) ya que permite mantener la colección ordenada, pues el timestamp es generado al momento de insertar en la colección, el cual es siempre insertado al final de la misma. Por esta razón no puede suceder que un timestamp a insertar, sea mas antiguo que uno ya perteneciente a la cola.  

### Unidad de tiempo
La medida de tiempo que utilicé para guardar los timestamps internamente es en nanosegundos, para disminiur la posiblidad de colisión.  
Podría haber usado milisegundos, pero no me parecía medir en segundos, ya que podrían surgir problemas en los casos de uso (Requests llegando en los limites).  

### Clock
La idea de inyectar al rate limiter, la interfaz de un reloj que te devuelve el timestamp actual, es para poder tener control del tiempo al momento de testear la correctitud.  

### Autenticación
Para la autenticación del usuario, al momento de hacer el request, se espera un Header llamado `userId`.  

### Testing
Se realizaron test unitarios para el rate limiter.  


## Requisitos
- SDK de Go 
- Se puede levantar con docker

Ejecutar
- `go build - go run`


## API del servicio
**Method:** 
> GET

**URL:** 
> {endpoint}/message/

**HEADER:**
> userId -> string
    
 
**SUCCESS RESPONSE - STATUS 200**
>      {
>        "message": "This is Fucking Awesome.",
>        "subtitle": "- userId"
>      }

**ERROR RESPONSE - STATUS 429 (Too Many Requests)**
>      "Try again later ;("
