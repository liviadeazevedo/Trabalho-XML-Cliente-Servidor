import locale
import os
from io import StringIO

from lxml import etree

from sys import exit
from socket import *
from threading import *

import time

def_cod = 'utf-8'
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

        self._hdr_num = 2 # numero de cabeçalhos
        self.pHdr_len = 2 # tamanho do proto cabeçalho
        self.recv_msg = b'' # msg recebida
        # self.onThread = True # indicador para desligar Thread

        try:
            addrsFileLines = open("addrs.txt", "r").readlines()
        except Exception as e:
            print(e)
        else:
            self.host, self.port = addrsFileLines[0][:-1], int(addrsFileLines[1]) # host e port padrôes lidos do arquivo


    def run(self):
        # Método que implementa o que a Thread roda
        self.recv_msg = self._receive()

    def defineAddrs(self):
        addrsFile = open("addrs.txt", "w")
        self.host = input("host:")
        self.port = int(input("port:"))
        addrsFile.write(self.host + "\n" + str(self.port))

    def connect(self, host=None, port=None):
        if host is not None:
            self.host = host
        if port is not None:
            self.port = port
        try:
            self.sock.connect((self.host, self.port))
            return True

        except Exception as e:
            print(e)
            return False

    def proto_define_connect(self):
        if self.connect():
            self.sock.send(str(self._hdr_num).encode("utf-8"))
            return True
        else:
            return False

    def send(self, msg):
        ''' Método para envio de mensagens'''
        msg_cod = msg.encode(def_cod)
        HDRLEN = str(len(msg_cod)) # tamanho do cabeçalho
        prt_hdr = str(len(HDRLEN)) # tamanho do proto cabeçalho
        if len(prt_hdr) < self.pHdr_len: # caso o proto cabeçalho tenha uam quantidade de algarismos menor do que a esperada, será concatenado a quantidade de zeros necessaria na frente
            prt_hdr = '0'*(self.pHdr_len-len(prt_hdr)) + prt_hdr
        elif len(prt_hdr) > self.pHdr_len:
            raise RuntimeError("Proto Cabeçalho maior que", self.pHdr_len,"bytes")
        data = (prt_hdr + HDRLEN).encode(def_cod) + msg_cod # Dado que será efetivamente enviado = (proto cabeçalho + cabeçalho).encode + mensagem codificada
        total_len = len(data)
        totalsent = 0

        while totalsent < total_len:
            try:
                sent = self.sock.send(data[totalsent:])
            except Exception as e:
                exit(1)

            else:
                if sent == 0:
                    raise RuntimeError("socket connection broken")

                totalsent += sent

    def _receive(self):
        ''' Método que roda na thread para receber os dados que são recebidos da rede e concatenar no self.recvBuffer'''
        hdr_len = None
        header = None
        msg = None

        while hdr_len is None:
            hdr_len = self._read_protoheader()

        if hdr_len is not None:
            if header is None:
                header = self._read_header(hdr_len)

        if header:
            if msg is None:
                msg = self._read_msg(header)
        else:
            print("Sem mensgagem para ler")

        return msg

    def _read_protoheader(self):
        hdr_len = self.sock.recv(self.pHdr_len)
        return int(hdr_len)

    def _read_header(self, hdr_len):
        hdr = self.sock.recv(hdr_len)
        # Interessante caso o cabeçalho tenha mais de um item
        # for reqhdr in ("byteorder","content-length","content-type","content-encoding",):
        #     if reqhdr not in self.jsonheader:
        #         raise ValueError(f'Missing required header "{reqhdr}".')
        return int(hdr)

    def _read_msg(self, msg_len):
        return self.sock.recv(msg_len)

    def read(self):
        self.start()
        while self.recv_msg == b'':
            pass
        msg = self.recv_msg
        return msg.decode(def_cod)

    def pscan(self, port):
        ''' Scaner de port'''
        try:
            self.sock.connect((self.host, port))
            return True
        except Exception as e:
            return False

    def close_con(self):
        # self.onThread = False
        self.sock.close()

class ClientSimple(ClientSocket):
    def __init__(self):
        super().__init__()
        self.pHdr_len = 5
        self._hdr_num = 1

    def send(self, msg):
        ''' Método para envio de mensagens'''
        hdr_len = str(len(msg))

        if len(hdr_len) < self.pHdr_len:
            hdr_len = ('0' * self.pHdr_len) + hdr_len

        elif len(hdr_len) > self.pHdr_len:
            raise RuntimeError("Proto Cabeçalho maior que " + str(self.pHdr_len) + " bytes")

        data = (hdr_len[-self.pHdr_len:] + msg).encode("utf-8")
        total_len = len(data)
        totalsent = 0

        while totalsent < total_len:
            try:
                sent = self.sock.send(data[totalsent:])
            except Exception as e:
                print(e)

            else:
                if sent == 0:
                    raise RuntimeError("socket connection broken")

                totalsent += sent

    def read(self):
        msg_len = None
        msg = None

        if msg_len is None:
            msg_len = self._read_protoheader()

        if msg_len is not None:
            if msg is None:
                msg = self._read_msg(msg_len)
        else:
            print("Sem mensgagem para ler")

        return msg

class Candidato():
    def __init__(self, cpf):
        self.cpf = cpf
        self.boletim = 'boletins\\' + self.cpf + '.xml'
        self.menssagem = ClientSocket()
        self.menssagem.proto_define_connect()

    def submeter(self):
        ctrlXML = ControladorXML()

        boletim = ctrlXML.lerXML(self.boletim)
        msg = ctrlXML.criarRequisicao("submeter", {'Boletim': ctrlXML.to_string(boletim).decode(def_cod)})
        self.menssagem.send(ctrlXML.to_string(msg).decode(def_cod))

    def consultaStatus(self):
        pass

class ControladorXML():
    # Caso seja necessário mais alguma funcionalidade que lide com XML/XSD implementar nesta classe
    def lerXML(self, nome_arq):
        # Abrindo o arquivo xml
        try:
            xml_arq = open(nome_arq, "r+", encoding=def_cod)
        except Exception as e:
            print(e)
            return None
        # Transformando o arquivo aberto em arvore de elementos
        return etree.parse(StringIO(xml_arq.read()))

    def lerXSD(self, nome_arq):
        # Abrindo o arquivo xsd
        xsd_arq = open(nome_arq, "r+", encoding=def_cod)

        # Transformando o arquivo aberto em arvore de elementos
        xsd_doc = etree.parse(xsd_arq)

        # Método da biblioteca que Guarda que o arquivo lido é um XML Schema que vai ser usado para validação
        self.xsd =  etree.XMLSchema(xsd_doc)
        return self.xsd

    def gravarXML(self):
        pass

    def validarXML(self, xml):
        if self.xsd.validate(xml):
            return True
        else:
            return False

    def to_string(self, xml):
        return etree.tostring(xml, encoding=def_cod)

    def criarRequisicao(self, nome_func, dict_param):
        root = etree.Element("requisicao")
        metodo = etree.SubElement(root, "metodo")
        etree.SubElement(metodo, "nome").text = nome_func
        parametros = etree.SubElement(metodo, "parametros")
        for nome, valor in dict_param.items():
            parametro = etree.SubElement(parametros, "parametro")
            etree.SubElement(parametro, "nome").text = nome
            etree.SubElement(parametro, "valor").text = etree.CDATA(valor)

        return root

    def imprimir(XMLdoHistorico):  # parametro: string

        # xml_arq = open(XMLdoHistorico, "r",-1,"utf-8")
        xml_arq = etree.parse(StringIO(XMLdoHistorico))
        xml = xml_arq.getroot()

        print("------------------------------------------------------------------------")
        print(xml.find('universidade').find('nome').text)
        print(xml.find('universidade').find('abreviacao').text)
        print("Curso: " + xml.find('curso').text)
        print("Aluno: " + xml.find('aluno').text)
        print("Matricula: " + xml.find('matricula').text)
        print("Cr medio: " + xml.find('crMedio').text)
        print("Data geracao: " + xml.find('dataGeracao').text)
        print("Hora geracao: " + xml.find('horaGeracao').text)
        print("Cod autenticacao: " + xml.find('codigoAutenticacao').text)
        print("------------------------------------------------------------------------")
        print("\n")
        listaPeriodos = xml.find('periodos').findall('Periodo')
        for i in range(len(listaPeriodos)):
            print("Ano Semestre: " + listaPeriodos[i].find('anoSemestre').text)
            print("Creditos solicitados: " + listaPeriodos[i].find('creditosSolicitados').text)
            print("Creditos acumulados: " + listaPeriodos[i].find('creditosAcumulados').text)
            print("Creditos obtidos: " + listaPeriodos[i].find('creditosObtidos').text)
            print("Cr periodo: " + listaPeriodos[i].find('crPeriodo').text)
            print("\n")

            listaDisciplinasAA = listaPeriodos[i].find('disciplinas').findall('AtividadeAcademica')
            for j in range(len(listaDisciplinasAA)):
                print("\tCodigo disciplina: " + listaDisciplinasAA[j].find('codigo').text)
                print("\tNome disciplina: " + listaDisciplinasAA[j].find('nome').text)
                print("\tCreditos: " + listaDisciplinasAA[j].find('creditos').text)
                print("\tNota: " + listaDisciplinasAA[j].find('nota').text)
                print("\tSituacao: " + listaDisciplinasAA[j].find('situacaoAA').text)
                print("\n")

            listaDisciplinas = listaPeriodos[i].find('disciplinas').findall('Disciplina')
            for j in range(len(listaDisciplinas)):
                print("\tCodigo disciplina: " + listaDisciplinas[j].find('codigo').text)
                print("\tNome disciplina: " + listaDisciplinas[j].find('nome').text)
                print("\tCreditos: " + listaDisciplinas[j].find('creditos').text)
                print("\tNota: " + listaDisciplinas[j].find('nota').text)
                print("\tSituacao: " + listaDisciplinas[j].find('situacao').text)
                print("\n")

            print("------------------------------------------------------------------------")

    def geraHtml(XMLdoHistorico):  # parametro: string

        # xml_arq = open(XMLdoHistorico, "r",-1,"utf-8")
        xml_arq = etree.parse(StringIO(XMLdoHistorico))
        xml = xml_arq.getroot()
        texto = []
        name = "historico"
        name = name + '.html'
        arq = open(name, 'w', -1, "utf-8")

        print('Gerando', name, '...\n')

        texto.append("<!DOCTYPE html>\n<html lang='pt-BR'>\n\n<html>\n\n")
        texto.append(
            "\t<head>\n\t\t<title>Histórico</title>\n\t\t<meta charset = 'utf-8'>\n\t\t<link rel=\"stylesheet\" type=\"text/css\" href=\"styles.css\">\n" + "\t\t<link rel=\"shortcut icon\" href=\"ufrrj.jpg\" type=\"image/jpg\"/>\n\t</head>\n")
        texto.append("\n\t<body>\n")
        texto.append("\n\t\t<header>\n")
        texto.append(
            "\t\t<a href=\"http://portal.ufrrj.br\" title=\"UFRRJ\"><img src=\"ufrrj.jpg\" class = \"imagem\" align = \"left\" alt=\"Falha na imagem\"></a>\n")
        texto.append("\t\t<h1><br>" + xml.find('universidade').find('nome').text + "</h1>\n")
        texto.append("\t\t<h1>" + xml.find('universidade').find('abreviacao').text + "</h1><br><br>\n")
        texto.append("\t\t</header>\n\n")
        texto.append("\n\t\t<div>\n")
        texto.append(
            "\t\t<img src=\"perfil.jpg\" class = \"imagemPerfil\" align = \"left\" title=\"Perfil\" alt=\"Falha na imagem\"><br>Curso: " + xml.find(
                'curso').text + "<br>\n")
        texto.append("\t\tAluno: " + xml.find('aluno').text + "<br>\n")
        texto.append("\t\tMatrícula: " + xml.find('matricula').text + "<br>\n")
        texto.append("\t\tCR médio: " + xml.find('crMedio').text + "<br>\n")
        texto.append("\t\tData geração: " + xml.find('dataGeracao').text + "<br>\n")
        texto.append("\t\tHora geração: " + xml.find('horaGeracao').text + "<br>\n")
        texto.append("\t\tCód. autenticação: " + xml.find('codigoAutenticacao').text + "<br><br>\n")
        texto.append("\t\t</div>\n")

        texto.append("\n\t\t<section>\n")
        listaPeriodos = xml.find('periodos').findall('Periodo')
        texto.append(
            "\t\t\t<img src=\"legenda.png\" class = \"imagemLegenda\" align = \"right\" title=\"legenda\" alt=\"Falha na imagem\">\n")
        for i in range(len(listaPeriodos)):
            texto.append(
                "\t\t\t<br><br><table>\n")  # \n\t\t\t<tr><th>Ano Semestre</th><th>Creditos solicitados</th><th>Creditos acumulados</th><th>Creditos obtidos</th><th>Cr periodo</th></tr><br>\n")
            texto.append("\t\t\t<tr><th>" + listaPeriodos[i].find('anoSemestre').text + "</th></tr>\n")

            texto.append(
                "\t\t\t\t<tr><th>Código disciplina</th>" + "<th>Nome disciplina</th>" + "<th>Créditos</th>" + "<th>Nota</th>" + "<th>Situação</th></tr>\n")

            listaDisciplinasAA = listaPeriodos[i].find('disciplinas').findall('AtividadeAcademica')
            for j in range(len(listaDisciplinasAA)):
                texto.append("\t\t\t\t<tr><td>" + listaDisciplinasAA[j].find('codigo').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinasAA[j].find('nome').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinasAA[j].find('creditos').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinasAA[j].find('nota').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinasAA[j].find('situacaoAA').text + "</td></tr>\n\n")

            listaDisciplinas = listaPeriodos[i].find('disciplinas').findall('Disciplina')
            for j in range(len(listaDisciplinas)):
                texto.append("\t\t\t\t<tr><td>" + listaDisciplinas[j].find('codigo').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinas[j].find('nome').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinas[j].find('creditos').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinas[j].find('nota').text + "</td>\n")
                texto.append("\t\t\t\t<td>" + listaDisciplinas[j].find('situacao').text + "</td></tr>\n\n")

            texto.append("\t\t\t</tr><td></td><td>Créditos solicitados: " + listaPeriodos[i].find(
                'creditosSolicitados').text + "</td>\n")
            texto.append(
                "\t\t\t<td>Créditos acumulados: " + listaPeriodos[i].find('creditosAcumulados').text + "</td>\n")
            texto.append("\t\t\t<td>Créditos obtidos: " + listaPeriodos[i].find('creditosObtidos').text + "</td>\n")
            texto.append("\t\t\t<td>CR período: " + listaPeriodos[i].find('crPeriodo').text + "</td></tr>\n\n")
            texto.append("\t\t\t</table>\n\n")

        texto.append("\n\t\t</section>\n")
        texto.append("\n\t\t<footer>\n")
        texto.append("\t\t\t<br>Grupinho de TEDB, 2019.")
        texto.append("\n\t\t</footer>\n")
        texto.append("\t</body>\n")
        texto.append("\n</html>")

        arq.writelines(texto)

        print(name, "gerado com sucesso.\nSalvo em:", os.path.abspath(name))

        arq.close()

        flag = False

        while (flag == False):

            option = input("\nDeseja abrir o historico gerado? (S / N)\n")

            if (option == 's' or option == 'S' or option == 'sim' or option == 'yes' or option == 'y'):
                os.startfile(os.path.abspath(name))
                flag = True

            elif (option == 'n' or option == 'N' or option == 'nao' or option == 'not' or option == 'no'):
                flag = True

            if (flag == False):
                print("Opção inválida. Tente novamente...")

        print("Visualização encerrada.")

def main():
    ''' Função para execução no terminal'''

    # no caso de implementação do menu isso deverá ser alterado
    #candidato = Candidato(input("digite o cpf do candidato a entrar:"))

    # client = ClientSocket()
    # client.connect()

    def login():
        while True:
            ans = input(
                "_______________________________"
                "1 - Fazer Login"
                "0 - Fechar Aplicação"
                "_______________________________")
            if ans == '1':
                cpf = input("Digite o CPF do Candidato:")
            elif ans == '0':
                print("Programa encerrado")
                exit(0)

    #candidato = Candidato(input())

    while True:
        ans = input(
        "_______________________________"
        "1 - submeter boletim"
        "2 - consultar status"
        "3 - entrar como outro candidato"
        "0 - sair do programa"
        "_______________________________")

        if ans == '1':
            # candidato.submeter()

            pass
        elif ans == '2':
            pass
        elif ans == '3':
            pass
        elif ans == '0':
            break
        else:
            print("Entrada não esperada")
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

def teste():
    '''Operações de teste'''
    c = Candidato("00000000002")
    c.submeter()
    # ctrlXML = ControladorXML()
    #
    # xml = ctrlXML.criarRequisicao("consultaStatus", {'cpf': '0001'})
    #
    # print(etree.tostring(xml))
    # pass

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

teste()