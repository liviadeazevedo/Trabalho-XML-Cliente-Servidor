import locale
from io import StringIO

from lxml import etree

from sys import exit
from socket import *
from threading import *
import time

TESTE = False

###########

received_msg = ''
lock = Lock()

class ClientSocket(Thread):
    def __init__(self, sock=None):
        super().__init__()

        if sock is None:
            self.sock = socket(AF_INET, SOCK_STREAM)
        else:
            self.sock = sock

        self.pHdr_len = 5
        self.recvBuffer = b''
        self.onThread = True

        addrsFileLines = open("addrs.txt", "r").readlines()
        self.host, self.port = addrsFileLines[0], int(addrsFileLines[1])

    def run(self):
        # Método que implementa o que a Thread roda
        # global on
        # global received_msg
        #
        while self.onThread:
            self._receive()

    def defineAddrs(self):
        addrsFile = open("addrs.txt", "w")
        self.host = input("host:")
        self.port = input("port:")
        addrsFile.write(self.host + "\n" + str(self.port))

    def connect(self, host=None, port=None):
        if host is not None:
            self.host = host
        if port is not None:
            self.port = port
        try:
            self.sock.connect((host, port))
            return True

        except Exception as e:
            return False

    def send(self, msg):
        ''' Método para envio de mensagens'''
        MSGLEN = len(msg)
        totalsent = 0
        while totalsent < MSGLEN:
            sent = self.sock.send(msg[totalsent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            totalsent = totalsent + sent

        self.start()

    def _receive(self):
        bytes_recv = 0

        try:
            data = self.sock.recv(4096)
        except BlockingIOError:
            pass
        else:
            if data:
                self.recvBuffer += data
                bytes_recv = len(data)
            else:
                raise RuntimeError("Peer closed.")

    def _read_protoheader(self):
        if len(self.recvBuffer) >= self.pHdr_len:
            hdr_len = self.recvBuffer[:self.pHdr_len]
            self._recv_buffer = self._recv_buffer[self.pHdr_len:]
            return int(hdr_len)

    def _read_header(self, hdr_len):
        header = 0
        if len(self._recv_buffer) >= hdr_len:
            header = self.recvBuffer[:hdr_len]
            self._recv_buffer = self._recv_buffer[hdr_len:]
            # Interessante caso o cabeçalho tenha mais de um item
            # for reqhdr in ("byteorder","content-length","content-type","content-encoding",):
            #     if reqhdr not in self.jsonheader:
            #         raise ValueError(f'Missing required header "{reqhdr}".')
        return int(header)

    def _read_msg(self, msg_len):
        msg = None
        if len(self._recv_buffer) >= msg_len:
            msg = self.recvBuffer[:msg_len]
            self._recv_buffer = self._recv_buffer[msg_len:]
        return msg

    def read(self):
        hdr_len = self._read_protoheader()

        header = self._read_header(hdr_len)

        return self._read_msg(header)

    def pscan(self, port):
        ''' Scaner de port'''
        try:
            self.sock.connect((self.host, port))
            return True
        except Exception as e:
            return False

class Candidato():
    def __init__(self, cpf):
        self.cpf = cpf
        self.boletim = Boletim(cpf)

    def submeter(self):
        pass

    def consultaStatus(self):
        pass

class Boletim():
    def __init__(self, cpf):
        self.ctrlXML = ControladorXML()
        self.xml_boletim = self.ctrlXML.lerXML(cpf)

class ControladorXML():
    # Caso seja necessário mais alguma funcionalidade que lide com XML/XSD implementar nesta classe
    def lerXML(self, nome_arq):
        # Abrindo o arquivo xml
        xml_arq = open(nome_arq, "r+")

        # Transformando o arquivo aberto em arvore de elementos
        return etree.parse(xml_arq)

    def lerXSD(self, nome_arq):
        # Abrindo o arquivo xsd
        xsd_arq = open(nome_arq, "r+")

        # Transformando o arquivo aberto em arvore de elementos
        xsd_doc = etree.parse(xsd_arq)

        # Método da biblioteca que Guarda que o arquivo lido é um XML Schema que vai ser usado para validação
        self.xsd =  etree.XMLSchema(xsd_doc)

    def gravarXML(self):
        pass

    def validarXML(self, xml):
        if self.xsd.validate(xml):
            return True
        else:
            return False

class FronteiraInter():
    def entrarCandidato(self):
        return

def main():
    ''' Função para execução no terminal'''

    # no caso de implementação do menu isso deverá ser alterado
    #candidato = Candidato(input("digite o cpf do candidato a entrar:"))

    # client = ClientSocket()
    # client.connect()

    # possível fazer uma menu do tipo:
    # 1 - submeter boletim
    # 2 - consultar status
    # 3 - entrar como outro candidato
    # 0 - sair do programa

    # perguntar se candidato gostaria de submeter boletim
    #candidato.submeter()

    # perguntar se candidato gostaria de consultar status
    # candidato.consultaStatus()

    pass


# Função comentada para consulta de operações necessárias
'''def run():

    global received_msg

    waiting_HE = False

    ctrl_XML = ControladorXML()

    xsd = ctrl_XML.lerXSD("he_schema.xsd")

    c = MySocket()

    c.connect(host, port)

    while waiting_HE:

        with lock:
            if received_msg != '':

                if TESTE:
                    print("Olha: " + received_msg)

                if validate(received_msg, xsd):
                    imprimir(received_msg)
                else:

                    print("Algo de errado não está certo, XML não corresponde ao Schema")


                received_msg = ''
                waiting_HE = False
'''