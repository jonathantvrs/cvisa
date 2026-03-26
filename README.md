# API de Gestão de Limite de Crédito

Este microsserviço é responsável por gerir exclusivamente o limite de crédito disponível para contas pré-existentes no sistema.

## Decisões Arquiteturais
### 1. Prevenção de Condições de Corrida (Concorrência e Locks)
Em sistemas financeiros, é comum que múltiplas transações ocorram no mesmo milissegundo. Se duas compras de R$ 80,00 chegarem exatamente ao mesmo tempo para uma conta com limite de R$ 100,00, um sistema sem proteção leria o limite como 100 para ambas as requisições e aprovaria as duas, deixando o saldo negativo em R$ -60,00.

Para evitar isso, implementamos uma **Estratégia de Bloqueio Pessimista (Pessimistic Locking)** usando Transações de Banco de Dados:
* Utilizamos a cláusula `SELECT ... FOR UPDATE` (via `clause.Locking` do GORM).
* **Como funciona:** Quando a primeira requisição chega, o banco de dados "tranca" aquela linha específica da conta. A segunda requisição fica em espera (na fila). Assim que a primeira transação termina (atualizando o limite para R$ 20,00) e liberta o *lock*, a segunda requisição lê o saldo atualizado (R$ 20,00) e é recusada corretamente.

### 2. Valores Financeiros em `int64` (Centavos)
Para garantir a máxima precisão financeira e evitar falhas críticas, este projeto **não utiliza** tipos de ponto flutuante (`float32` ou `float64`). O tipo `float` sofre de problemas de arredondamento binário na arquitetura dos computadores. Todos os valores monetários são armazenados e trafegados como **Inteiros em centavos**. 
* Exemplo: **R$ 100,00** é representado como `10000`.
* Exemplo: **R$ 30,50** é representado como `3050`.
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