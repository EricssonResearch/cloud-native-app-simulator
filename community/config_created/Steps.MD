# Hydragen

## Configuración

Primero que todo, hay que tener en cuenta que toca tener lo siguiente:

1. Hay que tener instalado el kind
   1. Kind se puede instalar con el siguiente comando `go install sigs.k8s.io/kind@v0.29.0`, este comando requiere tener instalado GoLang en el sistema
2. Hay que construir el json con las especificaciones de los microservicios
    1. Es importante colocar el objeto `settings: { "development": true  }` para que pueda halar la imagen del docker de la que se va a construir localmente
   2. Cuando se crea un cluster con _kind_, el mismo le coloca el prefijo de kind cuando Kubectl quiere referirse a el. (Quiere decir, si el cluster se creó con el nombre 'cluster-1', entonces va a aparecer con el nombre 'kind-cluster-1')
3. Correr el bash `generate.sh preset {inputfile}`, este va a crear cada uno de los servidores y las especificaciones de k8s para poder desplegar los microservicios en K8S, importante correr primero `newgrp docker` para que el docker no moleste por los permisos. 
4. Hay que desplegar o cargar las imagenes de docker al cluster de k8s. Para eso los scripts de `kind-push-image-to-clusters.sh` funcionan bastante bien, pero como solo se esta trabajando con un cluster, entonces solamente es necesario correr el comando una vez, con lo que se muestra en el bash.
5. Correr el bash `deploy.sh {inputfile}`, si el `logging = true` en el archivo del input, entonces va a tratar de desplegar algunos pods adicionales para elasticsearch, sin embargo, no pude verificar bien esta opción porque el SW doce forbidden por una extraña razón.
   1. Este bash va a desplegar el Deployment, Service y Config en el cluster de k8s.

## Como correrlo y monitorearlo

En el setup inicial hay que tener en cuenta varias cosas:

### Helm

Helm es un package manager para k8s, con este manejador es con lo que vamos a instalar par de las cosas que vamos a ver en esta sección. https://helm.sh/

### kubernetes-dashboard

https://github.com/kubernetes/dashboard

Esta es una herramienta UI para poder verificar cada uno de los recursos disponibles en el cluster de k8s. Esta herramienta se puede instalar con Helm. Sin embargo, para hacerla funcionar hay que referirse primero a crear un usuario para este dashboard

#### Crear un usuario para el dashboard

https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md

Siguiendo este tutorial se puede ver como se puede realizar, sin embargo, puede referirse a los archivos de:
- `create-role-binding-dashboard.yaml`
- `create-service-account-dashboard.yaml`

El crear un service account es como crear un usuario para poder acceder a la lectura de los recursos que tenemos activos. Los permisos estan como admin (Solo para desarrollo)

Cuando tengamos el usuario podemos correr el comando: `kubectl -n kubernetes-dashboard create token admin-user` donde el `admin-user` es el nombre del usuario definido en los archivos previamente mencionados.
Con el JWT que se genera, podemos acceder al dashboard

Con el comando `kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443`
Podemos acceder al dashboard por medio del localhost

## Pegarle a los endpoints definidos en Hydragen

Hay dos opciones para esto:
- Redirigir el puerto al servicio directamente
- Modificar un poco los ingress para poder hacerlo mediante los ingress.

Sin embargo, para los dos puntos hay que tener claro una cosa:

```yaml
apiVersion: v1
kind: Service
metadata:
    name: service1
    ...
spec:
    selector:
        app: service1
    ports:
        - name: http
          port: 80
          targetPort: 5000
```

Las ultimas tres lineas de la definición del servicio son muy importantes.
Note que el puerto que esta siendo expuesto por el **SERVICIO** es el puerto $80$, mientras que el `targetPort` es el puerto el cual tiene que estar escuchando el **contenedor** o el puerto que expone el contenedor.

Por lo tanto, siempre nos vamos a referir al puerto del servicio

### Redirigir el puerto al servicio directamente

Podemos redirigir el puerto que nosotos necesitemos de la siguiente manera: `kubectl port-forward svc/service1 8080:<Service Port>`, en este caso, el Service Port es el $80$.

De esta manera, si se realiza una llamada al `localhost:8080`, va a estar redirigiendo al puerto correctamente.

### Hacerlo por medio del Ingress

Para hacerlo por medio del Ingress, hay que tener en cuenta que cada ingress debe tener una implementación definida. Esta usualmente es la de Nginx la cual se puede instalar por helm

Esta implementación es conocida como el ingress-controller y necesita ser instalada antes de realizar cualquier definición de un Ingress. https://kubernetes.github.io/ingress-nginx/deploy/#quick-start

Una vez instalado el ingress-controller, este va a levantar dos servicios: 'admission' y el controller como tal. En este caso, como estamos en desarrollo local, el controller como tal se va a quedar esperando un proveedor externo a que le de una dirección IP, pero nosotros podemos modificar esto.

Para que nuestro ingress funcione, tenemos que modificarlo con el comando de: `kubectl edit svc ingress-nginx-controller -n ingress-nginx`, este va a editar el servicio y lo que tenemos que modificar es que este no sea de tipo `Load Balancer` si no que sea de tipo `NodePort`.

**En el caso especifico de Kind**, Kind no expone los puertos de manera automatica, entonces hay que exponerlo explicitamente de la siguiente forma: `kubectl port-forward -n ingress-nginx service/ingress-nginx-controller 8080:80` de esta manera, ahora el ingress esta tomando el puerto $8080$ que esta escuchando desde afuera.

Una vez expuesto el puerto, vamos a crear el ingress que queremos para que pueda redirigirlo al puerto. Este archivo se encuentra en en `ingress.yaml` y note que hay una propiedad que especifica que la clase es de `nginx` para que este controlador sea el que maneje el ingress. De esta manera, en `localhost:8080` va a entrar al ingress.

En este caso, es necesario que el puerto sea el 80, ya que es el que expone el servicio del nginx por defecto.
El flujo del ingress es el siguiente:
```
[navegador o curl] → [Ingress Controller (NGINX)] → [tu Service] → [tu Pod]
```

Cabe aclarar que hacer todo esto con minkube es mucho mas facil.
