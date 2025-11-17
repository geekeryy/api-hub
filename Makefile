.PHONY: gen-api gen-rpc model sync-api-doc lint test apifox build restart deploy

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
# Example: make model s=auth
model:
	goctl model mysql ddl --home tpl -s doc/sql/$(s).sql -d rpc/model/$(s)model


# Sync api doc to apifox
# ENV: export APIFOX_TOKEN=APS-xxxxxxxxxxxxxxxxxxx
# Example: make sync-api-doc s=gateway PROJECT_ID=xxx
sync-api-doc:
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

test:
	go test -v ./...

############################################ 辅助命令 #################################################

# Import api to apifox
# ENV: export APIFOX_TOKEN=APS-xxxxxxxxxxxxxxxxxxx
# Example: make apifox
apifox:
	make sync-api-doc s=gateway PROJECT_ID=6567759
	make sync-api-doc s=oms PROJECT_ID=7317140


############################################ Docker本地调试 #############################################

build:
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="-s -w" -o build/gateway api/gateway/gateway.go
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="-s -w" -o build/oms api/oms/oms.go
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="-s -w" -o build/auth rpc/auth/auth.go
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="-s -w" -o build/user rpc/user/user.go
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="-s -w" -o build/monitor rpc/monitor/monitor.go

restart: build
	docker compose -f deploy/local/docker-compose.yml up -d
	docker compose -f deploy/local/docker-compose.yml restart


############################################ Docker镜像构建 #############################################
IMAGE_REPO=ccr.ccs.tencentyun.com/private-0
IMAGE_NAME=api-hub-all
IMAGE_TAG=latest

tag-amd64:
	make build ARCH=amd64
	upx build/* || true
	docker buildx build --platform linux/amd64 -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-amd64 -f Dockerfile .
	docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-amd64

tag-arm64:
	make build ARCH=arm64
	upx build/* || true
	docker buildx build --platform linux/arm64 -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-arm64 -f Dockerfile .
	docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-arm64

tag: tag-amd64 tag-arm64
	docker manifest rm $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) || true
	docker manifest create $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-amd64 $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)-arm64; 
	docker manifest push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG); 
