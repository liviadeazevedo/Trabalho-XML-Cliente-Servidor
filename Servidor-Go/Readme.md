# Servidor - Linguagem Go

Para executar o servidor, algumas etapas precisam ser feitas(considerando um ambiente Linux-Ubuntu):

 1) Instale o Golang(versão 1.12):
    ```sh
    $ sudo add-apt-repository ppa:longsleep/golang-backports
    $ sudo apt-get update
    $ sudo apt-get install golang-go
    ```
    
 2) Instale a dependência "libxml2-dev" para futuramente usar os pacotes em Go para XML:
    ```sh
    $ sudo apt-get install libxml2-dev 
    ```

 3) Instale os pacotes Go necessários para o projeto:
    ```sh
    $ go get github.com/gorilla/websocket
    $ go get github.com/lestrrat-go/libxml2
    $ go get gopkg.in/xmlpath.v1
    ```
    Links para os projetos utilizados:
    
    - https://github.com/lestrrat-go/libxml2
    - https://github.com/gorilla/websocket
    - https://gopkg.in/xmlpath.v1
    
    O primeiro _go get_ executado ira criar um diretório denominado _go_, localizado em _/home/<nome-usuario>/go_.
    
4) Configurar o ambiente de execução do projeto:
    - Crie uma pasta chamada "servidorGoXML" no diretório _/home/<nome-usuario>/go/src/_ ;
    - Copie as pastas "serverLogic" e "serverConnection" do repositório para dentro da pasta recém-criada("servidorGoXML");
    - Entre em cada uma das pastas, e execute o comando:
        ```sh
        $ go build
        ```
        Se não houver algum erro, execute:
        ```sh
        $ go install
        ```
        Isso irá permitir que os dois pacotes podem ser usados em qualquer programa que os importe.
        
**Obs:** Todas as alterações nesses arquivos deverão ser feitas diretamente nos arquivos dentro desses diretórios na pasta do _go_. Caso haja alterações, substitua os arquivos antigos pelos novos com as alterações. Da mesma forma, para commitar no repositório, registre as suas alterações nos arquivos deste repositório e depois commite.

Por fim, para usar todas as funcionalidades disponíveis nesses pacotes, basta importá-los para seu programa:

```go
package main

import (
    (...outros imports...)
    "servidorGoXML/serverLogic"
	"servidorGoXML/serverConnection"
	(...outros imports...)
)

(...código...)

serverLogic.<nome-do-metodo>(...)
serverConnection.<nome-do-metodo>(...)
```

Para executar seu programa em Go(lembrando que ele  deve pertencer ao pacote _main_), execute o seguinte comando:
```sh
$ go run <nome-arquivo>.go
```