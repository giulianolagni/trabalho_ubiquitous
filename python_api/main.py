from fastapi import FastAPI, HTTPException, Depends
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker, declarative_base
from sqlalchemy import Column, Integer, String, select
import os

# --- Configuração de Banco de Dados Assíncrono ---
# Lemos as variáveis de ambiente (vêm do Docker)
DB_USER = os.getenv("DB_USER", "user_db")
DB_PASSWORD = os.getenv("DB_PASSWORD", "password_db")
DB_HOST = os.getenv("DB_HOST", "db")
DB_NAME = os.getenv("DB_NAME", "users_db")

# Note o prefixo: postgresql+asyncpg (Driver de alta performance)
DATABASE_URL = f"postgresql+asyncpg://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:5432/{DB_NAME}"

# Engine com configurações de pool para aguentar carga
engine = create_async_engine(DATABASE_URL, echo=False, pool_size=20, max_overflow=10)
AsyncSessionLocal = sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)
Base = declarative_base()

# --- Modelo de Dados (Tabela) ---
class UserDB(Base):
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, index=True)
    email = Column(String, unique=True, index=True)
    user = Column(String, unique=True, index=True)
    password = Column(String)

# --- Schemas Pydantic (Validação) ---
class UserCreate(BaseModel):
    name: str
    email: str
    user: str
    password: str

class UserResponse(UserCreate):
    id: int
    class Config:
        orm_mode = True

# --- Inicialização da App ---
app = FastAPI()

# Cria as tabelas no arranque
@app.on_event("startup")
async def startup():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

# Dependência para injetar a sessão do banco
async def get_db():
    async with AsyncSessionLocal() as session:
        yield session

# --- Rotas CRUD (Async) ---

@app.post("/users", response_model=UserResponse, status_code=201)
async def create_user(user: UserCreate, db: AsyncSession = Depends(get_db)):
    new_user = UserDB(**user.dict())
    db.add(new_user)
    try:
        await db.commit()
        await db.refresh(new_user)
    except Exception as e:
        await db.rollback()
        # Simplificação de erro para o exemplo
        raise HTTPException(status_code=400, detail="User or Email already exists")
    return new_user

@app.get("/users", response_model=list[UserResponse])
async def read_users(skip: int = 0, limit: int = 100, db: AsyncSession = Depends(get_db)):
    # Executa a query de forma assíncrona
    result = await db.execute(select(UserDB).offset(skip).limit(limit))
    return result.scalars().all()

@app.get("/users/{user_id}", response_model=UserResponse)
async def read_user(user_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(UserDB).filter(UserDB.id == user_id))
    user = result.scalars().first()
    if user is None:
        raise HTTPException(status_code=404, detail="User not found")
    return user

@app.put("/users/{user_id}", response_model=UserResponse)
async def update_user(user_id: int, user_update: UserCreate, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(UserDB).filter(UserDB.id == user_id))
    db_user = result.scalars().first()
    if db_user is None:
        raise HTTPException(status_code=404, detail="User not found")
    
    db_user.name = user_update.name
    db_user.email = user_update.email
    db_user.user = user_update.user
    db_user.password = user_update.password
    
    await db.commit()
    await db.refresh(db_user)
    return db_user

@app.delete("/users/{user_id}")
async def delete_user(user_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(UserDB).filter(UserDB.id == user_id))
    db_user = result.scalars().first()
    if db_user is None:
        raise HTTPException(status_code=404, detail="User not found")
    
    await db.delete(db_user)
    await db.commit()
    return {"message": "User deleted successfully"}