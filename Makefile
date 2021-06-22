.PHONY: init
init:
	$(MAKE) up
	$(MAKE) build

.PHONY: up
up:
	docker-compose up -d \

.PHONY: test
test:
	go test -v -mod=readonly -cover ./...

.PHONY: lint
lint:: \
	lint-rules \
	golangci-lint

.PHONY: lint-fix
lint-fix:: \
	lint-rules \
	golangci-lint-fix

GOLANGCI_LINT_VERSION=v1.31.0
GOLANGCI_LINT_DIR=$(shell go env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_BIN=$(GOLANGCI_LINT_DIR)/golangci-lint
$(GOLANGCI_LINT_BIN):
	curl -vfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_BIN)

.PHONY: golangci-lint
golangci-lint: install-golangci-lint
	$(GOLANGCI_LINT_BIN) -v run

.PHONY: golangci-lint-fix
golangci-lint-fix: install-golangci-lint
	$(GOLANGCI_LINT_BIN) -v run --fix

.PHONY: lint-rules
lint-rules:: ensure-command-pcregrep
	# Disallowed files.
	find . -name ".DS_Store" | xargs -n 1 -I {} sh -c 'echo {} && false'

	# Don't use upper case letter in file and directory name.
	# The convention for separator in name is:
	# - file: "_"
	# - directory in "/cmd": "-"
	# - other directory: shouldn't be separated
	! find . -name "*.go" | pcregrep "[[:upper:]]"

	# Don't use more than 2 block of packages in an import statement.
	! pcregrep -rnM --include=".+\.go$$" "import \(\n([\t \w\"/.-]+\n)+(\n([\t \w\"/.-]+\n)+){2,}\)" .
	# Don't use more than 1 import statement.
	! pcregrep -rnM --include=".+\.go$$" "((.*\n)+import){2,}" .

	# Don't export type/function/variable/constant in main package/test.
	! pcregrep -rnM --include=".+\.go$$" --exclude=".+_test\.go$$" "^package main\n(.*\n)*(type|func|var|const) [[:upper:]]" .
	! pcregrep -rnM --include=".+\.go$$" --exclude=".+_test\.go$$" "^package main\n(.*\n)*(var|const) \(\n((\t.*)?\n)*\t[[:upper:]]" .
	! pcregrep -rnM --include=".+_test\.go$$" "^(type|var|const) [[:upper:]]" .
	! pcregrep -rnM --include=".+_test\.go$$" "^(var|const) \(\n((\t.*)?\n)*\t[[:upper:]]" .
#	! pcregrep -rnM --include=".+_test\.go$$" "^func [[:upper:]]" . | pcregrep -v ":func (Test|Benchmark).*\((t|b) \*testing\.(T|B)\) {"

	# Write meaningful comments instead of "...".
	! pcregrep -rn --include=".+\.go$$" "\/\/.*\.\.\.$$" .

	# Don't add a space after "//nolint:"
	! pcregrep -rnF --include=".+\.go$$" "//nolint: " .

	# Don't declare a var block inside a function.
	! pcregrep -rn --include=".+\.go$$" "^\t+var \($$" .

	# Call golang-libraries/panichandler.Recover() with defer.
	! pcregrep -rnF --include=".+\.go$$" "panichandler.Recover()" . | pcregrep -vF "defer"


	# Don't use context.TODO().
	! pcregrep -rnF --include=".+\.go$$" "context.TODO()" .

	# For nil error do `t.Fatal("no error")` instead of using testutils.
	! pcregrep -rnM --include=".+_test\.go$$" "if err == nil \{\n.+testutils\." .


	# Use httptest.NewRequest() instead of http.NewRequest() in tests.
	! pcregrep -rnF --include=".+_test\.go$$" --exclude-dir="^httpclientrequest$$" "http.NewRequest(" .
	# Don't use http.DefaultServeMux directly or indirectly.
	! pcregrep -rn --include=".+\.go$$" --exclude-dir="^debugutils$$" "http\.(Handle\(|HandleFunc\(|DefaultServeMux)" .
	# Use httpclientrequest instead of http.DefaultClient|Get()|Head()|Post()|PostForm().
	! pcregrep -rn --include=".+\.go$$" --exclude-dir="^httpclientrequest$$" "http\.(DefaultClient|Get\(|Head\(|Post\(|PostForm\()" . | grep -vF "TODO"

	# Don't use JSON "omitempty" with a `time.Time` because it has no zero value. Convert it to a pointer or remove omitempty.
	! pcregrep -rn --include=".+\.go$$" "\stime.Time\s+\`.*json:\".*omitempty.*\`" .

	# Use mgobsonutils.ParseObjectIdHex() instead of bson.ObjectIdHex().
	! pcregrep -rnF --include=".+\.go$$" --exclude-dir="^mgobsonutils$$" "bson.ObjectIdHex" .
	# Don't mention "mgo" in projects migrated to the official MongoDB driver.
	(pcregrep -rn --include=".+\.go$$" "(github\.com\/globalsign\/mgo|github\.com\/DTSL\/golang-libraries\/mgo)" . > /dev/null) || (! pcregrep -rni --include=".+\.go$$" "mgo" .)
	# Use MongoDB options constructors and setters instead of the struct.
	! pcregrep -rn --include=".+\.go$$" "&options\.\w+\{" .

	# Use constructors from golang-libraries/redisutils.
	! pcregrep -rn --include=".+\.go$$" --exclude-dir="^redis(test|utils)$$" "redis\.New(Client|FailoverClient|ClusterClient|FailoverClusterClient|Ring|SentinelClient|UniversalClient)\(" .

	# Use jsonlog.NewErrorOptionalFile().
	! pcregrep -rn --include=".+\.go$$" "jsonlog\.(New|NewFile|NewOptionalFile)\(" .


.PHONY: mod-update
mod-update:
	go get -v -u -d all
	$(MAKE) mod-tidy

.PHONY: mod-tidy
mod-tidy:
	rm -f go.sum
	go mod tidy -v

.PHONY: clean
clean::
	go clean -cache -testcache

.PHONY: ensure-command-pcregrep
ensure-command-pcregrep:
	$(call ENSURE_COMMAND,pcregrep)

.PHONY: mysql-shell
mysql-shell:
	docker-compose exec mysql mysql -uroot -proot

GO_BUILD_DIR=build
.PHONY: build
build::
	mkdir -p $(GO_BUILD_DIR)
	go build -v -mod=readonly -ldflags="-s -w -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...
