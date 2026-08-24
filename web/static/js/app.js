(function () {
  function ink() {
    return document.documentElement.classList.contains("dark") ? "#e4e4e7" : "#0f1b33";
  }
  function muted() {
    return document.documentElement.classList.contains("dark") ? "rgba(228,228,231,0.35)" : "rgba(20,24,31,0.15)";
  }
  function gridColor() {
    return document.documentElement.classList.contains("dark") ? "rgba(255,255,255,0.06)" : "rgba(20,24,31,0.06)";
  }

  function isNarrow() {
    return window.matchMedia("(max-width: 640px)").matches;
  }

  function tickFont() {
    return { size: isNarrow() ? 10 : 11 };
  }

  function lineOptions() {
    var narrow = isNarrow();
    return {
      responsive: true,
      maintainAspectRatio: false,
      animation: false,
      interaction: { mode: "index", intersect: false },
      elements: { point: { radius: narrow ? 0 : 2, hoverRadius: 4 } },
      plugins: {
        legend: {
          labels: { color: ink(), boxWidth: 8, font: { size: narrow ? 10 : 12 }, padding: narrow ? 8 : 12 },
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
    ctl.moneyHint = function () {
      return this.region === "us" ? "50000 or 50k" : "50000 or 50k";
    };
    ctl.commitMoney = function (field, ev) {
      var n = parseCompactMoney(ev.target.value);
      if (!isFinite(n) || n < 0) n = 0;
      var min = this.rangeMin ? this.rangeMin(field) : 0;
      var max = this.rangeMax ? this.rangeMax(field) : 1e12;
      var range = ev.target.parentElement && ev.target.parentElement.querySelector("input[type=range]");
      if (range) {
        if (!this.rangeMin && range.min !== "") min = Number(range.min);
        if (!this.rangeMax && range.max !== "") max = Number(range.max);
      }
      if (n < min) n = min;
      if (n > max) n = max;
      this[field] = n;
      this._moneyFocus = "";
      ev.target.value = this.moneyBox(field);
      this.markDirty();
    };
    var recalc = ctl.recalc;
    ctl.recalc = function () {
      this.dirty = false;
      recalc.call(this);
    };
    return ctl;
  }

  function scheduleDraw(self) {
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

  function upsertLine(el, chart, labels, datasets) {
    if (!el || typeof Chart === "undefined") return chart;
    if (chart) {
      chart.data.labels = labels;
      chart.data.datasets = datasets;
      chart.options = lineOptions();
      chart.update("none");
      return chart;
    }
    return new Chart(el, { type: "line", data: { labels: labels, datasets: datasets }, options: lineOptions() });
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
      },
    };
  };

  function bindRegion(self) {
    self.region = currentRegion();
    window.addEventListener("region-change", function (e) {
      var r = e.detail === "us" ? "us" : "in";
      self.region = r;
      if (typeof self.applyRegion === "function") self.applyRegion(r, true);
      self.recalc();
    });
  }

  window.fireCalc = function () {
    return addDirty({
      age: 30,
      annualExpenses: 1200000,
      currentSavings: 1500000,
      monthlySavings: 50000,
      expectedReturn: 11,
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
      cityTier: 1,
      housing: "rent",
      region: "in",
      result: Calc.fire({
        age: 30,
        annualExpenses: 1200000,
        currentSavings: 1500000,
        monthlySavings: 50000,
        expectedReturn: 11,
        inflation: 6,
        swr: 4,
        housing: "rent",
        cityTier: 1,
        region: "in",
      }),
      _chart: null,
      init: function () {
        this.applyRegion(currentRegion(), false);
        this.recalc();
        bindRegion(this);
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
        this.expectedReturn = 11;
        this.inflation = 6;
        this.swr = 4;
      },
      rangeMin: function (field) {
        var us = this.region === "us";
        if (field === "annualExpenses") return us ? 12000 : 120000;
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
          cityTier: Number(this.cityTier) || 1,
          housing: this.housing || "rent",
          region: this.region || currentRegion(),
        };
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
          this.recalc();
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
        this.expectedReturn = 11;
        this.inflation = 6;
        this.swr = 4;
        this.recalc();
      },
      recalc: function () {
        this.result = Calc.fire(this.fireInput());
        scheduleDraw(this);
      },
      parkedCopy: function () {
        if (this.region === "us") {
          return "Educational parking rate 6%. Cash and savings only — no monthly contribution into this pot.";
        }
        return "Educational parking rate 6% (cash, FDs, liquid funds). No monthly SIP here. Old SIPs still sitting in funds go under Invested.";
      },
      goldCopy: function () {
        return "Educational ~8%. Gold funds and coins you could sell. Not jewellery.";
      },
      jewelleryCopy: function () {
        return "Grows at ~8% like gold, so net worth rises. You keep it — it does not shorten years-to-FIRE.";
      },
      investedCopy: function () {
        var r = Number(this.expectedReturn) || 0;
        if (this.region === "us") {
          return "This pot grows at " + r + "% — your expected return. Old contributions still invested, plus monthly contributions still running. Lock-ins, employer match, and tax wrappers are ignored — educational only.";
        }
        return "This pot grows at " + r + "% — your expected return. Old SIPs still in funds, plus SIPs still running. NPS ~9%, EPF ~8.25%, PPF ~7.1% — educational rates, lock-ins ignored. PPF cap is ₹1.5L/year.";
      },
      houseCopy: function () {
        var t = Number(this.cityTier) || 1;
        var us = this.region === "us";
        var cost = us
          ? t === 3 ? 220000 : t === 2 ? 400000 : 800000
          : t === 3 ? 4500000 : t === 2 ? 9000000 : 20000000;
        var label = us
          ? t === 3 ? "lower-cost" : t === 2 ? "mid-cost" : "high-cost"
          : t === 3 ? "tier-3" : t === 2 ? "tier-2" : "tier-1";
        var line = "Indicative modest house in a " + label + " city: " + window.fmtINR(cost) + ".";
        if (this.housing === "buy") return line + " Added to the number, then inflated with expenses. Keep today’s rent inside annual expenses until you buy; after that you can drop rent from spend.";
        if (this.housing === "own") return line + " Not added — you already live there. Keep the mortgage out of expenses if the home is paid.";
        return line + " Not added. You’ll keep renting, so put rent inside annual expenses.";
      },
      yearsCopy: function () {
        if (this.result.reachesFire && this.result.years === 0) return "Already independent on these assumptions.";
        if (!this.result.reachesFire) return "Does not reach FIRE within 80 years.";
        return this.result.years.toFixed(1) + " years · FI around age " + this.result.fiAge;
      },
      hookLine: function () {
        if (this.result.reachesFire && this.result.years === 0) return "Already there.";
        if (!this.result.reachesFire) return "Not in 80 yrs";
        var y = this.result.years.toFixed(1);
        return y + (Number(y) === 1 ? " year" : " years");
      },
      hookSub: function () {
        if (this.result.reachesFire && this.result.years === 0) return "Corpus already covers the number. " + window.fmtINR(this.result.fireNumber) + ".";
        if (!this.result.reachesFire) return "Save more, spend less, or both. Need " + window.fmtINR(this.result.fireNumber) + ".";
        return "FI around age " + this.result.fiAge + " · need " + window.fmtINR(this.result.fireNumber);
      },
      draw: function () {
        var el = document.getElementById("fire-chart");
        var labels = this.result.chart.map(function (p) { return String(p.age); });
        var sets = [
          { label: "Spendable", data: this.result.chart.map(function (p) { return p.corpus; }), borderColor: "#2b6cef", backgroundColor: "rgba(43,108,239,0.14)", fill: true, tension: 0.25 },
        ];
        if ((Number(this.jewelleryNow) || 0) > 0) {
          sets.push({ label: "Net worth", data: this.result.chart.map(function (p) { return p.netWorth; }), borderColor: "#38bdf8", tension: 0.25, pointRadius: 0 });
        }
        sets.push({ label: "FIRE number", data: this.result.chart.map(function (p) { return p.target; }), borderColor: ink(), borderDash: [6, 4], tension: 0.2, pointRadius: 0 });
        this._chart = upsertLine(el, this._chart, labels, sets);
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
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Invested", data: this.result.chart.map(function (p) { return p.invested; }), borderColor: muted(), tension: 0.2, pointRadius: 0 },
          { label: "Future value", data: this.result.chart.map(function (p) { return p.fv; }), borderColor: "#2b6cef", backgroundColor: "rgba(43,108,239,0.14)", fill: true, tension: 0.25 },
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
      draw: function () {
        var el = document.getElementById("emi-chart");
        var labels = this.result.chart.map(function (p) { return "Y" + p.year; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Scheduled balance", data: this.result.chart.map(function (p) { return p.balance; }), borderColor: ink(), tension: 0.2, pointRadius: 0 },
          { label: "With prepayment", data: this.result.chart.map(function (p) { return p.prepaidBalance; }), borderColor: "#2b6cef", tension: 0.2, pointRadius: 0 },
        ]);
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
        if (this.result.gap === 0) return "Target already met.";
        if (!this.result.reaches) return "Does not fill within 40 years on this top-up.";
        return this.result.monthsToFill + " months to fill.";
      },
      draw: function () {
        var el = document.getElementById("emergency-chart");
        var labels = this.result.chart.map(function (p) { return "M" + p.month; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Buffer", data: this.result.chart.map(function (p) { return p.balance; }), borderColor: "#2b6cef", fill: true, backgroundColor: "rgba(43,108,239,0.14)", tension: 0.2 },
          { label: "Target", data: this.result.chart.map(function (p) { return p.target; }), borderColor: ink(), borderDash: [6, 4], pointRadius: 0 },
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
      budgetCopy: function () {
        if (this.result.overspent) return "Overspent: actuals exceed take-home.";
        return "Unallocated " + window.fmtINR(this.result.unallocated) + ".";
      },
      draw: function () {
        var el = document.getElementById("budget-chart");
        if (!el || typeof Chart === "undefined") return;
        var data = {
          labels: ["Needs", "Wants", "Savings"],
          datasets: [
            { label: "Target", data: [this.result.targetNeeds, this.result.targetWants, this.result.targetSavings], backgroundColor: "rgba(43,108,239,0.35)" },
            { label: "Actual", data: [this.needs, this.wants, this.savings], backgroundColor: "#2b6cef" },
          ],
        };
        if (this._chart) {
          this._chart.data = data;
          this._chart.update("none");
          return;
        }
        this._chart = new Chart(el, {
          type: "bar",
          data: data,
          options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: false,
            plugins: { legend: { labels: { color: ink(), boxWidth: 8, font: { size: isNarrow() ? 10 : 12 } } } },
            scales: {
              x: { ticks: { color: ink(), font: tickFont() }, grid: { display: false } },
              y: {
                ticks: { color: ink(), maxTicksLimit: isNarrow() ? 5 : 8, font: tickFont(), callback: function (v) { return window.fmtINR(v); } },
                grid: { color: gridColor() },
              },
            },
          },
        });
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
    if (canvas._chart) canvas._chart.destroy();
    var colors = ["#2b6cef", "#38bdf8", "#1d4ed8", "#0f1b33", "#6366f1", "#7dd3fc", "#94a3b8"];
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

  setThemeColor(document.documentElement.classList.contains("dark"));
  syncHeaderH();
  window.addEventListener("resize", syncHeaderH);
  if (document.fonts && document.fonts.ready) document.fonts.ready.then(syncHeaderH);
})();
