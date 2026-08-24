# SPEC — the number

Personal-finance web app. Not registered investment advice. Every page carries that disclaimer.

## Product

- Know the FIRE corpus (“the number”) and years until independence.
- Live sliders: maths and charts update in the browser with **no server round-trip**.
- Numbers stay in the browser. Login, save, and dashboard are parked for now.
- Optional category-level allocation guidance (never product names).

Calculators and guidance are public. There is no account.

## Stack

| Layer | Choice |
| --- | --- |
| Language | Go 1.22+ |
| HTTP | Echo v4 |
| HTML | templ v0.2.793 |
| Partial updates | HTMX 2 (CDN) |
| Client reactivity | Alpine.js 3 (CDN) |
| CSS | Tailwind CSS 3.4 CLI (`darkMode: 'class'`), scan `.templ` |
| Charts | Chart.js 4 (CDN) |
| DB | SQLite via GORM + `glebarez/sqlite` |
| Auth | `gorilla/sessions` + `echo-contrib/session`; bcrypt cost 12 |
| Module | `github.com/thenumber/app` |

CDNs: htmx 2.0.4, alpinejs 3.14.8, chart.js 4.4.6.

Fonts: Fraunces + DM Sans.

templ: do not put `{` in HTML attributes. Use `x-on:click`, not `@click`. Alpine state lives in `web/static/js/app.js`. Load `calc.js` and `app.js` before Alpine (all `defer`).

Handlers call `internal/calc` and render templates. They must not contain FIRE/SIP/EMI formulas.

## Routes

| Method | Path | Notes |
| --- | --- | --- |
| GET | `/` | Landing |
| GET | `/disclaimer` | Limitations |
| POST | `/theme` | Cookie also set from JS |
| GET | `/calculators/{fire,sip,emi,emergency,budget}` | Live Alpine + Chart.js |
| GET/POST | `/guidance` | Questionnaire / HTMX fragment |
| GET | `/static/*` | `web/static` |

Bind `0.0.0.0:$PORT`, default **47321**.

## Auth

Disabled for now. Handler code remains in the repo but routes are not registered.

## Data

**users:** id, email unique, password_hash, created_at

**scenarios:** id, user_id, kind (`fire|sip|emi|emergency|budget`), title, inputs JSON, outputs JSON, created_at

On save, the server re-runs Go calc from posted form fields and stores both inputs and outputs.

GORM AutoMigrate on boot; `migrations/001_init.sql` kept in sync.

## Calculation rules

Implemented in Go and identically in `web/static/js/calc.js`.

### Currency

- ≥ 1 crore (1e7): `₹X.XX Cr` (trim trailing zeros)
- ≥ 1 lakh (1e5): `₹X.XX L`
- else Indian grouping: `₹12,345`

### FIRE

FIRE number today = `annualExpenses / (SWR/100)`. Lean 0.5×, Regular 1×, Fat 2×.

Monthly rate `r_m = (1 + r)^(1/12) - 1`; same for inflation. Years-to-FIRE: monthly loop, max 80 years. Already FI if starting corpus ≥ today’s FIRE number → 0 years. If never: `reaches_fire = false`.

Chart: yearly snapshots (12 monthly steps then inflate expenses annually) until ~FIRE+5 years (min 15, max 50; 40 if never).

Defaults: age 30, expenses 12,00,000, savings 15,00,000, monthly 50,000, return 11%, inflation 6%, SWR 4%.

### SIP

Ordinary annuity, month-end contribution, monthly compounding.

`FV_sip = P * ((1+r_m)^n - 1) / r_m` plus existing corpus grown `(1+r_m)^n`. If r ≈ 0, FV = principal sums.

Defaults: 10,000/mo, 0 existing, 12%, 15 years.

### EMI

`EMI = P * r * (1+r)^n / ((1+r)^n - 1)` with `r = annual/12/100`, `n = years*12`. If r ≈ 0, `P/n`. Optional extra monthly + day-one lump. Compare interest and months saved.

Defaults: principal 25,00,000, 8.5%, 20 years, extra 0, lump 0.

### Emergency fund

Target = monthly essentials × months of cover. Monthly parking rate = annual/12. Max 40 years.

Defaults: 60,000 expenses, 6 months, 80,000 buffer, 15,000 top-up, 6% parking.

### 50/30/20 budget

Targets: 50% needs, 30% wants, 20% savings of take-home.

Defaults: income 1,20,000; needs 60,000; wants 25,000; savings 25,000.

### Allocation heuristic

Risk bases (large, mid, small, debt, gold, intl, cash):

- conservative: 30, 5, 0, 40, 10, 5, 10
- moderate: 40, 15, 5, 20, 8, 10, 2
- aggressive: 40, 20, 10, 10, 5, 15, 0

Horizon ≤ 3 years: shift 20 (12 if aggressive) from equity (small→mid→large) into 70% debt / 30% cash. Horizon ≤ 7: shift 10 into debt. Age ≥ 50: shift 8 into 60% debt / 40% gold; zero small-cap.

Sleeves sum to 100. Categories only.

## HTMX vs Alpine

- **Alpine + Chart.js:** sliders, instant maths, live charts, dark toggle.
- **HTMX:** guidance POST.

Do not POST on every slider drag.

## Copy constraints

Footer (exact): Not registered investment advice. For educational purposes only.

Guidance banner (exact, persistent): Educational information only, not personalized financial advice. Consult a SEBI-registered investment advisor before investing.

Never recommend specific stocks, funds, AMCs, or tickers in UI or model output.
