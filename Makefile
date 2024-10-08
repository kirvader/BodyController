PROTO_FILES = $(shell find . -type f -name '*.proto' ! -path "*/node_modules/*")
PROJECT_ROOT = $(shell pwd)
GENERATED_DIR_NAME = generated

generate-output-openapi-directory:
	mkdir -p ./generated/openapiv2

generate-all-protos: generate-output-openapi-directory
	for file in $(PROTO_FILES); do \
		filepath_relative_to_project_root=$$(echo $$file | sed 's/^..//'); \
		protoc $$filepath_relative_to_project_root \
			--proto_path=${PROJECT_ROOT} \
			--go_out=. \
			--go-grpc_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			--grpc-gateway_out . \
			--grpc-gateway_opt paths=source_relative \
			--grpc-gateway_opt generate_unbound_methods=true \
			--openapiv2_out generated/openapiv2; \
	done

pack-monorepo-in-docker: generate-all-protos
	sudo docker build -t kirvader/body-controller:latest .

push-monorepo-image-to-hub: pack-monorepo-in-docker
	sudo docker push kirvader/body-controller:latest

.PHONY: pack-monorepo-in-docker generate-all-protos generate-output-openapi-directory
