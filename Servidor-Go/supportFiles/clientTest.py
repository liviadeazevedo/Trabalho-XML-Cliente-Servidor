# TCP Client Code

#arg[1] = opção 'file' ou 'local'
#arg[2] = se opção for 'file', dizer o nome do arquivo a ser lido(diretório referência já
# para a pasta 'Arquivos').
#arg[3] = se opção for '1', faz o protocolo 1, se for '2', faz o protocolo 2.

import os
import sys

host="127.0.0.1"            # Set the server address to variable host

port=4446               # Sets the variable port to 4444
protocol=sys.argv[3]

from socket import *             # Imports socket module

s=socket(AF_INET, SOCK_STREAM)      # Creates a socket

s.connect((host,port))          # Connect to server address

#msg=s.recv(1024)            # Receives data upto 1024 bytes and stores in variables msg
s.send(protocol.encode('utf-8'))

test_src = ""
taml5b = ""

ctrl = input("começar a enviar?")

def protocol1Communication(msg, tamStr):
    taml5b = (("0"*5)+tamStr)[-5:]
    
    s.send(taml5b.encode('utf-8'))

    print("tamanho enviado\nenviando mensagem\n\n", msg)
    s.send(msg.encode('utf-8'))
    
    print("\n\nmensagem enviada\nesperando tamanho")
    size = int(str(s.recv(5).decode('utf-8')))
    print("\ntamnho recebido\nesperando recebimento da mensagem de tamanho " + str(size))
    return str(s.recv(size))

def protocol2Communication(msg, tamStr):
    tamCabecalho = str(len(tamStr))
    tamCabecalho2b = ("0"+tamCabecalho)[-2:]

    print("\ntamCabecalho2b ", tamCabecalho2b, "\n tamMsg ", tamStr)
    s.send(tamCabecalho2b.encode('utf-8'))
    s.send(tamStr.encode('utf-8'))
    print("tamanho enviado\nenviando mensagem")
    s.send(msg.encode('utf-8'))

    
    print("\n\nmensagem enviada\nesperando tamanhos")
    sizeHeader = int(str((s.recv(2)).decode("utf-8")))
    print("\ntamnho recebido\nesperando recebimento do cabecalho " + str(sizeHeader))
    size = int(str((s.recv(sizeHeader)).decode("utf-8")))
    print("\ntamnho recebido\nesperando recebimento da mensagem " + str(size))
    return str(s.recv(size))

while(ctrl == ''):
    
    if sys.argv[1] == "file":
        f = open('../../Arquivos/' + sys.argv[2],'r')
        test_src = str(f.read(os.path.getsize(f.name)))
        tamStr = str(os.path.getsize(f.name)) 
    else:
        test_src = "<resposta><retorno>0</retorno></resposta>"
        tamStr = str(len(test_src))
    
    if protocol == "2":
        msg = protocol2Communication(test_src, tamStr)
    else:
        msg = protocol1Communication(test_src, tamStr)

    
    print(msg)
    ctrl = input("proximo envio?")
    

#s.send(bytes("out", 'utf-8'))
s.close()                            # Closes the socket 
# End of code