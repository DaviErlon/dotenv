# dotenv

Loader de arquivos `.env` em Go com inicialização automática, estado global e acesso seguro para concorrência.

## ✨ Features

- Carregamento automático do `.env` via `init()`
- Armazenamento em estado global
- Thread-safe (`sync.RWMutex`)
- Suporte a comentários (`#`)
- API simples e direta

---

## 📦 Instalação

```bash
go get github.com/DaviErlon/dotenv