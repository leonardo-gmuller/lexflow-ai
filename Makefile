include .env

PROJECT ?= lexflow-ai

DOCKER_COMPOSE_FILE_BUILD=build/docker-compose.yml

DOCKER_COMPOSE_FILE_LOCAL=docker-compose.yml

GOLANGCI_LINT_PATH=$$(go env GOPATH)/bin/golangci-lint
GOLANGCI_LINT_VERSION=1.59.0

MIGRATION_FOLDER_PATH=internal/app/gateway/postgres/migrations
GOLANG_MIGRATE_PATH=$$(go env GOPATH)/bin/golang-migrate
GOLANG_MIGRATE_VERSION=4.18.1

GO_PATH=$$(go env GOPATH)/bin

VERSION=v1.0.0


# Detect the OS automatically
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

# For Windows, we need a special case
ifeq ($(OS), mingw32)
	OS := windows
endif

build-docker:
	docker compose -f $(DOCKER_COMPOSE_FILE_BUILD) -p $(PROJECT) down --remove-orphans
	docker compose -f $(DOCKER_COMPOSE_FILE_BUILD) -p $(PROJECT) up --remove-orphans

start-docker:
	docker compose -f $(DOCKER_COMPOSE_FILE_LOCAL) -p $(PROJECT) down --remove-orphans
	docker compose -f $(DOCKER_COMPOSE_FILE_LOCAL) -p $(PROJECT) up --remove-orphans

mock-repositories:
	@find internal/app/gateway/postgres/repositories -name '*.go' | grep -v '_test.go' | while read file; do \
		if grep -q 'interface' $$file; then \
			mockgen -source=$$file -destination=internal/app/tests/mocks/mock_$$(basename $$file .go)_repository.go -package=mocks; \
		fi \
	done

generate-mocks:
	@mkdir -p internal/app/tests/mocks && \
	if [ -z "$(source)" ]; then \
		echo "Por favor, defina o source"; \
		exit 1; \
	else \
		find $(source) -name '*.go' | grep -v '_test.go' | while read file; do \
			if grep -q 'interface' $$file; then \
				mockgen -source=$$file -destination=internal/app/tests/mocks/mock_$$(basename $$file .go).go -package=mocks; \
			fi \
		done; \
	fi

lint:
	@echo "==> Installing golangci-lint"
ifeq (,$(findstring $(GOLANGCI_LINT_VERSION),$(shell which $(GOLANGCI_LINT_PATH) && eval $(GOLANGCI_LINT_PATH) version)))
	@echo "installing golangci-lint v$(GOLANGCI_LINT_VERSION)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION)
else
	@echo "already installed: $(shell eval $(GOLANGCI_LINT_PATH) version)"
endif
	@echo "==> Running golangci-lint"
	@$(GOLANGCI_LINT_PATH) run -c ./.golangci.yml --fix

migrate-up:
	migrate -path $(MIGRATION_FOLDER_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path $(MIGRATION_FOLDER_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

migrate-new:
	@echo "==> Installing golang-migrate"
	@echo "==> Checking for golang-migrate"
	@if [ -x "$(GOLANG_MIGRATE_PATH)" ]; then \
		echo "Found golang-migrate at $$(which $(GOLANG_MIGRATE_PATH))"; \
		CURRENT_VERSION=$$($(GOLANG_MIGRATE_PATH) --version 2>&1); \
		VERSION_NUM=$$(echo "$$CURRENT_VERSION" | awk '{print $$NF}'); \
		if [ "$$VERSION_NUM" = "$(GOLANG_MIGRATE_VERSION)" ]; then \
			echo "Already installed: $$VERSION_NUM"; \
			exit 0; \
		else \
			echo "Version mismatch, updating..."; \
		fi; \
	else \
		echo "golang-migrate not found, installing..."; \
	fi; \
	if [ "$(OS)" = "linux" ]; then \
		TEMP_DIR=$$(mktemp -d); \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v$(GOLANG_MIGRATE_VERSION)/migrate.linux-amd64.tar.gz | tar xvz -C $$TEMP_DIR; \
		mkdir -p $$(go env GOPATH)/bin; \
		mv $$TEMP_DIR/migrate $(GOLANG_MIGRATE_PATH); \
		chmod +x $(GOLANG_MIGRATE_PATH); \
		rm -rf $$TEMP_DIR; \
	elif [ "$(OS)" = "darwin" ]; then \
		TEMP_DIR=$$(mktemp -d); \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v$(GOLANG_MIGRATE_VERSION)/migrate.darwin-amd64.tar.gz | tar xvz -C $$TEMP_DIR; \
		mkdir -p $$(go env GOPATH)/bin; \
		mv $$TEMP_DIR/migrate $(GOLANG_MIGRATE_PATH); \
		chmod +x $(GOLANG_MIGRATE_PATH); \
		rm -rf $$TEMP_DIR; \
	elif [ "$(OS)" = "windows" ]; then \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v$(GOLANG_MIGRATE_VERSION)/migrate.windows-amd64.zip -o migrate.zip; \
		TEMP_DIR=$$(mktemp -d); \
		unzip migrate.zip -d $$TEMP_DIR; \
		mv $$TEMP_DIR/migrate.exe $$(go env GOPATH)/bin/golang-migrate.exe; \
		rm -rf $$TEMP_DIR migrate.zip; \
	fi
	@echo "==> Creating new migration files for ${name}..."
	$(GOLANG_MIGRATE_PATH) create -ext sql -dir $(MIGRATION_FOLDER_PATH) -seq ${name}


# ---- Config ----
PKG=./...
COVER_DIR=./.coverage
COVER_PROFILE=$(COVER_DIR)/coverage.out
COVER_HTML=$(COVER_DIR)/coverage.html
COVER_FUNC=$(COVER_DIR)/coverage.func.txt
COVER_PKG=$(COVER_DIR)/coverage.pkg.txt

# ---- Helpers ----
GO_TEST=go test
GO_TOOL=go tool

# ---- Targets ----
.PHONY: help
help:
	@echo "Targets:"
	@echo "  make test          - go test com -race -cover e profile consolidado"
	@echo "  make test_fast     - go test sem -race (mais rápido)"
	@echo "  make cover         - mostra resumo de cobertura por função"
	@echo "  make cover_html    - gera e abre relatório HTML"
	@echo "  make cover_pkg     - cobertura por pacote"
	@echo "  make tidy          - go mod tidy"
	@echo "  make clean         - limpa artefatos de cobertura"

$(COVER_DIR):
	mkdir -p $(COVER_DIR)

.PHONY: test
test: $(COVER_DIR)
	$(GO_TEST) -race -covermode=atomic -coverprofile=$(COVER_PROFILE) $(PKG)

.PHONY: test_fast
test_fast: $(COVER_DIR)
	$(GO_TEST) -covermode=atomic -coverprofile=$(COVER_PROFILE) $(PKG)

.PHONY: cover
cover: test
	$(GO_TOOL) cover -func=$(COVER_PROFILE) | tee $(COVER_FUNC)

.PHONY: cover_html
cover_html: test
	$(GO_TOOL) cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Relatório em: $(COVER_HTML)"
	@# tentar abrir automaticamente (macOS/Linux comuns)
	@(command -v open >/dev/null && open $(COVER_HTML)) || (command -v xdg-open >/dev/null && xdg-open $(COVER_HTML)) || echo "Por favor, abra o arquivo manualmente."

# Cobertura por pacote (usa go list + go test por pacote)
.PHONY: cover_pkg
cover_pkg: $(COVER_DIR)
	@echo "package,coverage" > $(COVER_PKG)
	@for p in $$(go list $(PKG)); do \
	  out=$$(mktemp); \
	  $(GO_TEST) -covermode=atomic -coverprofile=$$out $$p >/dev/null || exit 1; \
	  pct=$$($(GO_TOOL) cover -func=$$out | tail -n1 | awk '{print $$3}'); \
	  echo "$$p,$$pct" >> $(COVER_PKG); \
	  rm -f $$out; \
	done; \
	echo "Gerado: $(COVER_PKG)"; \
	column -t -s, $(COVER_PKG) | sed '1s/.*/\x1b[1m&\x1b[0m/'

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(COVER_DIR)


# ==== CONFIGURAÇÕES ====
TARGET_OS ?= linux
TARGET_ARCH ?= amd64
BINARY_NAME ?= lexflow-ai
MAIN_PKG ?= ./cmd

GO_MIN_VERSION = 1.24

GO_DOWNLOAD_URL_LINUX = https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
GO_DOWNLOAD_URL_DARWIN = https://go.dev/dl/go1.24.0.darwin-amd64.tar.gz


# ==== VERIFICAÇÃO DE GO ====

.PHONY: check-go
check-go:
	@echo "➡️  Verificando instalação do Go..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go não encontrado. Instalando..."; \
		$(MAKE) install-go; \
	else \
		echo "✔️ Go encontrado: $$(go version)"; \
		$(MAKE) check-go-version; \
	fi

.PHONY: check-go-version
check-go-version:
	@CURRENT=$$(go version | awk '{print $$3}' | sed 's/go//'); \
	REQ=$(GO_MIN_VERSION); \
	if [ "$$(printf '%s\n%s\n' $$REQ $$CURRENT | sort -V | head -n1)" != "$$REQ" ]; then \
		echo "⚠️  Versão do Go ($$CURRENT) é inferior ao mínimo requerido ($$REQ)."; \
		echo "Instalando nova versão..."; \
		$(MAKE) install-go; \
	else \
		echo "✔️ Versão do Go é compatível"; \
	fi

.PHONY: install-go
install-go:
	@echo "⬇️  Baixando Go ${GO_MIN_VERSION} ..."
	@if [ "$$(uname -s)" = "Linux" ]; then \
		curl -fsSL $(GO_DOWNLOAD_URL_LINUX) -o /tmp/go.tar.gz; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		curl -fsSL $(GO_DOWNLOAD_URL_DARWIN) -o /tmp/go.tar.gz; \
	else \
		echo "❌ Sistema operacional não suportado para instalação automática do Go"; \
		exit 1; \
	fi
	@echo "📦 Instalando Go..."
	@sudo rm -rf /usr/local/go
	@sudo tar -C /usr/local -xzf /tmp/go.tar.gz
	@rm /tmp/go.tar.gz
	@echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.bashrc
	@echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.zshrc
	@echo "✔️ Go instalado com sucesso!"
	@echo "ℹ️  Abra um novo terminal ou rode: source ~/.bashrc || source ~/.zshrc"


# ==== BUILD ====

.PHONY: build
build: check-go
	@echo "🚀 Building bin/$(BINARY_NAME)"
	@echo "➡️  Target: GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH)"
	@echo "📦 Version: $(VERSION)"
	@mkdir -p bin
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build \
	-ldflags="-s -w -X 'main.Version=$(VERSION)'" \
	-o bin/$(BINARY_NAME) $(MAIN_PKG)
	@echo "✔️ Build finalizado!"

.PHONY: build-macos
build-macos: check-go
	@echo "🚀 Building bin/$(BINARY_NAME)-macos"
	@echo "➡️  Target: GOOS=darwin GOARCH=$(TARGET_ARCH)"
	@echo "📦 Version: $(VERSION)"
	@mkdir -p bin
	GOOS=darwin GOARCH=$(TARGET_ARCH) go build \
	-ldflags="-s -w -X 'main.Version=$(VERSION)'" \
	-o bin/$(BINARY_NAME)-macos $(MAIN_PKG)
	@echo "✔️ Build finalizado!"

# ==== START ====

.PHONY: start
start: check-go
	@echo "▶️  Iniciando aplicação..."
	go run $(MAIN_PKG)

# ==== CLEAN BUILD ARTIFACTS ====
.PHONY: clean-build
clean-build:
	@echo "🧹 Limpando artefatos de build..."
	rm -rf bin/
	@echo "✔️ Artefatos de build removidos."


# ==== CREATE OPENAI KEY SCRIPT ====
.PHONY: create-openai-project
create-openai-project:
	@echo "▶️  Criando projeto OpenAI para tenant..."
	@chmod +x create_openai_project.sh
	@./create_openai_project.sh $(TENANT_NAME)
	@echo "✔️ Projeto OpenAI criado com sucesso!"

# ==== CREATE DB ====
.PHONY: create-db
create-db:
	@echo "▶️  Criando banco de dados PostgreSQL..."
	@psql -h $(DB_HOST) -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"
	@echo "✔️ Banco de dados criado com sucesso!"

# ==== DROP DB ====
.PHONY: drop-db
drop-db:
	@echo "▶️  Deletando banco de dados PostgreSQL..."
	@psql -h $(DB_HOST) -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "✔️ Banco de dados deletado com sucesso!"

# ==== AIR RUN ====
.PHONY: air-run
air-run: check-air
	@echo "▶️  Iniciando aplicação com Air..."
	$(GO_PATH)/air -c .air.toml
	@echo "✔️ Aplicação iniciada com Air!"

# ==== CHECK AIR ====
.PHONY: check-air
check-air:
	@echo "➡️  Verificando instalação do Air..."
	@if ! command -v $(GO_PATH)/air >/dev/null 2>&1; then \
		echo "❌ Air não encontrado. Instalando..."; \
		go install github.com/air-verse/air@latest; \
	else \
		echo "✔️ Air encontrado: $$(air -v)"; \
	fi