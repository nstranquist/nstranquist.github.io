# Catalog and public proof

## Sub-features

- The page lists the reviewed featured products in catalog order.
- A visitor can open and close each product write-up.
- Source, release, and live-demo links open the stated public proof.
- The page lists additional public repositories without presenting an untagged repository as a release.
- The page remains readable when JavaScript is unavailable.

## How to get to it (user POV)

1. Open the home page.
2. Select **Catalog** in the primary navigation.
3. Select a product name.
4. Use **Source**, **Release**, or **Live demo** in the open write-up.
5. Continue to **Also on GitHub** for other public repositories.

## Driving it with ndev browser

Preconditions: Run `make serve` from the repository root and use the managed browser session named `nstranquist-profile-verify`.

1. Open `http://127.0.0.1:8766/`.
2. Confirm that the page heading is `Nico Stranquist` and the Catalog table is visible.
3. Select the `docs-puller` catalog toggle.
4. Confirm that its detail region is visible and contains Source, Release, and Live demo links.
5. Select the toggle again and confirm that the detail region is hidden.
6. Save the browser evidence to `docs/verify/profile-site/evidence/catalog-local.json`.

## Gotchas

- `catalog.yaml` is the source. Run `make render`; do not edit generated `index.html` by hand.
- External links open a new tab. The verification run must return to the original tab before it continues.
- The browser run proves interaction. It does not prove adoption or a real visitor.
