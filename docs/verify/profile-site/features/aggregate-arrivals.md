# Privacy-safe aggregate arrivals

## Sub-features

- The home page sends one aggregate `pageview` event.
- GitHub and LinkedIn profile links send fixed profile-link events.
- Product source, proof, and demo links send a fixed event, surface, and public catalog product ID.
- Global Privacy Control and Do Not Track stop all sends.
- The site sends no cookie, visitor ID, IP address, user agent, referrer, query string, or free-form field.
- A telemetry failure does not block navigation.

## How to get to it (user POV)

1. Open the home page with Global Privacy Control and Do Not Track disabled.
2. Select a tracked public link.
3. Navigation continues even when the aggregate endpoint is unavailable.
4. Enable Global Privacy Control or Do Not Track and reload the page to opt out.

## Driving it with ndev browser

Preconditions: Run `make serve`, run the aggregate Worker on a managed local Cloudflare session when endpoint inspection is needed, and use the managed browser session named `nstranquist-profile-verify`.

1. Open `http://127.0.0.1:8766/`.
2. Confirm that page source contains only reviewed `data-arrival-*` values.
3. Select the `docs-puller` Source link.
4. Confirm that navigation is not blocked.
5. Inspect the managed browser network log and confirm that the telemetry request body contains only `event`, `surface`, and, for a product event, `product`.
6. Save sanitized evidence to `docs/verify/profile-site/evidence/aggregate-arrivals-local.json`.

## Gotchas

- Localhost cannot send production telemetry because the Worker accepts only the canonical GitHub Pages origin.
- A `204` response proves ingestion. It does not prove a unique person, adoption, or revenue.
- Analytics Engine contains aggregates. Do not attempt to reconstruct visitor identity.
