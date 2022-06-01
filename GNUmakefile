GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: testacc

build: fmt
	go build -v .

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

docs:
	go generate ./...

.PHONY: build test testacc vet fmt fmtcheck docs
