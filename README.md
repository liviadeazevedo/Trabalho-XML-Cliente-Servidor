# Trabalho-XML-Cliente-Servidor

A universidade UFRRJ está aceitando inscrições para o seu programa de Pós-Graduação. Para tanto, os candidatos que quiserem concorrer a uma vaga devem realizar a inscrição submetendo seus boletins escolares para o sistema da UFRRJ em formato XML, respeitando o XML Schema a ser definido, conforme conversado em aula. O sistema da UFRRJ deve disponibilizar dois métodos para serem invocados:

/*
envia um boletim como parâmetro e retorna um número inteiro (0 - sucesso, 1 - XML inválido, 2 - XML mal-formado, 3 - Erro Interno)
*/

**int submeter(Boletim)**

/*
consulta o status da inscrição do candidato com o CPF informado como parâmetro. Possíveis retornos: 0 - Candidato não encontrado, 1 - Em processamento, 2 - Candidato Aprovado e Selecionado, 3 - Candidato Aprovado e em Espera, 4 - Candidato Não Aprovado.
*/

**int consultaStatus(String CPF)**

Toda a comunicação deve ser feita através de XML, tanto do cliente para o servidor (request) quanto no caminho inverso (response).

O que deve ser apresentado no dia da entrega? Os sistemas se comunicando independentemente das linguagens em que foram implementados. Na implementação do método submeter, os servidores devem estar validando e verificando os XMLs dos boletins utilizando um XML Schema e retornando os códigos corretos, enquanto os clientes devem exibir a mensagem de retorno para o usuário. No método consultaStatus, o servidor deve ser capaz de identificar o cliente pelo CPF e retornar o código adequado e, de novo, os clientes devem ser capazes de exibir a mensagem correta. Vamos usar os seguintes valores de CPF para teste no método consultaStatus:

- 00000000000 -> Código 0
- 00000000001 -> Código 1
- 00000000002 -> Código 2
- 00000000003 -> Código 3
- 00000000004 -> Código 4

Obs: Para cpfs diferentes desses, retornar o código 0(Canditado não encontrado).

------------------

Para implmentação, usaremos as linguagens:

**CLiente:** Python
**Servidor:** Go