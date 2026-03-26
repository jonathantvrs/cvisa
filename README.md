# API de Gestão de Limite de Crédito

Este microsserviço é responsável por gerir exclusivamente o limite de crédito disponível para contas pré-existentes no sistema.

## Decisão Arquitetural: Valores em `int64` (Centavos)
Para garantir a máxima precisão financeira, este projeto **não utiliza** tipos de ponto flutuante (`float32` ou `float64`) para representar dinheiro.

O tipo `float` sofre de problemas de arredondamento. Para evitar bugs no saldo das contas, todos os valores monetários são armazenados e trafegados como **Inteiros em centavos**. 
* Exemplos: 
    - **R$ 100,00** é representado como `10000`.
    - **R$ 30,50** é representado como `3050`.
---
## Endpoints da API

1. Definir Limite Inicial

Cria o registro de limite para um Account_ID existente.

```bash
POST /accounts/:accountId/limit

Request Body:
JSON
{
  "available_credit_limit": 10000
}

Response (201 Created):
JSON
{
  "account_id": 5000,
  "available_credit_limit": 10000
}
```

2. Visualizar Limite

Retorna o limite de crédito disponível no momento.
```bash
GET /accounts/:accountId/limit

Response (200 OK):
JSON
{
  "account_id": 5000,
  "available_credit_limit": 10000
}
```

3. Processar Transação (Atualizar Limite)

Aplica a regra de negócio central: abate o limite em caso de saque (valor negativo) ou aumenta o limite em caso de pagamento (valor positivo). **O saldo nunca ficará menor que zero**.

```bash
PATCH /accounts/:accountId/limit

- Cenário A: Saque de R$ 30,00

Request Body:
{
  "amount": -3000
}

Response (200 OK):
{
  "message": "Limite atualizado com sucesso"
}

- Cenário B: Tentativa de Saque Acima do Limite (Ex: R$ 80,00 tendo apenas 70,00)

Request Body:
{
  "amount": -8000
}
Response (422 Unprocessable Entity):
{
  "error": "limite insuficiente"
}

- Cenário C: Pagamento de R$ 20,00

Request Body:
{
  "amount": 2000
}
Response (200 OK):
{
  "message": "Limite atualizado com sucesso"
}
```