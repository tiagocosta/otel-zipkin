# otel-zipkin

## Instruções para rodar o prjeto
### 1 - Na raiz do projeto execute o comendo docker compose up -d
### 2 - Acesse o container do serviço A e do serviço B executando os comandos: 
#### docker exec -it service-a bash
#### docker exec -it service-b bash
### 3 - Dentro de cada container navegue até a pasta src/app e execute o comando: 
#### go run cmd/main.go
### 4 - Faça uma requisição ao serviço A (que roda na porta 8081)
#### curl -XPOST -H 'Content-Type: application/json' -d '{"cep":"71218010"}' 'http://localhost:8081/weather'
### 5 - Acesse a url abaixo e pesquise pelo trace da requisição. Nele você encontrará os tepos de resposta do serviço A (busca de CEP) e do serviço B (busca de temperatura).
#### http://localhost:9411/zipkin
