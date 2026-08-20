# nstranquist.github.io

GitHub Pages site for [nstranquist.github.io](https://nstranquist.github.io/).

This is the public work page. It is not the profile README — that lives in
[`nstranquist/nstranquist`](https://github.com/nstranquist/nstranquist) and
renders on [github.com/nstranquist](https://github.com/nstranquist).

`catalog.yaml` is the source of truth. `tools/render` writes `index.html` and
`404.html`. Selected work has a GitHub Release. A second list names other
public source, including repos without a tag.

PR previews deploy to Cloudflare Pages and the URL is commented on the PR. Production stays https://nstranquist.github.io/

```sh
make verify
make serve   # http://127.0.0.1:8766/
```
