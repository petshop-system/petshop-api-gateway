# Petshop-api-gateway

Este projeto visa abstratir a comunicação entre o usuário e as diversas APIs que 
compõem a solução do projeto Petshop-System.

## Ferramentas

* Dockerfile
* Docker-compose
* Go
* Makefile

## Como inicializar

Para inicializar a aplicação pode-se fazer uso de comandos pré-configurados no Makefile, desta forma, as aplicações
que fazem parte da solução Petshop-System irão inicializar junto com Petshop-api-gateway.

Ex.: `$ make docker-compose-up `

## Como executar requests

Para executar um request a API Gateway, depois do host deve-se conter o path inicial que, neste caso, deve ser
_/petshop-system_.

Ex.:

```
 % curl -vL --location 'http://localhost:9999/petshop-system/address/search/1?test=testando'
* processing: http://localhost:9999/petshop-system/address/search/1?test=testando
*   Trying 127.0.0.1:9999...
* Connected to localhost (127.0.0.1) port 9999
> GET /petshop-system/address/search/1?test=testando HTTP/1.1
> Host: localhost:9999
> User-Agent: curl/8.2.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Length: 146
< Content-Type: application/json
< Date: Tue, 26 Dec 2023 17:13:02 GMT
< Request_id: 7a76e1c5-0cf6-4ecc-8c53-3c6b49423ab3.1703610782046515387
< 
{"message":"address found with success","result":{"id":1,"street":"Rua Jose Bonifácio","number":"1432"},"date":"2023-12-26T17:13:02.051459512Z"}
```

## Configurar redicionamentos

Para configurar um novo host deve-se configurar a rota de redirecionamento no banco de dados da aplicação,
na tabela de rotas.

Ex.:

### Coluna

* router: informa, de acordo com o path do request, qual configuração será utilizada;
* configuration: configuração para realizar o redirecionamento do request:
  * host: informa qual host será utilizado no redirecionamento do request;
  * replace-old-app-context: informa qual parte do path será substituído no request;
  * replace-new-app-context: informa pelo que será substiuído o path no request;

Ex.: 

**Address**

```
'address', // router
'{ // configuration
    "host": "http://petshop-api:5001", 
    "replace-old-app-context": "petshop-system", 
    "replace-new-app-context": "petshop-api"
}'
```

**Employee**

```
'employee',  // router
'{ // configuration
    "host": "http://petshop-admin-api:5002", 
    "replace-old-app-context": "petshop-system", 
    "replace-new-app-context": "petshop-admin-api"
}'
```

**BFF Desktop Service**

````
'bff-desktop-service', // router
'{ // configuration
    "host": "http://petshop-bff-desktop:9998", 
    "replace-old-app-context": "petshop-system/bff-desktop-service", 
    "replace-new-app-context": "petshop-bff-desktop"
}'
````

## Material utilizado

* [O que é API Gateway?](https://www.iugu.com/blog/api-gateway)
* [How to Create a Reverse Proxy using Golang](https://www.codedodle.com/go-reverse-proxy-example.html)
* [Go Simple and powerful reverse proxy](https://www.sobyte.net/post/2022-04/golang-reverse-proxy/)
* [Why is httputil.NewSingleHostReverseProxy causing an error on some www sites?](https://stackoverflow.com/questions/31715545/why-is-httputil-newsinglehostreverseproxy-causing-an-error-on-some-www-sites)