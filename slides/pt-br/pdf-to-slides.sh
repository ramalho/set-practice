#!/bin/bash
# Extract each page of a PDF into a separate PNG named <prefix>-NN.png
# Usage: ./pdf-to-slides.sh [pdf-file] [output-dir] [dpi]

set -euo pipefail

PDF="${1:-pratica-de-conjuntos-gophercon-br-best.pdf}"
OUTDIR="${2:-slides-png}"
DPI="${3:-72}"
PREFIX="go-slide"

command -v pdftoppm >/dev/null || { echo "pdftoppm not found (brew install poppler)" >&2; exit 1; }
[ -f "$PDF" ] || { echo "no such file: $PDF" >&2; exit 1; }

mkdir -p "$OUTDIR"
rm -f "$OUTDIR/$PREFIX"-*.png

# -png: PNG output; -r: resolution in dpi
pdftoppm -png -r "$DPI" "$PDF" "$OUTDIR/$PREFIX"

# pdftoppm pads the page number to the width of the page count;
# force exactly 2 digits (01, 02, ...) for PDFs with fewer than 100 pages.
for f in "$OUTDIR/$PREFIX"-*.png; do
    n="${f##*-}"; n="${n%.png}"
    printf -v want "%s/%s-%02d.png" "$OUTDIR" "$PREFIX" "$((10#$n))"
    [ "$f" = "$want" ] || mv "$f" "$want"
done

echo "wrote $(ls "$OUTDIR/$PREFIX"-*.png | wc -l | tr -d ' ') files to $OUTDIR/"
