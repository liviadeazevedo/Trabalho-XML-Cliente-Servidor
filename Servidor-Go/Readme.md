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
    $ go get github.com/lestrrat-go/libxml2
    $ go get gopkg.in/xmlpath.v1
    ```
    Links para os projetos utilizados:
    
    - https://github.com/lestrrat-go/libxml2
    - https://gopkg.in/xmlpath.v1
    
    O primeiro _go get_ executado ira criar um diretório denominado _go_, localizado em _/home/<nome-usuario>/go_.
    
4) Configurar o ambiente de execução do projeto:
    - Entre na pasta _/home/\<nome-usuario>/go/src_ e clone este projeto inteiro(utilizando o comando _git clone_)
	- abra o terminal em Servidor-Go e execute o bash:
        ```sh
        $ ./build.sh
        ```
        Isso irá permitir que os dois pacotes possam ser usados em qualquer programa que os importe.

Por fim, para usar todas as funcionalidades disponíveis nesses pacotes, basta importá-los para seu programa:

```go
package main

import (
    (...outros imports...)
    "Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLog"
    "Trabalho-XML-Cliente-Servidor/Servidor-Go/serverLogic"
    "Trabalho-XML-Cliente-Servidor/Servidor-Go/serverConnection"
    (...outros imports...)
)

(...código...)

serverLogic.<nome-do-metodo>(...)
serverConnection.<nome-do-metodo>(...)
serverLog.<nome-do-metodo>(...)
```

Para executar seu programa em Go(lembrando que ele  deve pertencer ao pacote _main_), execute o seguinte comando:
```sh
$ go run <nome-arquivo>.go
```
