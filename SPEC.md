# SPEC — the number

Personal-finance web app. Not registered investment advice. Every page carries that disclaimer.

## Product

- Know the FIRE corpus (“the number”) and years until independence.
- Live sliders: maths and charts update in the browser with **no server round-trip**.
- Numbers stay in the browser. No accounts, no database.
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
| DB | none |
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
| GET | `/calculators/{fire,sip,emi,emergency,budget}` | Alpine + Chart.js; results update on Calculate |
| GET/POST | `/guidance` | Questionnaire / HTMX fragment |
| GET | `/static/*` | `web/static` |

Bind `0.0.0.0:$PORT`, default **47321**.

## Auth and data

No accounts. Nothing is written except optional cookies: theme (`theme=dark|light`) and region (`region=in|us`). India is ₹; World is USD.

## Calculation rules

Implemented in Go and identically in `web/static/js/calc.js`.

### Currency

India (`region=in`):

- ≥ 1 crore (1e7): `₹X.XX Cr` (trim trailing zeros)
- ≥ 1 lakh (1e5): `₹X.XX L`
- ≥ 1 thousand (1e3): `₹X.XX k`
- else Indian grouping: `₹500`
- Amount boxes show full digits (`50000`). They also accept `50k`, `2L`, `1.5Cr`, and `2M`.

World (`region=us`):

- ≥ 1 million: `$X.XXM`
- ≥ 1,000: `$X.XXk`
- else Western grouping: `$500`

### FIRE

FIRE is the main product. India ₹ vs World $ toggle sits on the calculator (CalcNav and the home playground), not in the site header.

FIRE lifestyle corpus = `annualExpenses / (SWR/100)`. Housing:

- **own:** house add-on is 0
- **rent** (keep renting): house add-on is 0; put rent inside expenses
- **buy** (rent now, buy later): add an indicative modest house, and keep today’s rent in expenses until you buy

India house: tier 1 ₹2 Cr, tier 2 ₹90 L, tier 3 ₹45 L. World: high-cost $800k, mid-cost $400k, lower-cost $220k.

Lean 0.5× lifestyle + house, Regular = total, Fat 2× lifestyle + house.

**India pots:** Starting corpus = parked + gold funds + jewellery + NPS + EPF + PPF + foreign stocks + invested. Parked (cash, FDs, liquid funds) grows at ~6% and takes **no** monthly SIP. Gold funds ~8% (sellable). Jewellery grows at ~8% in net worth; you keep it, so it is ignored for years-to-FIRE. Invested = old SIPs still in funds + monthly SIPs still running, at expected return (same as foreign stocks). NPS ~9%, EPF ~8.25%, PPF ~7.1% (educational).

**World pots:** parked + gold funds + jewellery + retirement account (reuses the NPS fields, grows at expected return) + invested (already invested + monthly contributions). No EPF, PPF, or separate foreign-stock pot. Defaults: expenses $60k, parked $80k, monthly $2.5k, return 8%, inflation 3%, SWR 4%.

Monthly contributions add to the matching pot. Inflation lifts expenses and the house add-on. Years-to-FIRE: monthly loop, max 80 years. Already FI if starting corpus ≥ today’s number.

Chart: yearly snapshots until ~FIRE+5 years (min 15, max 50; 40 if never).

India defaults: age 30, expenses 12,00,000, other savings 15,00,000, monthly 50,000, return 11%, inflation 6%, SWR 4%, city tier 1, housing rent, other pots 0.

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

- **Alpine + Chart.js:** sliders (labels update live), Calculate button for results and charts, dark toggle.
- **HTMX:** guidance POST.

Do not POST on every slider drag.

## Copy constraints

Footer (exact): Not registered investment advice. For educational purposes only.

Guidance banner (exact, persistent): Educational information only, not personalized financial advice. Consult a SEBI-registered investment advisor before investing.

Never recommend specific stocks, funds, AMCs, or tickers in UI or model output.
