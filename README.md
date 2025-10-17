# Shortener URL API
O projeto consiste em reproduzir um encurtador de URL's simples.

## Decisões Tecnicas e Arquiteturais

### Linguagem utilizada - Go
Este projeto foi desenvolvido na linguagem Go devido ao seu meio termo entre simplicidade de codigo, eficiencia e desempenho.

### Docker
A utilizacao do Docker permite que consigamos executar a aplicacao em qualquer ambiente que possua o Docker instalado, garantindo que nao haja nenhuma inconsistencia ou incompatibilidade entre maquinas e sistemas operacionais.

Neste app, utilizamos a imagem customizada em `Dockerfile` para a construcao do container do servidor, utilizando uma imagem builder `alpine` para otimizar o tamanho da imagem. Como banco de dados, utilizamos a imagem `mysql:5.7`, que e utilizada dentro do servidor para persistencia de dados. E por fim tambem temos uma imagem customizada para executar testes automatizados do app, localizada em `Dockerfile.tests`

### Algoritmo de geracao de Alias
Foi utilizado o algoritmo `datetime + random + base62` para gerar novos alias. Seu uso foi escolhido pelo meio termo entre simplicidade na implementacao e sua unicidade para evitar colisoes com outros registros, pois sua formula faz com que seja quase impossivel de haver colisao, mesmo aumentando a escalabilidade do app.

## Casos de uso

### Criacao de URL encurtada
![diagrama de criacao de URL encurtada](/docs/img/create_case_diagram.png)

Endpoint: POST /?url=[url]&alias=[alias]
Parametros query:
* url - obrigatorio
* alias - opcional (se nao enviar, um alias aleatorio e gerado durante o cadastro)

Exemplo de resposta:
![exemplo de criacao de URL encurtada](/docs/img/create_response_example.png)


OBS: Ao fazer uma request para a URL retornada, e possivel obter a URL completa

### Obtencao de URL real utilizando o alias
![diagrama de Obtencao de URL real utilizando o alias](/docs/img/retrieve_by_alias_case_diagram.png)

Endpoint: GET /u/{alias}
Parametros URL:
* alias - obrigatorio

Exemplo de resposta:
![exemplo de resposta da Obtencao de URL real utilizando o alias](/docs/img/retrieve_by_alias_response_example.png)

### Obtencao das 10 URL mais acessadas
![diagrama de Obtencao de URL real utilizando o alias](/docs/img/retrieve_by_alias_case_diagram.png)

Endpoint: GET /most_acessed

Exemplo de resposta:
![exemplo de Obtencao das 10 URL mais acessadas](/docs/img/retrieve_10_most_accessed_urls_response_example.png)


## Instucoes para executar o app
1. Certifique-se de ter o Docker e docker-compose instalados em sua maquina.
2. Execute este comando no terminal para o build dos containers e a execucao do app:
```shell
docker-compose up --build -d app
```

3. Apos a execucao do comando, tanto o container do `MySQL` quanto do APP irao inicializar.
4. Observe os logs do container do app com o seguinte comando para verificar se o server comecou a executar:
```shell
docker logs shortener_url_app
```

OBS: E necessario visualizar um log com a seguinte mensagem `Server is running at port 8080`. Feito isso, o app ja pode ser utilizado.

## Instrucoes para executar os testes automatizados
1. Execute o comando no terminal para o build do container de teste:
```shell
docker-compose up --build -d tests
```
2. Execute o comando para entrar no terminal do container de testes:
```shell
docker-compose exec tests sh
```

3. Para executar os testes automatizados, execute o seguinte comando:
```shell
go test ./...
```

### Instrucoes para parar e remover containers e imagens criados
Execute o comando abaixo:
```shell
docker-compose down --rmi all
```