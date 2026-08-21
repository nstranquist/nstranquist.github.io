# Profile site feature map

This map defines the public behavior that must work before `nstranquist.github.io` is published. The generated HTML remains the source artifact.

- [Catalog and public proof](features/catalog-and-public-proof.md)
- [Privacy-safe aggregate arrivals](features/aggregate-arrivals.md)

## Automated production proof

- [Managed-browser scenario](scenarios/canonical-profile-arrival.yaml)
- Browser and NVR evidence belongs in `nicos-tools/nicos-dev/docs/verify/web-products/nstranquist-github-pages/`.
- Run the scenario from the `nicos-tools` repository root after GitHub Pages publishes the source commit.
- Do not commit generated evidence to this repository. Such a commit would start another Pages deployment and make the evidence stale.

Do not add cookies, visitor identifiers, IP addresses, user agents, referrers, or raw URLs to evidence or telemetry.
