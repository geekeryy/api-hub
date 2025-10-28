.PHONY: gen-api gen-rpc model lint apifox build restart

# Generate api go files
# Example: make gen-api s=oms
gen-api:
	goctl api go --home tpl -dir api/${s} -api api/${s}/doc/${s}.api
	goctl api format --dir=api/${s} --declare
	goctl api swagger --api api/${s}/doc/${s}.api --dir doc/swagger

# Generate rpc go files
# Example: make gen-rpc s=backendrpc
gen-rpc:
	goctl rpc --home tpl protoc rpc/${s}/${s}.proto --go_out=rpc/${s} --go-grpc_out=rpc/${s} --zrpc_out=rpc/${s} -m

# Generate model
# Example: make model
model:
	goctl model mysql ddl --home tpl -s doc/sql/auth.sql -d rpc/model/authmodel
	goctl model mysql ddl --home tpl -s doc/sql/member.sql -d rpc/model/membermodel
	goctl model mysql ddl --home tpl -s doc/sql/admin.sql -d rpc/model/adminmodel

# Import api to apifox
# ENV: export APIFOX_TOKEN=APS-xxxxxxxxxxxxxxxxxxx
# Example: make apifox s=gateway PROJECT_ID=6567759
apifox:
	@if [ -z "$(PROJECT_ID)" ]; then echo "PROJECT_ID is not set"; exit 1; fi
	@if [ -z "$(APIFOX_TOKEN)" ]; then echo "APIFOX_TOKEN is not set"; exit 1; fi
	@swaggerfile=./doc/swagger/$(s).json; \
	jsondata=$$(cat $$swaggerfile | tr -d '\n' | sed 's/\\/\\\\/g; s/"/\\"/g'); \
	curl --location --request POST "https://api.apifox.cn/api/v1/projects/$(PROJECT_ID)/import-data" \
		--header "X-Apifox-Version: 2022-11-16" \
		--header "Authorization: Bearer $(APIFOX_TOKEN)" \
		--header "Content-Type: application/json" \
		--data-raw "$$(printf '{"importFormat": "openapi","data":"%s"}' "$$jsondata")"

# 静态代码检查
# VSCode: "go.lintFlags": ["--config=./.golangci.yml"] 
lint:
	golangci-lint fmt
	golangci-lint run

############################################ Docker本地调试 #################################################

build:
	GOOS=linux go build -o build/gateway api/gateway/gateway.go
	GOOS=linux go build -o build/oms api/oms/oms.go

restart:
	docker compose -f deploy/local/docker-compose.yml up -d