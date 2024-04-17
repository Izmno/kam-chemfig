FROM golang:1.22-alpine3.19 as builder

RUN apk add --no-cache \
    make \
    ;

WORKDIR /build
COPY . .

RUN make build-go

FROM reitzig/texlive-base as base

RUN apk add --no-cache \
    bash \
    make \
    imagemagick \
    imagemagick-svg \
    ;

RUN tlmgr install \
    chemfig \
    dvisvgm \
    simplekv \
    ;

RUN tlmgr path add

COPY --from=builder /build/bin/ /usr/local/bin/

ENTRYPOINT []
