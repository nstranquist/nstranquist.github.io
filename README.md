# nstranquist.github.io

Public catalog site for [nstranquist.github.io](https://nstranquist.github.io/).

This is the GitHub Pages user site. It is not the profile README — that lives
in [`nstranquist/nstranquist`](https://github.com/nstranquist/nstranquist) and
renders on [github.com/nstranquist](https://github.com/nstranquist).

`catalog.yaml` is the source of truth. `tools/render` writes `index.html` and
`404.html`. Featured products are showcase-ready public extracts only.

```sh
make verify
make serve   # http://127.0.0.1:8766/
```
