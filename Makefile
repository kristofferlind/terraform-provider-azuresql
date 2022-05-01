HOSTNAME=registry.terraform.io
NAMESPACE=kristofferlind
NAME=azuresql
BINARY=terraform-provider-${NAME}
OS_ARCH=linux_amd64

VERSION=0.2.1

build:
	go mod download
	go build -v

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

generate:
	go mod tidy
	go generate ./...

test: install test-local-down test-local-up
	TF_ACC=local go test -cover -v ./...
	$(MAKE) test-local-down

test-release:
	goreleaser release --rm-dist --skip-publish --skip-sign --snapshot

test-local-up:
	docker-compose up -d
	docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'p@ssw0rd' -Q "CREATE DATABASE test"
	docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'p@ssw0rd' -Q "CREATE DATABASE test2"
	docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'p@ssw0rd' -Q "sp_configure 'contained database authentication', 1; RECONFIGURE"
	docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'p@ssw0rd' -Q "CREATE DATABASE contained_test CONTAINMENT = PARTIAL"
	docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'p@ssw0rd' -Q "CREATE DATABASE contained_test2 CONTAINMENT = PARTIAL"

test-local-down:
	docker-compose down

.PHONY: build install release-test generate test test-local-up test-local-down
