# Usa uma versão leve do Python
FROM python:3.11-slim

WORKDIR /app

# Variáveis de ambiente para evitar arquivos .pyc
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Instala dependências do sistema necessárias para compilar bibliotecas
RUN apt-get update && apt-get install -y build-essential libpq-dev

# Instala dependências do Python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copia o código
COPY . .

# Expõe a porta 8000 (Padrão do FastAPI/Uvicorn)
EXPOSE 8000

# Comando de inicialização
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]