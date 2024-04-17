DOCKER_REPOSITORY := docker.izmno.be
DOCKER_IMAGE      := kam-chemfig
VERSION           := 0.1.0

DOCKER_FULL_IMAGE := $(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(VERSION)
DOCKER_TAGS       := $(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):latest
DOCKER_TAGS       += $(DOCKER_FULL_IMAGE)

CHEMFIG_GENERATOR := /usr/local/bin/chemfig

OUTPUT_DIR        := out

TEX_SOURCES       := $(shell $(CHEMFIG_GENERATOR) --output-dir=$(OUTPUT_DIR) targets)
DVI_TARGETS       := $(patsubst %.tex,%.dvi,$(TEX_SOURCES))
SVG_TARGETS       := $(patsubst %.tex,%.svg,$(TEX_SOURCES))
PNG_TARGETS       := $(patsubst %.tex,%.png,$(TEX_SOURCES))

.PHONY: build

build: MAKE_TARGET=_build
build: docker-make

_build: sources $(DVI_TARGETS) $(SVG_TARGETS) $(PNG_TARGETS) _clean

_clean:
	@find $(OUTPUT_DIR) -type f ! \( \
		-name '*.tex' \
		-o -name '*.pdf' \
		-o -name '*.png' \
		-o -name '*.svg' \
		-o -name '*.dvi' \) \
		-delete

sources:
	$(CHEMFIG_GENERATOR) --output-dir=$(OUTPUT_DIR) generate

%.dvi: %.tex
	mkdir -p $(dir $@)
	latex --output-directory=$(dir $@) $<

%.svg: %.dvi
	mkdir -p $(dir $@)
	dvisvgm --no-fonts --no-styles --optimize=all --output=$@ $<

%.png: %.svg
	mkdir -p $(dir $@)
	convert -background none -density 1200 -resize 1200x $< $@

##
## Go utility
##
build-go:
	go build -o bin/chemfig ./cmd/main.go

##
## Packaging
##
## Create a zip archive with all files with the given extension
## in the given subdirectory
##
## Variables:
##   PACKAGE_DIR:      Subdirectory to package
##   PACKAGE_FILETYPE: File extension to package
##   PACKAGE_FILENAME: Name of the zip archive
## Example:
##   make package PACKAGE_DIR=2024-04-16 PACKAGE_FILETYPE=png
##
PACKAGE_DIR       ?= 2024-04-16
PACKAGE_FILETYPE  ?= png
PACKAGE_FILENAME  ?= structuren.zip

_PACKAGE_TARGETS  := $(shell find $(OUTPUT_DIR)/$(PACKAGE_DIR) -type f -name '*.$(strip $(PACKAGE_FILETYPE))')

package: $(PACKAGE_FILENAME)

$(PACKAGE_FILENAME): $(_PACKAGE_TARGETS)
	cd $(OUTPUT_DIR)/$(PACKAGE_DIR) && \
		find . -type f -name '*.$(strip $(PACKAGE_FILETYPE))' -exec zip $(abspath $(PACKAGE_FILENAME)) {} \;

##
## Docker image
##
## Build and push a docker image for building the LaTeX files.
## The image will be tagged docker.izmno.be/kam-chemfig:latest
##
## Usage
##   make docker-image
##   make docker-push
##   make docker-make MAKE_TARGET=_build
docker-image:
	docker build -t $(DOCKER_IMAGE) .
	@for tag in $(DOCKER_TAGS); do \
		docker tag $(DOCKER_IMAGE) $$tag; \
	done

docker-push: docker-image
	@for tag in $(DOCKER_TAGS); do \
		docker tag $(DOCKER_IMAGE) $$tag; \
		docker push $$tag; \
	done

MAKE_TARGET ?= build

docker-make:
	docker run --rm \
		-v $(PWD):/workdir \
		-w /workdir \
		$(DOCKER_FULL_IMAGE) make $(MAKE_TARGET)
