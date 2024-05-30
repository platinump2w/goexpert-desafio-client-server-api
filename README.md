# Desafio 1 - Client-Server-API
Desafio para aplicar conceitos de webserver HTTP, contextos, banco de dados e manipulação de arquivos usando a linguagem Go.

## Descrição
O objetivo deste desafio é criar um client e um server em Go que se comunicam através de uma API de câmbio de moedas. O client deve solicitar a cotação do dólar ao server, que por sua vez deve consumir a API de câmbio e retornar o valor atual do câmbio. O servidor deve salvar a cotação em um banco de dados e o client deve salvar a cotação em um arquivo cotacao.txt.

## Requisitos Gerais
- Utilizar 2 arquivos: `client.go` e `server.go`;

### Client
- Realizar uma requisição HTTP no `server.go` solicitando a cotação do dólar;
- Receber do `server.go` o valor atual do câmbio (campo `bid` do JSON);
- Salvar a cotação atual em um arquivo cotacao.txt no formato: `Dólar: {valor}`;
- Usar o package `context` para definir um timeout máximo de `300ms` para receber o resultado do `server.go`;
- Registrar erros nos logs caso o tempo de execução seja insuficiente.

### Server
- Consumir a API de câmbio Dólar/Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL;
- Retornar o resultado no formato JSON para o `client.go`;
- Usar o package `context` para:
  - Definir um timeout máximo de `200ms` para chamar a API de cotação do dólar;
  - Registrar no banco de dados SQLite cada cotação recebida, com um timeout máximo de `10ms`;
  - Registrar erros nos logs caso o tempo de execução seja insuficiente;
- Criar um endpoint `/cotacao` na porta `8080`.

## Execução
- Acessar a pasta `server` e executar o comando `go run server.go`;
- Acessar a pasta `client` e executar o comando `go run client.go`;

## Observações
- O arquivo `cotacao.txt` será criado na pasta `client` e será atualizado a cada execução do `client.go`;
- O banco de dados SQLite será criado na pasta `server` e será chamado de `exchange_rates.db`;
