(function () {
  function isDark() {
    return document.documentElement.classList.contains("dark");
  }
  function ink() {
    return isDark() ? "#f4f4f5" : "#0f1b33";
  }
  function gridColor() {
    return isDark() ? "rgba(244,244,245,0.32)" : "rgba(15,27,51,0.18)";
  }
  function brandLine() {
    return isDark() ? "#7dd3fc" : "#1d4ed8";
  }
  function brandFill() {
    return isDark() ? "rgba(125,211,252,0.30)" : "rgba(29,78,216,0.18)";
  }
  function accentLine() {
    return isDark() ? "#fbbf24" : "#b45309";
  }
  function compareLine() {
    return isDark() ? "#fdba74" : "#c2410c";
  }
  function potStroke(i) {
    var colors = isDark()
      ? ["#7dd3fc", "#fdba74", "#86efac", "#c4b5fd", "#f9a8d4", "#fbbf24", "#5eead4", "#fb7185"]
      : ["#1d4ed8", "#c2410c", "#15803d", "#7c3aed", "#be123c", "#b45309", "#0f766e", "#e11d48"];
    return colors[i % colors.length];
  }
  function potFill(i) {
    var colors = isDark()
      ? ["rgba(125,211,252,0.42)", "rgba(253,186,116,0.38)", "rgba(134,239,172,0.38)", "rgba(196,181,253,0.40)", "rgba(249,168,212,0.38)", "rgba(251,191,36,0.36)", "rgba(94,234,212,0.36)", "rgba(251,113,133,0.36)"]
      : ["rgba(29,78,216,0.28)", "rgba(194,65,12,0.24)", "rgba(21,128,61,0.24)", "rgba(124,58,237,0.24)", "rgba(190,18,60,0.22)", "rgba(180,83,9,0.24)", "rgba(15,118,110,0.24)", "rgba(225,29,72,0.22)"];
    return colors[i % colors.length];
  }
  function targetLine() {
    return isDark() ? "#fafafa" : "#111827";
  }

  function isNarrow() {
    return window.matchMedia("(max-width: 640px)").matches;
  }

  function tickFont() {
    return { size: isNarrow() ? 10 : 11 };
  }

  var chartAnimateMs = 0;

  function prefersReducedMotion() {
    try {
      return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    } catch (e) {
      return false;
    }
  }

  function motionOpts() {
    if (!chartAnimateMs) return false;
    return { duration: chartAnimateMs, easing: "easeOutQuart" };
  }

  function withChartMotion(fn) {
    chartAnimateMs = prefersReducedMotion() ? 0 : 780;
    try {
      fn();
    } finally {
      chartAnimateMs = 0;
    }
  }

  function lineOptions() {
    var narrow = isNarrow();
    applyChartTheme();
    return {
      responsive: true,
      maintainAspectRatio: false,
      animation: motionOpts(),
      color: ink(),
      interaction: { mode: "index", intersect: false },
      elements: { line: { borderWidth: 3 }, point: { radius: narrow ? 0 : 2, hoverRadius: 5 } },
      plugins: {
        legend: {
          labels: { color: ink(), boxWidth: 12, font: { size: narrow ? 11 : 13, weight: "600" }, padding: narrow ? 8 : 12 },
        },
        tooltip: {
          backgroundColor: isDark() ? "#18181b" : "#ffffff",
          titleColor: ink(),
          bodyColor: isDark() ? "#e4e4e7" : "#1f2937",
          borderColor: isDark() ? "rgba(255,255,255,0.25)" : "rgba(15,27,51,0.15)",
          borderWidth: 1,
        },
      },
      scales: {
        x: {
          ticks: { color: ink(), maxTicksLimit: narrow ? 5 : 10, maxRotation: 0, autoSkip: true, font: tickFont() },
          grid: { color: gridColor() },
        },
        y: {
          ticks: {
            color: ink(),
            maxTicksLimit: narrow ? 5 : 8,
            font: tickFont(),
            callback: function (v) { return window.fmtINR(v); },
          },
          grid: { color: gridColor() },
        },
      },
    };
  }

  function groupedBarOptions() {
    var opts = lineOptions();
    opts.elements = { bar: { borderWidth: 0 } };
    opts.scales.x.stacked = false;
    opts.scales.x.grid = { display: false };
    opts.scales.y.stacked = false;
    opts.plugins.tooltip.callbacks = {
      label: function (ctx) {
        return ctx.dataset.label + ": " + window.fmtINR(ctx.parsed.y);
      },
    };
    return opts;
  }

  function parseCompactMoney(raw) {
    var s = String(raw || "").trim().replace(/,/g, "").replace(/[₹$]/g, "");
    if (!s) return 0;
    var m = s.match(/^(-?\d*\.?\d+)\s*(crores?|crs?|lakhs?|lacs?|l|k|m)?$/i);
    if (!m) {
      var fallback = parseFloat(s);
      return isFinite(fallback) ? fallback : NaN;
    }
    var n = parseFloat(m[1]);
    if (!isFinite(n)) return NaN;
    var u = (m[2] || "").toLowerCase();
    if (u.indexOf("cr") === 0) return n * 1e7;
    if (u === "l" || u.indexOf("lac") === 0 || u.indexOf("lakh") === 0) return n * 1e5;
    if (u === "k") return n * 1e3;
    if (u === "m") return n * 1e6;
    return n;
  }

  function isBareAmount(raw) {
    var s = String(raw || "").trim().replace(/,/g, "").replace(/[₹$]/g, "");
    return /^-?\d*\.?\d+$/.test(s);
  }

  // India annual/corpus boxes show ₹12 L. Typing 12 (or 6, 20) means lakhs, not rupees.
  function inferIndiaLakhs(field, n, raw) {
    if (!isBareAmount(raw) || n <= 0 || n >= 1000) return n;
    if (field === "annualExpenses" || field === "currentSavings") return n * 1e5;
    return n;
  }

  function addDirty(ctl) {
    ctl.dirty = false;
    ctl._moneyFocus = "";
    ctl.markDirty = function () {
      this.dirty = true;
    };
    ctl.moneyBox = function (field) {
      var n = Math.round(Number(this[field]) || 0);
      if (n < 0) n = 0;
      return String(n);
    };
    ctl.lang = "en";
    ctl.t = function (key, vars) {
      void this.lang;
      return window.t(key, vars);
    };
    ctl.moneyHint = function () {
      return this.t(this.region === "us" ? "money_hint_us" : "money_hint_in");
    };
    ctl.commitMoney = function (field, ev, scroll) {
      var n = parseCompactMoney(ev.target.value);
      if (!isFinite(n) || n < 0) n = 0;
      if (this.region !== "us") n = inferIndiaLakhs(field, n, ev.target.value);
      var min = this.rangeMin ? this.rangeMin(field) : 0;
      var max = this.rangeMax ? this.rangeMax(field) : 1e12;
      if (n < min) n = min;
      if (n > max) n = max;
      this[field] = n;
      this._moneyFocus = "";
      ev.target.value = this.moneyBox(field);
      this.recalc(scroll);
    };
    ctl.revealed = false;
    ctl.working = false;
    ctl.justCalculated = false;
    ctl._calcTimer = 0;
    ctl._holdDraw = false;
    ctl.resultRevealClass = function () {
      var c = [];
      if (this.working) c.push("is-working");
      if (this.justCalculated) c.push("calc-result-in");
      return c.join(" ");
    };
    var recalc = ctl.recalc;
    ctl.recalc = function (scroll) {
      this.dirty = false;
      if (scroll) {
        runCalculate(this, recalc);
        return;
      }
      recalc.call(this);
    };
    return ctl;
  }

  function runCalculate(self, inner) {
    if (self.working) return;
    if (self._calcTimer) {
      clearTimeout(self._calcTimer);
      self._calcTimer = 0;
    }
    self.working = true;
    var first = !self.revealed;
    var wait = prefersReducedMotion() ? 0 : 480;
    var finish = function () {
      self._calcTimer = 0;
      self._holdDraw = true;
      inner.call(self);
      self._holdDraw = false;
      self.working = false;
      self.revealed = true;
      self.justCalculated = first;
      if (first && !prefersReducedMotion()) {
        setTimeout(function () {
          self.justCalculated = false;
        }, 500);
      } else {
        self.justCalculated = false;
      }
      requestAnimationFrame(function () {
        requestAnimationFrame(function () {
          withChartMotion(function () {
            if (typeof self.draw === "function") self.draw();
          });
          scrollToResult();
        });
      });
    };
    if (wait === 0) finish();
    else self._calcTimer = setTimeout(finish, wait);
  }

  function scrollToResult() {
    var el = document.getElementById("calc-result");
    if (!el) return;
    var header = document.getElementById("app-chrome");
    var topChrome = (header ? header.offsetHeight : 0) + 8;
    var stat = document.getElementById("calc-sticky-stat");
    if (stat && stat.offsetParent !== null) topChrome += stat.offsetHeight;
    var bottomChrome = 16;
    var dock = document.querySelector(".dock");
    if (dock && window.getComputedStyle(dock).display !== "none") {
      bottomChrome += dock.offsetHeight;
    }
    var avail = window.innerHeight - topChrome - bottomChrome;
    var y = el.getBoundingClientRect().top + window.scrollY;
    var h = el.offsetHeight;
    var target = y - topChrome;
    if (h < avail) target = y - topChrome - Math.min(24, (avail - h) / 2);
    window.scrollTo({ top: Math.max(0, target), behavior: "smooth" });
    requestAnimationFrame(function () {
      resizeCharts();
    });
  }

  function scheduleDraw(self) {
    if (self._holdDraw) return;
    if (self._raf) return;
    self._raf = requestAnimationFrame(function () {
      self._raf = 0;
      self.draw();
    });
  }

  function syncHeaderH() {
    var el = document.getElementById("app-chrome");
    if (!el) return;
    document.documentElement.style.setProperty("--app-header-h", el.offsetHeight + "px");
  }

  function applyChartTheme() {
    if (typeof Chart === "undefined") return;
    Chart.defaults.color = ink();
    Chart.defaults.borderColor = gridColor();
  }

  function stackedFireOptions() {
    var opts = lineOptions();
    opts.scales.y.stacked = true;
    opts.elements.line.borderWidth = 2;
    opts.elements.point.radius = 0;
    opts.elements.point.hoverRadius = 4;
    return opts;
  }

  function upsertLine(el, chart, labels, datasets, options) {
    if (!el || typeof Chart === "undefined") return null;
    var existing = (typeof Chart.getChart === "function" && Chart.getChart(el)) || chart;
    if (existing) existing.destroy();
    return new Chart(el, { type: "line", data: { labels: labels, datasets: datasets }, options: options || lineOptions() });
  }

  function potsSlopeOptions() {
    var opts = lineOptions();
    opts.elements = { line: { borderWidth: 3 }, point: { radius: 5, hoverRadius: 7 } };
    opts.scales.x.ticks.maxTicksLimit = 2;
    opts.scales.x.grid = { display: false };
    opts.plugins.tooltip.callbacks = {
      label: function (ctx) {
        return ctx.dataset.label + ": " + window.fmtINR(ctx.parsed.y);
      },
    };
    return opts;
  }

  function upsertPotsLine(el, chart, labels, datasets) {
    if (!el || typeof Chart === "undefined") return null;
    var existing = (typeof Chart.getChart === "function" && Chart.getChart(el)) || chart;
    if (existing) existing.destroy();
    return new Chart(el, { type: "line", data: { labels: labels, datasets: datasets }, options: potsSlopeOptions() });
  }

  function upsertBar(el, chart, labels, datasets) {
    if (!el || typeof Chart === "undefined") return null;
    var existing = (typeof Chart.getChart === "function" && Chart.getChart(el)) || chart;
    if (existing) existing.destroy();
    return new Chart(el, { type: "bar", data: { labels: labels, datasets: datasets }, options: groupedBarOptions() });
  }

  function resizeCharts() {
    if (typeof Chart === "undefined" || !Chart.getChart) return;
    ["fire-chart", "fire-pots-chart", "sip-chart", "emi-chart", "emergency-chart", "budget-chart"].forEach(function (id) {
      var el = document.getElementById(id);
      var ch = el && Chart.getChart(el);
      if (ch) ch.resize();
    });
  }

  function currentRegion() {
    var m = document.cookie.match(/(?:^|; )region=([^;]*)/);
    var r = m ? decodeURIComponent(m[1]) : (document.documentElement.getAttribute("data-region") || "");
    return r === "us" || r === "usd" ? "us" : "in";
  }

  window.fmtINR = function (n) {
    return Calc.formatMoney(Number(n) || 0, currentRegion());
  };

  window.regionCtl = function () {
    return {
      region: currentRegion(),
      init: function () {
        var self = this;
        window.addEventListener("region-change", function (e) {
          self.region = e.detail === "us" ? "us" : "in";
        });
      },
      set: function (r) {
        this.region = r === "us" ? "us" : "in";
        document.cookie = "region=" + this.region + "; path=/; max-age=31536000; samesite=lax";
        document.documentElement.setAttribute("data-region", this.region);
        window.dispatchEvent(new CustomEvent("region-change", { detail: this.region }));
      },
    };
  };

  function setThemeColor(dark) {
    var meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute("content", dark ? "#071018" : "#eef3fb");
  }

  window.themeCtl = function () {
    return {
      dark: document.documentElement.classList.contains("dark"),
      toggle: function () {
        this.dark = !this.dark;
        document.documentElement.classList.toggle("dark", this.dark);
        document.cookie = "theme=" + (this.dark ? "dark" : "light") + "; path=/; max-age=31536000; samesite=lax";
        setThemeColor(this.dark);
        window.dispatchEvent(new Event("theme-change"));
      },
    };
  };

  function bindRegion(self) {
    self.region = currentRegion();
    window.addEventListener("region-change", function (e) {
      var r = e.detail === "us" ? "us" : "in";
      var switched = self.region !== r;
      self.region = r;
      if (switched && typeof self.applyRegion === "function") self.applyRegion(r, true);
      self.recalc();
    });
    window.addEventListener("theme-change", function () {
      scheduleDraw(self);
    });
    self.lang = "en";
  }

  window.fireCalc = function () {
    return addDirty({
      age: 30,
      annualExpenses: 1200000,
      currentSavings: 1500000,
      monthlySavings: 50000,
      expectedReturn: 12,
      inflation: 6,
      swr: 4,
      npsNow: 0,
      npsMonthly: 0,
      ppfNow: 0,
      ppfMonthly: 0,
      epfNow: 0,
      epfMonthly: 0,
      foreignNow: 0,
      foreignMonthly: 0,
      stoppedNow: 0,
      goldNow: 0,
      goldMonthly: 0,
      jewelleryNow: 0,
      stepUp: 10,
      cityTier: 1,
      housing: "rent",
      region: "in",
      moves: [],
      options: [],
      result: Calc.fire({
        age: 30,
        annualExpenses: 1200000,
        currentSavings: 1500000,
        monthlySavings: 50000,
        expectedReturn: 12,
        inflation: 6,
        swr: 4,
        housing: "rent",
        cityTier: 1,
        region: "in",
      }),
      _chart: null,
      _potsChart: null,
      init: function () {
        var handed = this.overlayFireQuery();
        if (!handed) this.applyRegion(currentRegion(), false);
        this.recalc();
        bindRegion(this);
        if (handed) {
          this.revealed = true;
          var self = this;
          requestAnimationFrame(function () {
            requestAnimationFrame(function () {
              withChartMotion(function () {
                if (typeof self.draw === "function") self.draw();
              });
            });
          });
        }
      },
      isIN: function () {
        return this.region !== "us";
      },
      isUS: function () {
        return this.region === "us";
      },
      applyRegion: function (region, reset) {
        this.region = region === "us" ? "us" : "in";
        if (!reset && this.region === "in") return;
        if (this.region === "us") {
          this.age = 30;
          this.annualExpenses = 60000;
          this.currentSavings = 80000;
          this.monthlySavings = 2500;
          this.npsNow = 0;
          this.npsMonthly = 0;
          this.ppfNow = 0;
          this.ppfMonthly = 0;
          this.epfNow = 0;
          this.epfMonthly = 0;
          this.foreignNow = 0;
          this.foreignMonthly = 0;
          this.stoppedNow = 0;
          this.goldNow = 0;
          this.goldMonthly = 0;
          this.jewelleryNow = 0;
          this.cityTier = 1;
          this.housing = "rent";
          this.expectedReturn = 8;
          this.inflation = 3;
          this.swr = 4;
          return;
        }
        this.age = 30;
        this.annualExpenses = 1200000;
        this.currentSavings = 1500000;
        this.monthlySavings = 50000;
        this.npsNow = 0;
        this.npsMonthly = 0;
        this.ppfNow = 0;
        this.ppfMonthly = 0;
        this.epfNow = 0;
        this.epfMonthly = 0;
        this.foreignNow = 0;
        this.foreignMonthly = 0;
        this.stoppedNow = 0;
        this.goldNow = 0;
        this.goldMonthly = 0;
        this.jewelleryNow = 0;
        this.cityTier = 1;
        this.housing = "rent";
        this.expectedReturn = 12;
        this.inflation = 6;
        this.swr = 4;
        this.stepUp = 10;
      },
      rangeMin: function () {
        return 0;
      },
      rangeMax: function (field) {
        var us = this.region === "us";
        if (field === "monthlySavings") return us ? 30000 : 1000000;
        if (field === "annualExpenses") return us ? 400000 : 20000000;
        if (field === "currentSavings") return us ? 5000000 : 100000000;
        if (field === "npsNow" || field === "stoppedNow" || field === "goldNow" || field === "jewelleryNow") return us ? 3000000 : 50000000;
        if (field === "npsMonthly" || field === "goldMonthly") return us ? 15000 : 150000;
        if (field === "epfNow" || field === "foreignNow") return 50000000;
        if (field === "epfMonthly") return 200000;
        if (field === "ppfNow") return 20000000;
        if (field === "ppfMonthly") return 12500;
        if (field === "foreignMonthly") return 500000;
        return us ? 50000 : 1000000;
      },
      rangeStep: function (field) {
        var us = this.region === "us";
        if (field === "monthlySavings" || field === "npsMonthly" || field === "goldMonthly") return us ? 50 : 500;
        if (field === "annualExpenses") return us ? 1000 : 10000;
        if (field === "currentSavings" || field === "npsNow" || field === "stoppedNow" || field === "goldNow" || field === "jewelleryNow") return us ? 1000 : 10000;
        return us ? 100 : 500;
      },
      fireInput: function () {
        return {
          age: Number(this.age) || 0,
          annualExpenses: Number(this.annualExpenses) || 0,
          currentSavings: Number(this.currentSavings) || 0,
          monthlySavings: Number(this.monthlySavings) || 0,
          expectedReturn: Number(this.expectedReturn) || 0,
          inflation: Number(this.inflation) || 0,
          swr: Number(this.swr) || 0,
          npsNow: Number(this.npsNow) || 0,
          npsMonthly: Number(this.npsMonthly) || 0,
          ppfNow: this.region === "us" ? 0 : Number(this.ppfNow) || 0,
          ppfMonthly: this.region === "us" ? 0 : Number(this.ppfMonthly) || 0,
          epfNow: this.region === "us" ? 0 : Number(this.epfNow) || 0,
          epfMonthly: this.region === "us" ? 0 : Number(this.epfMonthly) || 0,
          foreignNow: this.region === "us" ? 0 : Number(this.foreignNow) || 0,
          foreignMonthly: this.region === "us" ? 0 : Number(this.foreignMonthly) || 0,
          stoppedNow: Number(this.stoppedNow) || 0,
          goldNow: Number(this.goldNow) || 0,
          goldMonthly: Number(this.goldMonthly) || 0,
          jewelleryNow: Number(this.jewelleryNow) || 0,
          stepUp: Number(this.stepUp) || 0,
          cityTier: Number(this.cityTier) || 1,
          housing: this.housing || "rent",
          region: this.region || currentRegion(),
        };
      },
      fullFireHref: function () {
        function n(v) {
          return String(Math.round(Number(v) || 0));
        }
        var p = new URLSearchParams();
        p.set("age", n(this.age));
        p.set("expenses", n(this.annualExpenses));
        p.set("parked", n(this.currentSavings));
        p.set("monthly", n(this.monthlySavings));
        p.set("step", n(this.stepUp));
        p.set("ret", n(this.expectedReturn));
        p.set("inf", n(this.inflation));
        p.set("swr", String(Number(this.swr) || 0));
        p.set("nps", n(this.npsNow));
        p.set("npsm", n(this.npsMonthly));
        p.set("ppf", n(this.ppfNow));
        p.set("ppfm", n(this.ppfMonthly));
        p.set("epf", n(this.epfNow));
        p.set("epfm", n(this.epfMonthly));
        p.set("fx", n(this.foreignNow));
        p.set("fxm", n(this.foreignMonthly));
        p.set("stopped", n(this.stoppedNow));
        p.set("gold", n(this.goldNow));
        p.set("goldm", n(this.goldMonthly));
        p.set("jew", n(this.jewelleryNow));
        p.set("tier", String(Number(this.cityTier) || 1));
        p.set("housing", this.housing === "own" || this.housing === "buy" ? this.housing : "rent");
        p.set("region", this.region === "us" ? "us" : "in");
        return "/calculators/fire?" + p.toString();
      },
      overlayFireQuery: function () {
        var q = new URLSearchParams(location.search);
        var hasRegion = q.get("region") === "us" || q.get("region") === "in";
        var hasNum = q.has("age") || q.has("expenses") || q.has("parked") || q.has("monthly") || q.has("step");
        if (!hasRegion && !hasNum) return false;
        var region = hasRegion ? (q.get("region") === "us" ? "us" : "in") : currentRegion();
        if (hasRegion) {
          document.cookie = "region=" + region + "; path=/; max-age=31536000; samesite=lax";
          document.documentElement.setAttribute("data-region", region);
        }
        this.applyRegion(region, true);
        var self = this;
        function take(key, field, min, max) {
          if (!q.has(key)) return;
          var n = parseFloat(String(q.get(key) || "").replace(/,/g, ""));
          if (!isFinite(n) || n < 0) return;
          if (n < min) n = min;
          if (n > max) n = max;
          self[field] = n;
        }
        take("age", "age", 18, 70);
        take("expenses", "annualExpenses", 0, this.rangeMax("annualExpenses"));
        take("parked", "currentSavings", 0, this.rangeMax("currentSavings"));
        take("monthly", "monthlySavings", 0, this.rangeMax("monthlySavings"));
        take("step", "stepUp", 0, 20);
        take("ret", "expectedReturn", 0, 20);
        take("inf", "inflation", 0, 15);
        take("swr", "swr", 2, 8);
        take("nps", "npsNow", 0, this.rangeMax("npsNow"));
        take("npsm", "npsMonthly", 0, this.rangeMax("npsMonthly"));
        take("ppf", "ppfNow", 0, this.rangeMax("ppfNow"));
        take("ppfm", "ppfMonthly", 0, this.rangeMax("ppfMonthly"));
        take("epf", "epfNow", 0, this.rangeMax("epfNow"));
        take("epfm", "epfMonthly", 0, this.rangeMax("epfMonthly"));
        take("fx", "foreignNow", 0, this.rangeMax("foreignNow"));
        take("fxm", "foreignMonthly", 0, this.rangeMax("foreignMonthly"));
        take("stopped", "stoppedNow", 0, this.rangeMax("stoppedNow"));
        take("gold", "goldNow", 0, this.rangeMax("goldNow"));
        take("goldm", "goldMonthly", 0, this.rangeMax("goldMonthly"));
        take("jew", "jewelleryNow", 0, this.rangeMax("jewelleryNow"));
        take("tier", "cityTier", 1, 3);
        this.cityTier = Math.round(Number(this.cityTier) || 1);
        if (this.cityTier < 1) this.cityTier = 1;
        if (this.cityTier > 3) this.cityTier = 3;
        var housing = q.get("housing");
        if (housing === "own" || housing === "rent" || housing === "buy") this.housing = housing;
        return hasNum;
      },
      applyPreset: function (name) {
        this.applyRegion(this.region, true);
        if (this.region === "us") {
          if (name === "start") {
            this.age = 24;
            this.annualExpenses = 36000;
            this.currentSavings = 0;
            this.monthlySavings = 800;
          } else if (name === "ahead") {
            this.age = 35;
            this.annualExpenses = 90000;
            this.currentSavings = 180000;
            this.monthlySavings = 4000;
            this.npsNow = 120000;
          } else {
            this.age = 30;
            this.annualExpenses = 60000;
            this.currentSavings = 80000;
            this.monthlySavings = 2500;
          }
          this.recalc(true);
          return;
        }
        this.npsNow = 0;
        this.npsMonthly = 0;
        this.ppfNow = 0;
        this.ppfMonthly = 0;
        this.epfNow = 0;
        this.epfMonthly = 0;
        this.foreignNow = 0;
        this.foreignMonthly = 0;
        this.stoppedNow = 0;
        this.goldNow = 0;
        this.goldMonthly = 0;
        this.jewelleryNow = 0;
        this.cityTier = 1;
        this.housing = "rent";
        if (name === "start") {
          this.age = 24;
          this.annualExpenses = 480000;
          this.currentSavings = 0;
          this.monthlySavings = 15000;
        } else if (name === "ahead") {
          this.age = 35;
          this.annualExpenses = 1500000;
          this.currentSavings = 4000000;
          this.monthlySavings = 80000;
          this.epfNow = 2500000;
          this.npsNow = 800000;
          this.ppfNow = 700000;
        } else {
          this.age = 30;
          this.annualExpenses = 1200000;
          this.currentSavings = 1500000;
          this.monthlySavings = 50000;
        }
        this.expectedReturn = 12;
        this.inflation = 6;
        this.swr = 4;
        this.recalc(true);
      },
      recalc: function () {
        var inp = this.fireInput();
        this.result = Calc.fire(inp);
        this.moves = Calc.yearMoves(inp);
        this.options = Calc.yearOptions(inp);
        scheduleDraw(this);
      },
      parkedCopy: function () {
        return this.t(this.region === "us" ? "parked_us_copy" : "parked_in_copy");
      },
      goldCopy: function () {
        return this.t("gold_copy");
      },
      jewelleryCopy: function () {
        return this.t("jew_copy");
      },
      investedCopy: function () {
        return this.t(this.region === "us" ? "inv_us" : "inv_in", { rate: Number(this.expectedReturn) || 0 });
      },
      houseCopy: function () {
        var t = Number(this.cityTier) || 1;
        var us = this.region === "us";
        var cost = us
          ? t === 3 ? 220000 : t === 2 ? 400000 : 800000
          : t === 3 ? 4500000 : t === 2 ? 9000000 : 20000000;
        var city = us
          ? t === 3 ? this.t("city_lo") : t === 2 ? this.t("city_mid") : this.t("city_hi")
          : t === 3 ? this.t("city_t3") : t === 2 ? this.t("city_t2") : this.t("city_t1");
        var line = this.t("house_line", { city: city, cost: window.fmtINR(cost) });
        if (this.housing === "buy") return line + this.t("house_buy");
        if (this.housing === "own") return line + this.t("house_own");
        return line + this.t("house_rent");
      },
      fiYear: function () {
        return new Date().getFullYear() + Math.round(Number(this.result.years) || 0);
      },
      crossingYear: function () {
        return new Date().getFullYear() + Math.round(Number(this.result.crossingYears) || 0);
      },
      yearsCopy: function () {
        if (this.result.reachesFire && this.result.years === 0) return this.t("already_indep");
        if (!this.result.reachesFire) return this.t("never_indep");
        return this.t("fi_line", { year: this.fiYear(), age: this.result.fiAge });
      },
      crossingCopy: function () {
        var r = this.result || {};
        if ((Number(r.monthlyIn) || 0) <= 0) {
          if (r.reachesCrossing) return this.t("crossing_already");
          return this.t("crossing_none");
        }
        if (!r.reachesCrossing) return this.t("crossing_never");
        if (r.crossingYears === 0) return this.t("crossing_already");
        if (r.crossingYears <= 1) return this.t("crossing_this");
        var line = { year: this.crossingYear(), age: r.crossingAge };
        if (r.reachesFire && r.years > 0 && r.crossingYears + 1e-9 < r.years) {
          return this.t("crossing_before", line);
        }
        return this.t("crossing_line", line);
      },
      numberLaterCopy: function () {
        return this.t("number_later", { amount: window.fmtINR(this.result.fireNumberLater) });
      },
      stillNeedCopy: function () {
        var r = this.result || {};
        var line = {
          target: window.fmtINR(r.fireNumber),
          have: window.fmtINR(r.spendableNow),
        };
        if ((Number(r.stillNeed) || 0) <= 0.5 && r.reachesFire && r.years === 0) {
          return this.t("still_need_done", line);
        }
        return this.t("still_need", line);
      },
      yearsStickyCopy: function () {
        return this.yearsCopy();
      },
      hookLine: function () {
        if (this.result.reachesFire && this.result.years === 0) return this.t("already_there");
        if (!this.result.reachesFire) return this.t("not_80");
        var y = this.result.years.toFixed(1);
        return y + this.t(Number(y) === 1 ? "year_one" : "year_many");
      },
      hookSub: function () {
        var r = this.result || {};
        if (r.reachesFire && r.years === 0) return this.t("hook_already", { amount: window.fmtINR(r.fireNumber) });
        if (!r.reachesFire) return this.t("hook_never", { amount: window.fmtINR(r.stillNeed) });
        return this.t("hook_fi", { year: this.fiYear(), age: r.fiAge, amount: window.fmtINR(r.stillNeed) });
      },
      formatShift: function (row) {
        var base = this.result || {};
        if (base.reachesFire && base.years === 0 && row.reaches && row.years === 0) return this.t("shift_already");
        if (base.reachesFire && !row.reaches) return this.t("shift_never");
        if (!base.reachesFire && row.reaches) return this.t("shift_reaches");
        if (!base.reachesFire && !row.reaches) return this.t("shift_never");
        var months = Math.round((Number(row.deltaYears) || 0) * 12);
        if (months === 0) return this.t("shift_same");
        if (months <= -12) return this.t("shift_earlier_y", { n: Math.abs(Math.round(row.deltaYears)) });
        if (months < 0) return this.t("shift_earlier_m", { n: Math.abs(months) });
        if (months >= 12) return this.t("shift_later_y", { n: Math.round(row.deltaYears) });
        return this.t("shift_later_m", { n: months });
      },
      moveRows: function () {
        var rows = this.moves || [];
        var self = this;
        return rows.map(function (row) {
          var shift = self.formatShift(row);
          var line;
          if (row.kind === "sip") line = self.t("moves_sip", { amount: window.fmtINR(row.amount), shift: shift });
          else if (row.kind === "return") line = self.t("moves_return", { shift: shift });
          else if (self.housing === "buy") line = self.t("moves_house_rent", { shift: shift });
          else line = self.t("moves_house_buy", { shift: shift });
          return { kind: row.kind, line: line };
        });
      },
      optionLine: function (row) {
        if (row.reaches && row.years === 0) return this.t("optional_already");
        if (!row.reaches) return this.t("optional_never");
        var year = new Date().getFullYear() + Math.round(Number(row.years) || 0);
        return this.t("optional_fi", { year: year, age: row.fiAge });
      },
      optionRows: function () {
        var rows = this.options || [];
        var self = this;
        var none = (Number(this.result && this.result.monthlyIn) || 0) <= 0;
        return rows.map(function (row) {
          var line;
          if (none && (row.kind === "paycut" || row.kind === "parttime")) line = self.t("optional_none");
          else if (row.kind === "paycut") line = self.t("optional_paycut", { line: self.optionLine(row) });
          else if (row.kind === "parttime") line = self.t("optional_part", { line: self.optionLine(row) });
          else line = self.t("optional_pause", { line: self.optionLine(row) });
          return { kind: row.kind, line: line };
        });
      },
      shareText: function () {
        var r = this.result || {};
        if (r.reachesFire && r.years === 0) return this.t("share_already");
        if (r.reachesFire && r.reachesCrossing && r.crossingYears > 0) {
          return this.t("share_fi", { year: this.fiYear(), cross: this.crossingYear() });
        }
        if (r.reachesFire) return this.t("share_fi_only", { year: this.fiYear() });
        if (r.reachesCrossing && r.crossingYears > 0) return this.t("share_cross_only", { year: this.crossingYear() });
        return this.t("share_never");
      },
      shareHeadline: function () {
        var r = this.result || {};
        if (r.reachesFire && r.years === 0) return this.t("share_already");
        if (r.reachesFire) return this.t("share_fi_only", { year: this.fiYear() }).replace(/\.$/, "");
        return this.t("share_never").replace(/\.$/, "");
      },
      shareSubhead: function () {
        var r = this.result || {};
        if (r.reachesCrossing && r.crossingYears > 0 && !(r.reachesFire && r.years === 0)) {
          return this.t("share_cross_only", { year: this.crossingYear() }).replace(/\.$/, "");
        }
        return "";
      },
      drawShareCard: function (canvas) {
        var ctx = canvas.getContext("2d");
        var w = 1200;
        var h = 630;
        var dark = document.documentElement.classList.contains("dark");
        ctx.fillStyle = dark ? "#071018" : "#eef3fb";
        ctx.fillRect(0, 0, w, h);
        ctx.fillStyle = "#2b6cef";
        ctx.fillRect(0, 0, 16, h);
        ctx.fillStyle = dark ? "#7dd3fc" : "#2b6cef";
        ctx.font = "600 28px 'DM Sans', ui-sans-serif, sans-serif";
        ctx.fillText(this.t("share_kicker"), 72, 92);
        ctx.fillStyle = dark ? "#f4f4f5" : "#0f1b33";
        ctx.font = "600 64px Fraunces, ui-serif, Georgia, serif";
        ctx.fillText(this.shareHeadline(), 72, 220);
        var sub = this.shareSubhead();
        if (sub) {
          ctx.font = "500 36px 'DM Sans', ui-sans-serif, sans-serif";
          ctx.fillStyle = dark ? "rgba(244,244,245,0.82)" : "rgba(15,27,51,0.72)";
          ctx.fillText(sub, 72, 300);
        }
        ctx.fillStyle = dark ? "rgba(244,244,245,0.5)" : "rgba(15,27,51,0.4)";
        ctx.font = "400 22px 'DM Sans', ui-sans-serif, sans-serif";
        ctx.fillText(this.t("share_footer"), 72, 560);
      },
      downloadShare: function (blob) {
        var a = document.createElement("a");
        a.href = URL.createObjectURL(blob);
        a.download = "the-number-year.png";
        document.body.appendChild(a);
        a.click();
        a.remove();
        setTimeout(function () {
          URL.revokeObjectURL(a.href);
        }, 1500);
      },
      shareYear: function () {
        var canvas = document.getElementById("fire-share-card");
        if (!canvas || !canvas.getContext) return;
        var self = this;
        var paint = function () {
          self.drawShareCard(canvas);
          var text = self.shareText();
          canvas.toBlob(function (blob) {
            if (!blob) return;
            var file = new File([blob], "the-number-year.png", { type: "image/png" });
            if (navigator.share && navigator.canShare) {
              try {
                if (navigator.canShare({ files: [file] })) {
                  navigator.share({ files: [file], title: self.t("share_kicker"), text: text }).catch(function () {
                    self.downloadShare(blob);
                  });
                  return;
                }
              } catch (e) {}
            }
            self.downloadShare(blob);
          }, "image/png");
        };
        if (document.fonts && document.fonts.ready) document.fonts.ready.then(paint);
        else paint();
      },
      potDefs: function () {
        var us = this.region === "us";
        if (us) {
          return [
            { key: "parked", label: "chart_parked" },
            { key: "nps", label: "chart_retire" },
            { key: "invested", label: "chart_invested_sip" },
            { key: "gold", label: "chart_gold" },
            { key: "jewellery", label: "chart_jew" },
          ];
        }
        return [
          { key: "parked", label: "chart_parked" },
          { key: "nps", label: "chart_nps" },
          { key: "epf", label: "chart_epf" },
          { key: "ppf", label: "chart_ppf" },
          { key: "foreign", label: "chart_foreign" },
          { key: "invested", label: "chart_invested_sip" },
          { key: "gold", label: "chart_gold" },
          { key: "jewellery", label: "chart_jew" },
        ];
      },
      livePots: function () {
        var chart = (this.result && this.result.chart) || [];
        return this.potDefs().filter(function (d) {
          return chart.some(function (p) { return (Number(p[d.key]) || 0) > 1; });
        });
      },
      hasPotLines: function () {
        return this.livePots().length > 0;
      },
      potLaterPoint: function () {
        var chart = (this.result && this.result.chart) || [];
        if (!chart.length) return null;
        for (var i = 0; i < chart.length; i++) {
          if ((Number(chart[i].year) || 0) >= 20) return chart[i];
        }
        return chart[chart.length - 1];
      },
      potLaterLabel: function () {
        var later = this.potLaterPoint();
        var y = later ? Number(later.year) || 0 : 0;
        if (y >= 20) return this.t("pots_later_20");
        if (y > 0) return this.t("pots_later_years", { years: y });
        return this.t("pots_today");
      },
      potsLead: function () {
        var live = this.livePots();
        var chart = (this.result && this.result.chart) || [];
        var later = this.potLaterPoint();
        if (!live.length || !chart[0] || !later) return "";
        var today = chart[0];
        var best = live[0];
        var bestLater = -1;
        live.forEach(function (d) {
          var b = Number(later[d.key]) || 0;
          if (b > bestLater) {
            bestLater = b;
            best = d;
          }
        });
        var a = Number(today[best.key]) || 0;
        var b = Number(later[best.key]) || 0;
        if (b <= 1 && a <= 1) return this.t("pots_lead_flat");
        return this.t("pots_lead", {
          pot: this.t(best.label),
          today: window.fmtINR(a),
          later: window.fmtINR(b),
          when: this.potLaterLabel(),
        });
      },
      potRows: function () {
        var live = this.livePots();
        var chart = (this.result && this.result.chart) || [];
        var later = this.potLaterPoint();
        if (!live.length || !chart[0] || !later) return [];
        var today = chart[0];
        var when = this.potLaterLabel();
        var self = this;
        return live.map(function (d) {
          var a = Number(today[d.key]) || 0;
          var b = Number(later[d.key]) || 0;
          var g = Math.max(0, b - a);
          return {
            name: self.t(d.label),
            line: self.t("pots_span", {
              gain: window.fmtINR(g),
              today: window.fmtINR(a),
              later: window.fmtINR(b),
              when: when,
            }),
          };
        });
      },
      draw: function () {
        var el = document.getElementById("fire-chart");
        var labels = this.result.chart.map(function (p) { return String(p.age); });
        var full = !!document.getElementById("fire-pots-chart");
        var sets = full ? this.fullFireSets() : this.homeFireSets();
        this._chart = upsertLine(el, this._chart, labels, sets, full ? stackedFireOptions() : lineOptions());
        this.drawPots();
      },
      homeFireSets: function () {
        var sets = [
          { label: this.t("chart_spend"), data: this.result.chart.map(function (p) { return p.corpus; }), borderColor: brandLine(), backgroundColor: brandFill(), fill: true, tension: 0.25 },
        ];
        if ((Number(this.jewelleryNow) || 0) > 0) {
          sets.push({ label: this.t("chart_nw"), data: this.result.chart.map(function (p) { return p.netWorth; }), borderColor: accentLine(), tension: 0.25, pointRadius: 0 });
        }
        sets.push({ label: this.t("chart_fire"), data: this.result.chart.map(function (p) { return p.target; }), borderColor: targetLine(), borderDash: [6, 4], tension: 0.2, pointRadius: 0 });
        return sets;
      },
      fullFireSets: function () {
        var self = this;
        var spendable = this.livePots().filter(function (d) { return d.key !== "jewellery"; });
        var sets = spendable.map(function (d, i) {
          return {
            label: self.t(d.label),
            data: self.result.chart.map(function (p) { return Number(p[d.key]) || 0; }),
            borderColor: potStroke(i),
            backgroundColor: potFill(i),
            fill: true,
            stack: "pots",
            tension: 0.2,
            pointRadius: 0,
          };
        });
        if (!sets.length) {
          sets.push({
            label: this.t("chart_spend"),
            data: this.result.chart.map(function (p) { return p.corpus; }),
            borderColor: brandLine(),
            backgroundColor: brandFill(),
            fill: true,
            stack: "pots",
            tension: 0.2,
            pointRadius: 0,
          });
        }
        if ((Number(this.jewelleryNow) || 0) > 0) {
          sets.push({
            label: this.t("chart_nw"),
            data: this.result.chart.map(function (p) { return p.netWorth; }),
            borderColor: accentLine(),
            fill: false,
            stack: "nw",
            tension: 0.25,
            pointRadius: 0,
          });
        }
        sets.push({
          label: this.t("chart_fire"),
          data: this.result.chart.map(function (p) { return p.target; }),
          borderColor: targetLine(),
          borderDash: [6, 4],
          fill: false,
          stack: "fire",
          tension: 0.2,
          pointRadius: 0,
        });
        return sets;
      },
      drawPots: function () {
        var el = document.getElementById("fire-pots-chart");
        if (!el) return;
        var live = this.livePots();
        var today = this.result.chart[0];
        var later = this.potLaterPoint();
        if (!live.length || !today || !later) return;
        var self = this;
        var labels = live.map(function (d) { return self.t(d.label); });
        this._potsChart = upsertBar(el, this._potsChart, labels, [
          {
            label: this.t("pots_today"),
            data: live.map(function (d) { return Number(today[d.key]) || 0; }),
            backgroundColor: isDark() ? "#64748b" : "#94a3b8",
            borderWidth: 0,
          },
          {
            label: this.potLaterLabel(),
            data: live.map(function (d) { return Number(later[d.key]) || 0; }),
            backgroundColor: isDark() ? "#7dd3fc" : "#1d4ed8",
            borderWidth: 0,
          },
        ]);
      },
    });
  };

  window.sipCalc = function () {
    return addDirty({
      monthly: 10000,
      existing: 0,
      expectedReturn: 12,
      years: 15,
      region: currentRegion(),
      result: Calc.sip({ monthly: 10000, existing: 0, expectedReturn: 12, years: 15 }),
      _chart: null,
      init: function () { bindRegion(this); this.recalc(); },
      recalc: function () {
        this.result = Calc.sip({
          monthly: Number(this.monthly) || 0,
          existing: Number(this.existing) || 0,
          expectedReturn: Number(this.expectedReturn) || 0,
          years: Number(this.years) || 0,
        });
        scheduleDraw(this);
      },
      draw: function () {
        var el = document.getElementById("sip-chart");
        var labels = this.result.chart.map(function (p) { return "Y" + p.year; });
        var last = labels.length - 1;
        var endDot = labels.map(function (_, i) { return i === last ? 5 : 0; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: this.t("chart_invested"), data: this.result.chart.map(function (p) { return p.invested; }), borderColor: compareLine(), tension: 0.2, pointRadius: 0, fill: false },
          { label: this.t("chart_fv"), data: this.result.chart.map(function (p) { return p.fv; }), borderColor: brandLine(), backgroundColor: brandFill(), fill: "-1", tension: 0.25, pointRadius: endDot, pointBackgroundColor: brandLine(), pointBorderColor: brandLine() },
        ]);
      },
    });
  };

  window.emiCalc = function () {
    return addDirty({
      principal: 2500000,
      annualRate: 8.5,
      years: 20,
      extraMonthly: 0,
      lump: 0,
      region: currentRegion(),
      result: Calc.emi({ principal: 2500000, annualRate: 8.5, years: 20, extraMonthly: 0, lump: 0 }),
      _chart: null,
      init: function () { bindRegion(this); this.recalc(); },
      recalc: function () {
        this.result = Calc.emi({
          principal: Number(this.principal) || 0,
          annualRate: Number(this.annualRate) || 0,
          years: Number(this.years) || 0,
          extraMonthly: Number(this.extraMonthly) || 0,
          lump: Number(this.lump) || 0,
        });
        scheduleDraw(this);
      },
      emiAfterCopy: function () {
        return this.t("emi_after", {
          interest: window.fmtINR(this.result.prepaidInterest),
          pay: window.fmtINR((Number(this.principal) || 0) + (this.result.prepaidInterest || 0)),
          saved: window.fmtINR(this.result.interestSaved),
        });
      },
      draw: function () {
        var el = document.getElementById("emi-chart");
        var labels = this.result.chart.map(function (p) { return "Y" + p.year; });
        var sets = [
          { label: this.t("chart_sched"), data: this.result.chart.map(function (p) { return p.balance; }), borderColor: compareLine(), tension: 0.2, pointRadius: 0, fill: false },
        ];
        var hasPrepay = (Number(this.extraMonthly) || 0) > 0 || (Number(this.lump) || 0) > 0;
        if (hasPrepay) {
          sets.push({ label: this.t("chart_prepay"), data: this.result.chart.map(function (p) { return p.prepaidBalance; }), borderColor: brandLine(), backgroundColor: brandFill(), fill: "-1", tension: 0.2, pointRadius: 0 });
        }
        this._chart = upsertLine(el, this._chart, labels, sets);
      },
    });
  };

  window.emergencyCalc = function () {
    return addDirty({
      monthlyEssentials: 60000,
      monthsCover: 6,
      currentBuffer: 80000,
      monthlyTopup: 15000,
      parkingReturn: 6,
      region: currentRegion(),
      result: Calc.emergency({ monthlyEssentials: 60000, monthsCover: 6, currentBuffer: 80000, monthlyTopup: 15000, parkingReturn: 6 }),
      _chart: null,
      init: function () { bindRegion(this); this.recalc(); },
      recalc: function () {
        this.result = Calc.emergency({
          monthlyEssentials: Number(this.monthlyEssentials) || 0,
          monthsCover: Number(this.monthsCover) || 0,
          currentBuffer: Number(this.currentBuffer) || 0,
          monthlyTopup: Number(this.monthlyTopup) || 0,
          parkingReturn: Number(this.parkingReturn) || 0,
        });
        scheduleDraw(this);
      },
      fillCopy: function () {
        if (this.result.gap === 0) return this.t("fill_met");
        if (!this.result.reaches) return this.t("fill_never");
        return this.t("fill_months", { n: this.result.monthsToFill });
      },
      draw: function () {
        var el = document.getElementById("emergency-chart");
        var labels = this.result.chart.map(function (p) { return "M" + p.month; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: this.t("chart_buffer"), data: this.result.chart.map(function (p) { return p.balance; }), borderColor: brandLine(), fill: true, backgroundColor: brandFill(), tension: 0.2 },
          { label: this.t("chart_target"), data: this.result.chart.map(function (p) { return p.target; }), borderColor: targetLine(), borderDash: [6, 4], pointRadius: 0 },
        ]);
      },
    });
  };

  window.budgetCalc = function () {
    return addDirty({
      income: 120000,
      needs: 60000,
      wants: 25000,
      savings: 25000,
      region: currentRegion(),
      result: Calc.budget({ income: 120000, needs: 60000, wants: 25000, savings: 25000 }),
      _chart: null,
      init: function () { bindRegion(this); this.recalc(); },
      recalc: function () {
        this.result = Calc.budget({
          income: Number(this.income) || 0,
          needs: Number(this.needs) || 0,
          wants: Number(this.wants) || 0,
          savings: Number(this.savings) || 0,
        });
        scheduleDraw(this);
      },
      verdict: function () {
        if (this.result.overspent) return this.t("overspent");
        if (this.result.deltaNeeds > 1) return this.t("needs_over");
        if (this.result.deltaWants > 1) return this.t("wants_over");
        if (this.result.deltaSavings < -1) return this.t("save_under");
        return this.t("on_track");
      },
      savingsLine: function () {
        return this.t("save_line", {
          pct: (this.result.savingsRate || 0).toFixed(1),
          amount: window.fmtINR(this.result.targetSavings),
        });
      },
      bucketLine: function (which) {
        var d = which === "wants" ? this.result.deltaWants : which === "savings" ? this.result.deltaSavings : this.result.deltaNeeds;
        if (Math.abs(d) < 1) {
          return which === "savings" ? this.t("hit_20") : this.t("on_cap");
        }
        if (which === "savings") {
          return this.t(d > 0 ? "more_20" : "short_20", { amount: window.fmtINR(Math.abs(d)) });
        }
        return this.t(d > 0 ? "over_cap" : "under_cap", { amount: window.fmtINR(Math.abs(d)) });
      },
      budgetCopy: function () {
        if (this.result.overspent) return this.t("bud_over");
        if (this.result.unallocated > 1) return this.t("bud_left", { amount: window.fmtINR(this.result.unallocated) });
        if (this.result.unallocated < -1) return this.t("bud_over_sum", { amount: window.fmtINR(-this.result.unallocated) });
        return this.t("bud_ok");
      },
      draw: function () {
        var el = document.getElementById("budget-chart");
        if (!el || typeof Chart === "undefined") return;
        var data = {
          labels: [this.t("chart_needs"), this.t("chart_wants"), this.t("chart_savings")],
          datasets: [
            { label: this.t("chart_target"), data: [this.result.targetNeeds, this.result.targetWants, this.result.targetSavings], backgroundColor: isDark() ? "#d4d4d8" : "#334155" },
            { label: this.t("chart_actual"), data: [this.needs, this.wants, this.savings], backgroundColor: brandLine() },
          ],
        };
        var barOpts = {
          responsive: true,
          maintainAspectRatio: false,
          animation: motionOpts(),
          plugins: {
            legend: { labels: { color: ink(), boxWidth: 10, font: { size: isNarrow() ? 11 : 13, weight: "500" } } },
            tooltip: {
              backgroundColor: isDark() ? "#18181b" : "#ffffff",
              titleColor: ink(),
              bodyColor: isDark() ? "#e4e4e7" : "#1f2937",
              borderColor: isDark() ? "rgba(255,255,255,0.25)" : "rgba(15,27,51,0.15)",
              borderWidth: 1,
            },
          },
          scales: {
            x: { ticks: { color: ink(), font: tickFont() }, grid: { display: false } },
            y: {
              ticks: { color: ink(), maxTicksLimit: isNarrow() ? 5 : 8, font: tickFont(), callback: function (v) { return window.fmtINR(v); } },
              grid: { color: gridColor() },
            },
          },
        };
        if (this._chart) {
          this._chart.data = data;
          this._chart.options = barOpts;
          this._chart.update(chartAnimateMs ? undefined : "none");
          return;
        }
        this._chart = new Chart(el, { type: "bar", data: data, options: barOpts });
      },
    });
  };

  window.dashList = function () {
    return {
      selected: [],
      reset: function () {
        this.selected = [];
      },
      toggle: function (id, ev) {
        id = String(id);
        if (ev.target.checked) {
          if (this.selected.length >= 2) {
            ev.target.checked = false;
            return;
          }
          this.selected.push(id);
        } else {
          this.selected = this.selected.filter(function (x) { return x !== id; });
        }
      },
      compareHref: function () {
        if (this.selected.length !== 2) return "#";
        return "/dashboard/compare?id=" + this.selected[0] + "&id=" + this.selected[1];
      },
      compareClass: function () {
        return this.selected.length === 2 ? "" : "pointer-events-none opacity-40";
      },
    };
  };

  window.renderDonut = function (canvas, labels, values) {
    if (!canvas || typeof Chart === "undefined") return;
    applyChartTheme();
    if (canvas._chart) canvas._chart.destroy();
    var colors = isDark()
      ? ["#93c5fd", "#fbbf24", "#38bdf8", "#c4b5fd", "#f472b6", "#4ade80", "#fb923c"]
      : ["#1d4ed8", "#0f766e", "#0369a1", "#7c3aed", "#be123c", "#15803d", "#c2410c"];
    canvas._chart = new Chart(canvas, {
      type: "doughnut",
      data: {
        labels: labels,
        datasets: [{ data: values, backgroundColor: colors.slice(0, values.length), borderWidth: 0 }],
      },
      options: {
        maintainAspectRatio: true,
        plugins: {
          legend: {
            position: "bottom",
            labels: { color: ink(), boxWidth: 10, font: { size: isNarrow() ? 10 : 12 }, padding: isNarrow() ? 8 : 12 },
          },
        },
      },
    });
  };

  window.guideForm = function () {
    var y = parseInt(new URLSearchParams(location.search).get("years") || "10", 10);
    if (!isFinite(y) || y < 1) y = 10;
    if (y > 50) y = 50;
    return { horizon: String(y) };
  };

  document.addEventListener("htmx:beforeSwap", function (e) {
    if (e.detail.xhr && e.detail.xhr.status === 401) {
      e.detail.shouldSwap = true;
      e.detail.isError = false;
    }
  });

  document.addEventListener("htmx:afterSwap", function (e) {
    var root = e.detail.elt;
    if (!root) return;
    if (window.Alpine) Alpine.initTree(root);
    var wrap = root.hasAttribute && root.hasAttribute("data-donut-values") ? root : root.querySelector && root.querySelector("[data-donut-values]");
    if (wrap) {
      var labels = (wrap.getAttribute("data-donut-labels") || "").split("|");
      var values = (wrap.getAttribute("data-donut-values") || "").split("|").map(Number);
      var canvas = wrap.querySelector("#guidance-donut") || wrap.querySelector("canvas");
      window.renderDonut(canvas, labels, values);
    }
  });

  function redrawDonut() {
    var wrap = document.querySelector("[data-donut-values]");
    if (!wrap) return;
    var labels = (wrap.getAttribute("data-donut-labels") || "").split("|");
    var values = (wrap.getAttribute("data-donut-values") || "").split("|").map(Number);
    var canvas = wrap.querySelector("#guidance-donut") || wrap.querySelector("canvas");
    window.renderDonut(canvas, labels, values);
  }
  window.addEventListener("theme-change", redrawDonut);

  setThemeColor(document.documentElement.classList.contains("dark"));
  syncHeaderH();
  window.addEventListener("resize", function () {
    syncHeaderH();
    resizeCharts();
  });
  if (document.fonts && document.fonts.ready) document.fonts.ready.then(syncHeaderH);
})();
