# TCP Client Code

#arg[1] = opção 'file' ou 'local'
#arg[2] = se opção for 'file', dizer o nome do arquivo a ser lido(diretório referência já
# para a pasta 'Arquivos').

import os
import sys

host="127.0.0.1"            # Set the server address to variable host

port=4446               # Sets the variable port to 4444

from socket import *             # Imports socket module

s=socket(AF_INET, SOCK_STREAM)      # Creates a socket

s.connect((host,port))          # Connect to server address

#msg=s.recv(1024)            # Receives data upto 1024 bytes and stores in variables msg

test_src = ""
taml5b = ""

ctrl = input("começar a enviar?")

while(ctrl == ''):
    
    if sys.argv[1] == "file":
        f = open('../../Arquivos/' + sys.argv[2],'r')
        test_src = str(f.read(os.path.getsize(f.name)))
        taml5b = (("0"*5)+str(os.path.getsize(f.name)))[-5:]
    else:
        test_src = "<resposta><retorno>0</retorno></resposta>"
        taml5b = (("0"*5)+str(len(test_src)))[-5:]
    
    s.send(taml5b.encode('utf-8'))

    print("tamanho enviado\nenviando mensagem")
    s.send(test_src.encode('utf-8'))
    
    print("\n\nmensagem enviada\nesperando tamanho")
    size = int(str((s.recv(5)).decode("utf-8")))
    print("\ntamnho recebido\nesperando recebimento da mensagem de tamanho " + str(size))
    msg = str(s.recv(size))
    
    print(msg)
    ctrl = input("proximo envio?")
    

#s.send(bytes("out", 'utf-8'))
s.close()                            # Closes the socket 
# End of code