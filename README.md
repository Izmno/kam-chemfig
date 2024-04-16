# kam-chemfig

Images of chemical compounds for educational purposes. Used by Atheneum Mariakerke

This project includes a docker image with all software dependencies installed. Build the image with
`make docker-image` if you have no access to the repository docker.izmno.be.

With the docker image available, build all images of the tex files in the source directory with `make build`.

There is an included packaging utility. Create a zip archive of the images of any source (sub)directory with `make package PACKAGE_DIR=2024-04-16`.

## Template

```tex
\documentclass{minimal}

% TikZ configuration. Load the DVI 2 SVG TikZ driver
\def\pgfsysdriver{pgfsys-dvisvgm.def}

\usepackage{chemfig}

\begin{document}
\begin{tikzpicture}
    \node at (0,0) {\chemfig{
        % !! chemfig code here
    }};
\end{tikzpicture}
\end{document}
```

## Chemfig cheat sheet

### Bond specification

```
[<angle>,<length>,<departure_atom>,<arrival_atom>,<tikz_code>]
```

### Bond types

see Chemfig manual II.4 (Operation of chemfig -> Diffeent types of bonds)

```

-     Single bond
=     Double bond
~     Triple bond
>     Right cram, plain
<     Left cram, plain
>:    Right cram, dashed
<:    Left cram, dashed
>|    Right cram, hollow
<|    Left cram, hollow
```
