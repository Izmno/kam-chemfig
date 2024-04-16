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

ENTRYPOINT []
