# Trabalho Ubiquitous: Benchmark Go vs Python

Comparação de performance entre APIs em **Go** e **Python**, utilizando **Docker**, **Nginx** e **PostgreSQL**.

## 🚀 Como Rodar

Abra o terminal na pasta do projeto e execute:

```bash
docker-compose up --build
````

## 🔗 Links para Teste (Navegador)

  * **API Python:** http://localhost:4000/python/users
  * **API Go:** http://localhost:4000/go/users

## 🛠️ Como Usar (Endpoints)

Todas as requisições devem ser feitas na porta **4000** (Gateway).

| Ação | Método | URL |
| :--- | :---: | :--- |
| **Criar Usuário** | `POST` | `http://localhost:4000/go/users` ou `/python/users` |
| **Listar Usuários** | `GET` | `http://localhost:4000/go/users` ou `/python/users` |

**Exemplo de JSON para POST:**

```json
{
  "name": "Teste Rápido",
  "email": "teste@exemplo.com",
  "user": "usuario_teste",
  "password": "123"
}
```

## 📊 Testes de Desempenho

Para reproduzir os testes de carga, utilize o arquivo `teste_ubiquitous.jmx` com o **Apache JMeter**.

-----

**Autores:** Giuliano Chiochetta Lagni e Hugo Pizzatto
