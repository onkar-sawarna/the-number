# SPEC — the number

Personal-finance web app. Not registered investment advice. Every page carries that disclaimer.

## Product

- Know the FIRE corpus (“the number”) and years until independence.
- Live inputs: maths and charts update in the browser with **no server round-trip**.
- Numbers stay in the browser. No accounts, no database. Last FIRE inputs can stay in `localStorage` on this device.
- Optional category-level allocation guidance (never product names).
- After results on the full FIRE calculator, the UPI card (Google Pay QR, id, pay, copy) plus a link to `/support`. The home playground stays the number. Not required to see the year. VPA from `SUPPORT_UPI` (default `onkarsawarna-3@okicici`). Remark “Paid for FIRE — the number”.

Calculators and guidance are public. There is no account.

## Stack

| Layer | Choice |
| --- | --- |
| Language | Go 1.22+ |
| HTTP | Echo v4 |
| HTML | templ v0.2.793 |
| Partial updates | HTMX 2.0.4 (`web/static/js/vendor/htmx.min.js`) |
| Client reactivity | Alpine.js 3.14.8 (`web/static/js/vendor/alpine.min.js`) |
| CSS | Tailwind CSS 3.4 CLI (`darkMode: 'class'`), scan `.templ` |
| Charts | Chart.js 4.4.6 (`web/static/js/vendor/chart.umd.min.js`) |
| Fonts | DM Sans + Fraunces latin woff2 (`web/static/fonts/`, `fonts.css`) |
| DB | none |
| Module | `github.com/thenumber/app` |

Scripts and fonts are vendored. `make vendor` re-downloads them (`scripts/vendor.py`). Do not load Alpine, HTMX, Chart.js, or Google Fonts from a CDN at runtime.

templ: do not put `{` in HTML attributes. Use `x-on:click`, not `@click`. Alpine state lives in `web/static/js/app.js`. Load `calc.js` and `app.js` before Alpine (all `defer`).

English copy lives in `web/static/js/i18n.js`. `make generate` writes `web/templates/copy.go` from that file — do not edit `copy.go` by hand. `make og` renders the 1200×630 share card to `web/static/img/og.png`. Layout JSON-LD is a `WebApplication` + `Person` graph (`JSONLDHTML()`); do not put the email address in structured data.

Handlers call `internal/calc` and render templates. They must not contain FIRE/SIP/EMI formulas.

## Routes

| Method | Path | Notes |
| --- | --- | --- |
| GET | `/` | Landing |
| GET | `/about` | Who built this; link to the blog |
| GET | `/support` | Why a UPI; Google Pay QR; keep the workshop running |
| GET | `/crossing` | What the crossing is; jewellery is out; it is not FIRE |
| GET | `/disclaimer` | Limitations |
| POST | `/theme` | Cookie also set from JS |
| GET | `/calculators/{fire,sip,emi,emergency,budget}` | Alpine + Chart.js; results update on Calculate. FIRE accepts a plan query string (home handoff or **Copy this plan**) |
| GET/POST | `/guidance` | Questionnaire / HTMX fragment |
| GET | `/static/*` | `web/static` |
| GET | `/robots.txt` | Allow crawlers; points at sitemap |
| GET | `/sitemap.xml` | Public pages only |
| GET | any other path | HTML 404 in the same layout |

Bind `0.0.0.0:$PORT`, default **47321**.

## Auth and data

No accounts. Nothing is written to a server. Optional cookies: theme (`theme=dark|light`) and region (`region=in|us`). India is ₹; World is USD. The last FIRE inputs stay in `localStorage` (`tn-fire-v1`) on this device after you use the calculator. A query string on `/calculators/fire` reopens the same plan. Clear site data to forget them.

## Calculation rules

Implemented in Go and identically in `web/static/js/calc.js`. `go test ./internal/calc` runs the same fixtures through both (`parity_test.go`); skip only if `node` is missing.

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

FIRE lifestyle corpus = `annualExpenses / (SWR/100)`. That is the **corpus you need**. The result **headline is still to go** (`corpus − spendable pots today`), so parked cash, old SIPs, gold funds, NPS/EPF/PPF, and foreign stocks all change the big figure and years-to-FIRE. Housing:

- **own:** house add-on is 0
- **rent** (keep renting): house add-on is 0; put rent inside expenses
- **buy** (rent now, buy later): add an indicative modest house, and keep today’s rent in expenses until you buy

India house: tier 1 ₹2 Cr, tier 2 ₹90 L, tier 3 ₹45 L. World: high-cost $800k, mid-cost $400k, lower-cost $220k.

Lean FIRE is 20× today’s expenses + house (5% withdrawal). Regular FIRE is `annualExpenses / (SWR/100)` — default 4% / 25×. Fat FIRE is 50× today’s expenses + house.

**India pots:** Starting corpus = parked + gold funds + jewellery + NPS + EPF + PPF + foreign stocks + invested. Parked (cash, FDs, liquid funds) is spendable and grows at ~6%, so it shortens years-to-FIRE, but it takes **no** monthly SIP — you cannot contribute into parked in this calculator. Running SIPs and old SIPs still in funds belong under Invested (expected return). On the full FIRE page, parked and invested sit next to each other. Gold funds ~8% (sellable). Jewellery grows at ~8% in net worth; you keep it, so it is ignored for years-to-FIRE. NPS ~9%, EPF ~8.25%, PPF ~7.1% (educational). Foreign stocks grow at expected return.

**World pots:** parked (no monthly contribution) + gold funds + jewellery + retirement account (reuses the NPS fields, grows at expected return) + invested (already invested + monthly contributions). Parked and invested sit next to each other on the full page. No EPF, PPF, or separate foreign-stock pot. Defaults: expenses $60k, parked $80k, monthly $2.5k, return 8%, inflation 3%, SWR 4%.

Monthly contributions add to the matching pot. Optional yearly SIP step-up raises every monthly pot by that % after each completed year (PPF capped at ₹1.5L/year). Inflation lifts expenses and the house add-on. Years-to-FIRE: monthly loop, max 80 years. Already FI if starting spendable corpus ≥ today’s number.

The **crossing** is the first completed year where spendable pots earn more than you contribute that year (growth = end − start − contributions; jewellery is kept, so it is ignored). Contributions match the monthly adds in the FIRE loop. Crossing is usually before FIRE. No monthly in and a stash already counts as crossed; no monthly in and nothing saved has no crossing. Durable explainer: `/crossing`.

**Coast FIRE** is one line: the earliest you can stop SIPs (`StopAfterMonths`) and still reach FIRE by age 60 — or by your FIRE age if that is later. If spendable pots already grow to the number with no further deposits, you can stop today. If you must contribute until independence, the line says so. It is not a second calculator.

After Calculate, **what moves the year** re-runs FIRE with (1) ₹5,000 / $100 extra monthly SIP, (2) 1% lower expected return, (3) buy vs rent the house. **Salary-optional** re-runs with 40% less going in, half going in, and 24 months of zero deposits then resume (`ContribScale`, `PauseMonths`). The share card is a 1200×630 PNG of the FI year and crossing year — never the corpus. **Copy this plan** copies a `/calculators/fire?...` link that reopens the boxes.

Chart: yearly snapshots until ~FIRE+5 years (min 20, max 50; 40 if never). The FIRE chart labels age on the x-axis, pots versus the corpus you need, and marks the independent year. A key under the chart says what the filled spendable area, the jewellery net-worth line, and the dashed target are.

India defaults: age 30, expenses 12,00,000, other savings 15,00,000, monthly 50,000, return 12%, inflation 6%, SWR 4% (25× expenses), yearly SIP step-up 10%, city tier 1, housing rent, other pots 0. Lean 20×, Fat 50×.

### SIP

Ordinary annuity, month-end contribution, monthly compounding.

`FV_sip = P * ((1+r_m)^n - 1) / r_m` plus existing corpus grown `(1+r_m)^n`. If r ≈ 0, FV = principal sums. Optional yearly step-up raises the monthly SIP by that % after each completed year (same `steppedMonthly` as FIRE). Default step-up 10%.

Defaults: 10,000/mo, 0 existing, 12%, 15 years, 10% yearly step-up.

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

- **Alpine + Chart.js:** amount boxes (type rupees or `50k` / `2L`), number boxes for age / years / percents, Calculate button for a short working beat then results and charts, dark toggle. Preset chips run Calculate. Home **Open the full FIRE calculator** and **Copy this plan** use the same query string. FIRE results include the coast line, year-moves, salary-optional scenarios, and a downloadable/shareable year card.
- **HTMX:** guidance POST.

Do not POST on every keystroke.

## Copy constraints

Footer (exact): Not registered investment advice. For educational purposes only.

Guidance banner (exact, persistent): Educational information only, not personalized financial advice. Consult a SEBI-registered investment advisor before investing.

Never recommend specific stocks, funds, AMCs, or tickers in UI or model output.
