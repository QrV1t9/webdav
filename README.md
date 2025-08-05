# 🌐 Simple WebDAV Server in Go

---

## ✨ Features

- 🔐 **Strong Argon2 password hashing**
- 📦 **Docker & Docker Compose support**
- 🔒 **TLS encryption support**
- ⚙️ **Minimal setup**

---

##  Installation & Usage

### Without Docker

1. **Clone repository**
    ```bash
    git clone https://github.com/qrv1t9/webdav
    cd webdav
    ```
2. **Generate argon hash from your password**
    ```bash
    go run ./cmd/argon/main.go --config ./config/prod.yaml --pass "your password"
    ```
3. **Change username and password in config file (config -> prod.yaml)**
   - Edit `config/prod.yaml`
    - Replace default **username** and **password hash**
    - Set the **folder** to share via WebDAV

4. **Start the server**
    ```bash
    go run ./cmd/main.go --config ./config/prod.yaml
    ```
   
---

## With docker
1. **Clone the repository**
    ```bash
    git clone https://github.com/qrv1t9/webdav
    cd webdav
    ```

2. **Generate Argon2 hash**
    ```bash
    go run ./cmd/argon/main.go --config ./config/prod.yaml --pass "your password"
    ```

3. **Update config**
    - Edit `config/prod.yaml`:
        - Replace default **username** and **password hash**
        - Set the **folder** to share

4. **Edit `docker-compose.yaml`**
    - Update the volume path:
      ```yaml
      volumes:
        - /your/local/folder:/app/folder
      ```

5. **Start with Docker Compose**
    ```bash
    docker compose up -d
    ```
---

## Notes

- Make sure ports used in the config are open and available.

- TLS certificates must be valid and correctly referenced in `prod.yaml`.