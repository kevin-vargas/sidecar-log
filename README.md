# Objetivos

Este sidecar nos permite obtener el contenido por stdout de una aplicación por medio de la interaccion con la api de kubernetes, y lo envia por tramas a un pub sub para poder procesar estos logs de forma asincronica, se selecciono un mqtt broker debido a la infraestractura en cual se va a implementar.

# Development
Para poder ejecutar el sidecar de forma local es necesario generar el un archivo `.env` para las variables de entorno del proceso.
```
SCOPE=dev
HOSTNAME=libre-job-56d8df4d66-swfb7
```
Donde `hostname` es el pod al cual contiene el container con la aplicacion a extender la funcionalidad.

Tambien es necesario un archivo de config en ${HOME}/.kube/config con las credenciales necesarias para poder acceder a las funcionalidades expuestas por la api de kubernetes

Para poder uso de forma local de la funcionalidad pub y sub podemos ejecutar un container para este fin.
```
docker run -d --name emqx -p 18083:18083 -p 1883:1883 emqx/emqx:latest
```