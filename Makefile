.PHONY: import
setup:
	brew install tilt
	brew install kind
	brew install kubectx
	brew install k9s
	brew install txn2/tap/kubefwd
	brew install entr
	brew install coreutils
	brew install kustomize
cluster:
	./scripts/kind-with-registry.sh
import:
	./scripts/import-labs.sh
