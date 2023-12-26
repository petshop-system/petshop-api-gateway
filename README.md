# Petshop-api-gateway

Este projeto visa abstratir a comunicação entre o usuário e as diversas APIs que 
compõem a solução do projeto Petshop-System.

## Ferramentas

* Dockerfile
* Docker-compose
* Go
* Makefile

## Como iniciar

Para inicializar a aplicação pode-se fazer uso de comandos pré-configurados no Makefile, desta forma, as aplicações
que fazem parte da solução Petshop-System irão inicializar junto com Petshop-api-gateway.

Ex.: `$ make docker-compose-up `