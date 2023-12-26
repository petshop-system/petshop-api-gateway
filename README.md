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

Para configurar um novo host deve-se configurar a rota de redirecionamento 
no arquivo _/configuration/router/router.json_.

Ex.:

```
{
  "customer": {
    "host": "http://petshop-api:5001",
    "app-context": "petshop-api"
  },
  "address": {
    "host": "http://petshop-api:5001",
    "app-context": "petshop-api"
  },
  "employee": {
    "host": "petshop-admin:5002",
    "app-context": "petshop-admin"
  }
}
```

## Material utilizado

* [O que é API Gateway?](https://www.iugu.com/blog/api-gateway)
* [How to Create a Reverse Proxy using Golang](https://www.codedodle.com/go-reverse-proxy-example.html)
* [Go Simple and powerful reverse proxy](https://www.sobyte.net/post/2022-04/golang-reverse-proxy/)
* [Why is httputil.NewSingleHostReverseProxy causing an error on some www sites?](https://stackoverflow.com/questions/31715545/why-is-httputil-newsinglehostreverseproxy-causing-an-error-on-some-www-sites)