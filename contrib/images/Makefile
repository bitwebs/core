all: iqd-env

iqd-env:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --tag bitwebs/iqd-env iqd-env

.PHONY: all iqd-env
