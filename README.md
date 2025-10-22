APIs del sexto cuatrimestre

API_RASPBERRY:
Recibe peticiones ya sea de /atraccion o /visitas las cuales son almacenados en una base de datos
de MySQL para luego enviarse a un broker de RabbitMQ
API2_ZOO:
Consume los datos enviados por el broker para insertarlos en su propia base de datos la cual servira
para hacer multiples consultas del front
API3_ZOO:
API basica de CRUD de usuarios

Como iniciarlo:
 - Despues de clonar el repositorio acceda a cada api de manera individual, osea abra una nueva terminal para cada
api, una para API_RASPBERRY, API2_ZOO y API3_ZOO
 - Una vez estando en cada terminal ejecute - go mod tidy - para instalar los paquetes necesarios de cada api
 - Ahora prepare los servicios que usara cada api estos fueron generados en docker por lo que se le proporcionara los comandos necesarios para levantarlos:
    MySQL: La primera base de datos para ahorrar recursos sera usada por la API_RASPBERRY y la API3_ZOO:
    - [ docker run -d --name mysql-container1 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=mydb -p 3306:3306 mysql:8.0 ]
    - Para la segunda base de datos sera usada unicamente por la API2_ZOO:
    - [ docker run -d --name mysql-container2 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=data -p 3307:3306 mysql:8.0 ]
    RabbitMQ: Este servicio sera usado por la API_RASPBERRY y la API2_ZOO:
    - [ docker run -d --name rabbitmq -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=password -p 5672:5672 -p 15672:15672 -v rabbitmq_data:/var/lib/rabbitmq --restart unless-stopped rabbitmq:3-management ]
Nota: Se puede generalizar todo en una sola base de datos por lo que si en verdad quiere ahorrarse hacer multiples bases de datos 
puede modifiar la api para usar una sola base de datos, sin embargo se recomienda tenerlo en diferentes para tener una idea del flujo de datos que se tenia propuesto

Una vez teniendo todos los contenedores levantados y las apis listas puede proceder a iniciar cada una con [ go run main.go ], este comando sirve para cada api por lo que solo requiere copiar y pegar en cada terminal para luego ejecutarlo, se recomienda inciar desde
la API_RASPBERRY, API2_ZOO y API3_ZOO en ese orden por seguridad, una vez iniciado deberia de ver una ruta para poder acceder a la documentacion de cada api y disponer los endpoints disponibles para su uso.
Lo primero que se debe de hacer es crear un nuevo usuario para poder acceder a los metodos de las apis.
PARA PODER USAR LOS METODOS SE NECESITARA UN TOKEN EL CUAL UNA VEZ OBTENIDO SE INSERTARA EN EL APARTADO DE AUTHORIZE MAS ESPECIFICO EL BOTON VERDE QUE APARECERA ABAJO DE LA DESCRIPCION DE SWAGGER CON UN CANDADO INSERTAR EL TOKEN EN EL APARTADO VALUE Y PULSAR EL BOTON DE AUTHORIZE UNA VEZ HECHO ESO DEBERIA DE PODER CERRAR LA VENTANA SIN PROBLEMAS Y USAR LOS METODOS DE AHORA EN ADELANTE.
Muchas gracias.