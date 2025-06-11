# Generate api go files
# Example: make gen-api s=oms
gen-api:
	goctl api go --home tpl-1.8.4 -dir api/${s} -api api/${s}/doc/${s}.api
	goctl api format --dir=api/${s} --declare
	goctl api swagger --api api/${s}/doc/${s}.api --dir doc/swagger

# Generate rpc go files
# Example: make gen-rpc s=backendrpc
gen-rpc:
	goctl rpc --home tpl protoc rpc/${s}/${s}.proto --go_out=rpc/${s} --go-grpc_out=rpc/${s} --zrpc_out=rpc/${s} -m

# Generate model
# Example: make model table=stock dir=testmodel
model:
	goctl model pg datasource --table="${table}" -dir ./rpc/model/${dir} --url="postgres://$(PG)" --schema="public" -home=tpl

apifox:
	go run doc/swagger/main.go