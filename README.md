# the number

A calm personal-finance workshop for Indian savers. Calculate a FIRE corpus (“the number”), years until independence, SIP future value, EMI with prepayment, an emergency-fund gap, and a 50/30/20 budget. Optional category-level allocation guidance — never a stock, fund, AMC, or ticker.

**Not registered investment advice. For educational purposes only.**

There is no login. Calculators and guidance are public; numbers stay in the browser. Typed amounts update maths and charts with no server round-trip.

Live: [https://number.onkarsawarna.dev](https://number.onkarsawarna.dev).

## Run

```bash
# once
go install github.com/a-h/templ/cmd/templ@v0.2.793
npm install

# generate templates + CSS (CSS is already committed)
templ generate
npx tailwindcss -i web/static/css/input.css -o web/static/css/app.css --minify

# serve
make run
```

Or: `PORT=47321 go run ./cmd/server` after `templ generate`.

Open [http://127.0.0.1:47321](http://127.0.0.1:47321).

## Deploy (Render, free)

The live host is [Render](https://render.com): Docker web service, **Free** instance, no database. The public URL is [https://number.onkarsawarna.dev](https://number.onkarsawarna.dev).

1. Sign in at [dashboard.render.com](https://dashboard.render.com) with GitHub (no card needed for Free).
2. **New → Blueprint** and select `onkar-sawarna/the-number`, or **New → Web Service**, repo `onkar-sawarna/the-number`, runtime **Docker**, instance **Free**, region **Singapore**.
3. Deploy, then in Render add the custom domain `number.onkarsawarna.dev`. On Namecheap, CNAME host `number` → the `*.onrender.com` hostname Render shows for **this** service.

Free instances sleep after ~15 minutes idle; the first request after that can take 30–60 seconds. Do not use `the-number.onrender.com` — that hostname is a different app.

Hot reload (optional): `make install` then `air`.

```bash
go test ./internal/calc ./internal/ai
```

## Environment

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `47321` | HTTP, bound to `0.0.0.0` |
| `OPENAI_API_KEY` | empty | Guidance model (JSON mode) |
| `OPENAI_MODEL` | `gpt-4o-mini` | |
| `ANTHROPIC_API_KEY` | empty | Used only if OpenAI is unset |
| `ANTHROPIC_MODEL` | `claude-3-5-haiku-latest` | |

No API keys are hardcoded. Without keys, guidance uses the on-device heuristic.

## Layout

```
cmd/server/main.go          # entrypoint
internal/calc/              # PURE maths + unit tests
internal/handlers/          # Echo handlers (no FIRE/SIP/EMI formulas)
internal/ai/                # OpenAI / Anthropic + JSON sanitiser
web/templates/              # templ views
web/static/                 # CSS, calc.js, app.js, theme-boot.js
```

Module path: `github.com/thenumber/app`.

## Stack

Go 1.22+, Echo v4, templ, HTMX 2, Alpine.js 3, Tailwind 3.4 (`darkMode: 'class'`), Chart.js 4. No database.

Theme is a cookie (`theme=dark|light`), never `localStorage`. `theme-boot.js` runs in `<head>` to avoid a flash of the wrong theme.

## Disclaimer

This is educational arithmetic on your assumptions. It is not a SEBI-registered advisor. The 4% rule can fail. Taxes are omitted. See `/disclaimer`.
