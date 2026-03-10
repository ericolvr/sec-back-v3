# Deploy do NR1 Backend para Google Cloud Platform

## Pré-requisitos

### 1. Verificar se existe gcloud
```bash
gcloud --version
```

### 2. Configurar gcloud CLI
```bash
# Fazer login
gcloud auth login

# Configurar projeto (criar um novo projeto primeiro)
gcloud projects create nr1-back --name="NR1 Backend"
gcloud config set project nr1-back
gcloud config get-value project

# Verificar configuração
gcloud config list
```

## Fluxo recomendado (do zero)

Este repositório possui dois workflows do GitHub Actions:

- **`.github/workflows/gcp-setup.yml`**
  - Cria projeto (se não existir), habilita APIs, cria Cloud SQL e cria Secrets no Secret Manager.
  - Billing precisa ser vinculado manualmente (o workflow só orienta).
- **`.github/workflows/deploy.yml`**
  - Faz build da imagem e faz deploy no Cloud Run no serviço **`nr1-back-api`**.

Você pode seguir:

- **Opção A (recomendada)**: Rodar `gcp-setup` e depois rodar `deploy`.
- **Opção B**: Fazer tudo manualmente via `gcloud` (passos abaixo).

### 3. Habilitar APIs Necessárias
```bash
# Habilitar todas as APIs de uma vez
gcloud services enable \
    cloudbuild.googleapis.com \
    run.googleapis.com \
    containerregistry.googleapis.com \
    artifactregistry.googleapis.com \
    sqladmin.googleapis.com \
    secretmanager.googleapis.com \
    iam.googleapis.com \
    logging.googleapis.com \
    monitoring.googleapis.com
```

### 4. Configurar Autenticação Docker
```bash
# Configurar Docker para usar gcloud
gcloud auth configure-docker
```

## Setup do Banco de Dados

### 5. Criar instância Cloud SQL PostgreSQL
```bash
# Criar instância PostgreSQL
gcloud sql instances create nr1-back-db \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1 \
    --storage-type=SSD \
    --storage-size=10GB \
    --backup-start-time=03:00 \
    --enable-bin-log \
    --maintenance-window-day=SUN \
    --maintenance-window-hour=04 \
    --maintenance-release-channel=production
```

### 6. Criar usuário e banco de dados
```bash
# Definir senha do usuário root
gcloud sql users set-password root \
    --instance=nr1-back-db \
    --password=SUA_SENHA_ROOT_AQUI

# Criar usuário da aplicação
gcloud sql users create nr1_back_user \
    --instance=nr1-back-db \
    --password=SUA_SENHA_USER_AQUI

# Criar banco de dados
gcloud sql databases create nr1_back_production \
    --instance=nr1-back-db
```

### 7. Obter IP da instância Cloud SQL
```bash
gcloud sql instances describe nr1-back-db \
    --format='value(ipAddresses[0].ipAddress)'
```
Anote este IP para usar nas variáveis de ambiente.

### 8. Obter seu IP para autorização
```bash
curl ifconfig.me
```

### 9. Autorizar acesso ao banco
```bash
# Substitua SEU_IP pelo IP obtido no passo anterior
gcloud sql instances patch nr1-back-db \
  --authorized-networks=SEU_IP/32,0.0.0.0/0 \
  --quiet
```

### 10. Criar as tabelas no banco
```bash
# Conectar ao banco e executar o script SQL
gcloud sql connect nr1-back-db --user=nr1_back_user --database=nr1_back_production
```
Execute o arquivo `scripts/init.sql` no banco.

## Setup dos Secrets

### 11. Criar secrets no Secret Manager
```bash
# Senha do banco
echo -n "SUA_SENHA_USER_AQUI" | gcloud secrets create nr1-back-db-password --data-file=-

# JWT Secret (gerar uma chave aleatória)
openssl rand -base64 32 | gcloud secrets create nr1-back-jwt-secret --data-file=-

# SMTP password
echo -n "SUA_SENHA_SMTP_AQUI" | gcloud secrets create nr1-back-smtp-password --data-file=-

# Twilio credentials (se usar)
echo -n "SEU_TWILIO_SID" | gcloud secrets create nr1-back-twilio-sid --data-file=-
echo -n "SEU_TWILIO_TOKEN" | gcloud secrets create nr1-back-twilio-token --data-file=-
```

### 12. Dar permissões aos secrets
```bash
# Obter o número do projeto
PROJECT_NUMBER=$(gcloud projects describe nr1-back --format='value(projectNumber)')

# Dar permissão para o Compute Engine acessar os secrets
gcloud secrets add-iam-policy-binding nr1-back-db-password \
    --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding nr1-back-jwt-secret \
    --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding nr1-back-smtp-password \
    --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding nr1-back-twilio-sid \
    --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding nr1-back-twilio-token \
    --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"
```

## Configurar GitHub Actions (Deploy Automático)

### 13. Criar Service Account para GitHub Actions
```bash
# Criar service account
gcloud iam service-accounts create github-actions \
    --display-name="GitHub Actions Deploy"

# Obter o email da service account
PROJECT_ID=$(gcloud config get-value project)
SA_EMAIL="github-actions@${PROJECT_ID}.iam.gserviceaccount.com"

# Dar permissões necessárias
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/iam.serviceAccountUser"

# Criar chave JSON
gcloud iam service-accounts keys create github-actions-key.json \
    --iam-account=${SA_EMAIL}
```

### 14. Configurar Secret no GitHub
```bash
# 1. Copiar conteúdo do arquivo JSON
cat github-actions-key.json

# 2. No GitHub:
#    - Ir em Settings > Secrets and variables > Actions
#    - Clicar em "New repository secret"
#    - Nome: GCP_SA_KEY
#    - Valor: Colar todo o conteúdo do JSON
#    - Clicar em "Add secret"

# 3. Deletar o arquivo local (segurança)
rm github-actions-key.json
```

### 15. Habilitar Billing no GCP
```bash
# Verificar se billing está habilitado
gcloud beta billing projects describe $PROJECT_ID

# Se não estiver, habilitar via console:
# https://console.cloud.google.com/billing
```

### 16. Fazer push para main para disparar deploy
```bash
git add .
git commit -m "Deploy to production"
git push origin main
```

## Deploy via GitHub Actions (recomendado)

1. Rode o workflow **GCP Infrastructure Setup** (`gcp-setup.yml`)
   - `project_id`: `nr1-back`
   - `db_instance_name`: `nr1-back-db`
   - `db_user`: `nr1_back_user`
   - `db_name`: `nr1_back_production`

2. Vincule o Billing ao projeto `nr1-back` no Console.

3. Configure no GitHub:
   - Secrets: `GCP_SA_KEY`, `DB_ROOT_PASSWORD`, `DB_PASSWORD`, `JWT_SECRET`, `SMTP_PASSWORD`, `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`
   - Variables: `GCP_PROJECT_ID=nr1-back`, `GCP_REGION=us-central1` e demais `DB_*`, `SMTP_*`, `GCS_BUCKET_NAME`, `FRONTEND_URL`, `TWILIO_FROM`

4. Rode o workflow **Deploy to Cloud Run** (`deploy.yml`) ou faça push para `main`.

## Deploy Manual (Alternativa)

### 17. Executar o deploy via Cloud Build
```bash
# Fazer o deploy usando Cloud Build
gcloud builds submit --config=cloudbuild.yaml \
    --substitutions=_DB_HOST=IP_DO_CLOUD_SQL,_DB_USER=nr1_back_user,_DB_NAME=nr1_back_production,_TWILIO_FROM=+15551234567
```

### 14. Verificar se o deploy funcionou
```bash
# Listar serviços Cloud Run
gcloud run services list --region=us-central1

# Ver logs da aplicação
gcloud run services logs read nr1-back-api --region=us-central1
```

## Comandos Úteis

### Verificar builds em andamento
```bash
gcloud builds list --ongoing
```

### Cancelar build
```bash
gcloud builds cancel BUILD_ID
```

### Ver logs do Cloud Run
```bash
gcloud run services logs tail nr1-back-api --region=us-central1
```

### Deletar serviço (se necessário)
```bash
gcloud run services delete nr1-back-api --region=us-central1
```

### Conectar ao banco para debug
```bash
gcloud sql connect nr1-back-db --user=nr1_back_user --database=nr1_back_production
```

## Estrutura de Custos Estimados

- **Cloud SQL (db-f1-micro)**: ~$7/mês
- **Cloud Run**: Pay-per-use (muito baixo para desenvolvimento)
- **Container Registry**: ~$0.10/GB/mês
- **Secret Manager**: $0.06 por 10.000 acessos

**Total estimado**: ~$10-15/mês para ambiente de desenvolvimento/teste.
