# desafio_rate-limiter

O Rate Limiter opera como um **middleware HTTP** que intercepta cada requisição recebida pelo servidor. Sua lógica de funcionamento é a seguinte:

1.  **Identificação da Requisição**:
    * Primeiro, verifica se o cabeçalho `API_KEY` está presente.
    * Se estiver, a requisição é identificada pelo **Token de Acesso**.
    * Se não estiver, a requisição é identificada pelo **Endereço IP** do cliente.

2.  **Determinação do Limite**:
    * **Para Tokens**:
        * Se o token for um dos **tokens específicos configurados** (ex: `TOKEN_LIMIT_meu_token`), ele usará o limite definido para esse token.
        * Se o token estiver presente, mas **não for um token específico configurado**, ele usará um limite padrão para tokens (`RATE_LIMIT_TOKEN_DEFAULT`).
    * **Para IPs (quando não há token)**:
        * Ele usará o limite padrão por IP (`RATE_LIMIT_IP`).
    * **Prioridade**: A configuração do limite do **token de acesso sempre se sobrepõe** à configuração do limite por IP. Se um token é fornecido, a regra do token é aplicada, independentemente do limite de IP.

3.  **Contagem e Bloqueio**:
    * Cada requisição incrementa um contador associado à chave (IP ou Token) no Redis.
    * Se o contador exceder o limite configurado dentro de uma janela de tempo específica (definida por `BLOCK_DURATION_SECONDS`):
        * O IP ou Token é marcado como **bloqueado** no Redis por essa mesma duração.
        * Requisições subsequentes para aquela chave (IP ou Token) serão bloqueadas imediatamente até que o tempo de bloqueio expire.

4.  **Respostas HTTP**:
    * **`HTTP 200 OK`**: Requisição permitida.
    * **`HTTP 429 Too Many Requests`**: Requisição negada por exceder o limite ou por estar bloqueada. A mensagem retornada informa o motivo (`Você atingiu o número máximo de requisições. Tente novamente mais tarde.` ou `Você foi bloqueado devido a muitas requisições.`).

5.  **Persistência**: Todas as contagens e estados de bloqueio são armazenados e consultados no Redis.

---

## Como Configurar o Rate Limiter

As configurações são realizadas através do `.env` na pasta raiz do projeto.

Crie um arquivo `.env` com as seguintes variáveis:

* **`RATE_LIMIT_IP`**: (Obrigatório) Número máximo de requisições permitidas por IP.
    * Exemplo: `RATE_LIMIT_IP=10`
* **`RATE_LIMIT_TOKEN_DEFAULT`**: (Obrigatório) Número máximo de requisições permitidas para tokens de acesso que não possuem uma configuração específica.
    * Exemplo: `RATE_LIMIT_TOKEN_DEFAULT=100`
* **`TOKEN_LIMIT_<NOME_DO_TOKEN>`**: (Opcional) Para configurar limites específicos para tokens de acesso. Substitua `<NOME_DO_TOKEN>` pelo valor exato do token que será enviado no cabeçalho `API_KEY`. Você pode ter múltiplas entradas para diferentes tokens.
    * Exemplo: `TOKEN_LIMIT_my_secret_token=50`, `TOKEN_LIMIT_another_token=200`
* **`BLOCK_DURATION_SECONDS`**: (Obrigatório) Duração em segundos que um IP ou Token permanecerá bloqueado após exceder o limite. Esta também define a "janela" de tempo para a contagem de requisições.
    * Exemplo: `BLOCK_DURATION_SECONDS=300` (5 minutos)
* **`REDIS_ADDR`**: (Obrigatório) Endereço do servidor Redis (ex: `localhost:6379` ou `redis:6379` se usando Docker).
* **`REDIS_PASSWORD`**: (Opcional) Senha para autenticação no Redis. Deixe em branco se não houver senha.
* **`REDIS_DB`**: (Opcional) Número do banco de dados Redis a ser utilizado (padrão é `0`).

**Exemplo de arquivo `.env` completo:**

RATE_LIMIT_IP=5
RATE_LIMIT_TOKEN_DEFAULT=10
TOKEN_LIMIT_token_secreto_a=15
BLOCK_DURATION_SECONDS=300

REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

**Pra rodar**
docker-compose up --build -d

# Testar limite por IP (exceder o limite de RATE_LIMIT_IP)
for i in $(seq 1 6); do curl -s -w "%{http_code}\n" http://localhost:8080/; done

# Testar limite com token específico (ex: token_vip com limite 20)
for i in $(seq 1 16); do curl -s -H "token_secreto_a" -w "%{http_code}\n" http://localhost:8080/; done

# Testar sobreposição (IP bloqueado, mas token_alto permitido)
# 1. Bloqueie o IP primeiro (faça mais de RATE_LIMIT_IP requisições sem token)
# 2. Depois, faça requisições com o token:
curl -s -H "API_KEY: token_secreto_b" -w "%{http_code}\n" http://localhost:8080/

docker-compose down