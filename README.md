# otel-zipkin

## Instruções para rodar o prjeto
### 1 - Na raiz do projeto execute o comendo docker compose up -d
### 2 - Faça uma requisição ao serviço A (que roda na porta 8081). Exemplo:
#### curl -XPOST -H 'Content-Type: application/json' -d '{"cep":"71218010"}' 'http://localhost:8081/weather'
### 3 - Acesse a url abaixo e pesquise pelo trace da requisição. Nele você encontrará os tepos de resposta do serviço A (busca de CEP) e do serviço B (busca de temperatura).
#### http://localhost:9411/zipkin
