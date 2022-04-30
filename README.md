# mac-backend

### Generate PEM and PUB 
```bash
openssl genpkey -algorithm RSA -aes-256-cbc -outform PEM -out private_key.pem -pkeyopt rsa_keygen_bits:2048
chmod 0400 private_key.pem 
ssh-keygen -y -f private_key.pem > key.pub

```
 

### Deployment
Env var

#### db config
- DB_USER=
- DB_PASS=
- DB_HOST=
- DB_PORT=
- DB_NAME=
- LOG_PATH=

#### app config
- LICENCE_KEY=
- JWT_TOKEN_SECRET=
- JWT_REFRESH_TOKEN_SECRET=


#### build config
- GO111MODULE=on
- GOFLAGS=-mod=vendor


Build docker image
```
docker build -t dmglab/mac-backend:0.1 .
```