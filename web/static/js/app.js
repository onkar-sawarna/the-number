(function () {
  function ink() {
    return document.documentElement.classList.contains("dark") ? "#e4e4e7" : "#14181f";
  }
  function muted() {
    return document.documentElement.classList.contains("dark") ? "rgba(228,228,231,0.35)" : "rgba(20,24,31,0.15)";
  }
  function gridColor() {
    return document.documentElement.classList.contains("dark") ? "rgba(255,255,255,0.06)" : "rgba(20,24,31,0.06)";
  }

  function lineOptions() {
    return {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { labels: { color: ink() } } },
      scales: {
        x: { ticks: { color: ink() }, grid: { color: gridColor() } },
        y: { ticks: { color: ink(), callback: function (v) { return Calc.formatINR(v); } }, grid: { color: gridColor() } },
      },
    };
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

  window.fmtINR = function (n) {
    return Calc.formatINR(Number(n) || 0);
  };

  window.themeCtl = function () {
    return {
      toggle: function () {
        var dark = !document.documentElement.classList.contains("dark");
        document.documentElement.classList.toggle("dark", dark);
        document.cookie = "theme=" + (dark ? "dark" : "light") + "; path=/; max-age=31536000; samesite=lax";
      },
    };
  };

  window.fireCalc = function () {
    return {
      age: 30,
      annualExpenses: 1200000,
      currentSavings: 1500000,
      monthlySavings: 50000,
      expectedReturn: 11,
      inflation: 6,
      swr: 4,
      result: Calc.fire({
        age: 30,
        annualExpenses: 1200000,
        currentSavings: 1500000,
        monthlySavings: 50000,
        expectedReturn: 11,
        inflation: 6,
        swr: 4,
      }),
      _chart: null,
      init: function () {
        this.recalc();
      },
      recalc: function () {
        this.result = Calc.fire({
          age: Number(this.age) || 0,
          annualExpenses: Number(this.annualExpenses) || 0,
          currentSavings: Number(this.currentSavings) || 0,
          monthlySavings: Number(this.monthlySavings) || 0,
          expectedReturn: Number(this.expectedReturn) || 0,
          inflation: Number(this.inflation) || 0,
          swr: Number(this.swr) || 0,
        });
        var self = this;
        this.$nextTick(function () {
          self.draw();
        });
      },
      yearsCopy: function () {
        if (this.result.reachesFire && this.result.years === 0) return "Already independent on these assumptions.";
        if (!this.result.reachesFire) return "Does not reach FIRE within 80 years.";
        return this.result.years.toFixed(1) + " years · FI around age " + this.result.fiAge;
      },
      draw: function () {
        var el = document.getElementById("fire-chart");
        var labels = this.result.chart.map(function (p) { return String(p.age); });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Corpus", data: this.result.chart.map(function (p) { return p.corpus; }), borderColor: "#0f766e", backgroundColor: "rgba(15,118,110,0.12)", fill: true, tension: 0.25 },
          { label: "FIRE number", data: this.result.chart.map(function (p) { return p.target; }), borderColor: ink(), borderDash: [6, 4], tension: 0.2, pointRadius: 0 },
        ]);
      },
    };
  };

  window.sipCalc = function () {
    return {
      monthly: 10000,
      existing: 0,
      expectedReturn: 12,
      years: 15,
      result: Calc.sip({ monthly: 10000, existing: 0, expectedReturn: 12, years: 15 }),
      _chart: null,
      init: function () { this.recalc(); },
      recalc: function () {
        this.result = Calc.sip({
          monthly: Number(this.monthly) || 0,
          existing: Number(this.existing) || 0,
          expectedReturn: Number(this.expectedReturn) || 0,
          years: Number(this.years) || 0,
        });
        var self = this;
        this.$nextTick(function () { self.draw(); });
      },
      draw: function () {
        var el = document.getElementById("sip-chart");
        var labels = this.result.chart.map(function (p) { return "Y" + p.year; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Invested", data: this.result.chart.map(function (p) { return p.invested; }), borderColor: muted(), tension: 0.2, pointRadius: 0 },
          { label: "Future value", data: this.result.chart.map(function (p) { return p.fv; }), borderColor: "#0f766e", backgroundColor: "rgba(15,118,110,0.12)", fill: true, tension: 0.25 },
        ]);
      },
    };
  };

  window.emiCalc = function () {
    return {
      principal: 2500000,
      annualRate: 8.5,
      years: 20,
      extraMonthly: 0,
      lump: 0,
      result: Calc.emi({ principal: 2500000, annualRate: 8.5, years: 20, extraMonthly: 0, lump: 0 }),
      _chart: null,
      init: function () { this.recalc(); },
      recalc: function () {
        this.result = Calc.emi({
          principal: Number(this.principal) || 0,
          annualRate: Number(this.annualRate) || 0,
          years: Number(this.years) || 0,
          extraMonthly: Number(this.extraMonthly) || 0,
          lump: Number(this.lump) || 0,
        });
        var self = this;
        this.$nextTick(function () { self.draw(); });
      },
      draw: function () {
        var el = document.getElementById("emi-chart");
        var labels = this.result.chart.map(function (p) { return "Y" + p.year; });
        this._chart = upsertLine(el, this._chart, labels, [
          { label: "Scheduled balance", data: this.result.chart.map(function (p) { return p.balance; }), borderColor: ink(), tension: 0.2, pointRadius: 0 },
          { label: "With prepayment", data: this.result.chart.map(function (p) { return p.prepaidBalance; }), borderColor: "#0f766e", tension: 0.2, pointRadius: 0 },
        ]);
      },
    };
  };

  window.emergencyCalc = function () {
    return {
      monthlyEssentials: 60000,
      monthsCover: 6,
      currentBuffer: 80000,
      monthlyTopup: 15000,
      parkingReturn: 6,
      result: Calc.emergency({ monthlyEssentials: 60000, monthsCover: 6, currentBuffer: 80000, monthlyTopup: 15000, parkingReturn: 6 }),
      _chart: null,
      init: function () { this.recalc(); },
      recalc: function () {
        this.result = Calc.emergency({
          monthlyEssentials: Number(this.monthlyEssentials) || 0,
          monthsCover: Number(this.monthsCover) || 0,
          currentBuffer: Number(this.currentBuffer) || 0,
          monthlyTopup: Number(this.monthlyTopup) || 0,
          parkingReturn: Number(this.parkingReturn) || 0,
        });
        var self = this;
        this.$nextTick(function () { self.draw(); });
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
          { label: "Buffer", data: this.result.chart.map(function (p) { return p.balance; }), borderColor: "#0f766e", fill: true, backgroundColor: "rgba(15,118,110,0.12)", tension: 0.2 },
          { label: "Target", data: this.result.chart.map(function (p) { return p.target; }), borderColor: ink(), borderDash: [6, 4], pointRadius: 0 },
        ]);
      },
    };
  };

  window.budgetCalc = function () {
    return {
      income: 120000,
      needs: 60000,
      wants: 25000,
      savings: 25000,
      result: Calc.budget({ income: 120000, needs: 60000, wants: 25000, savings: 25000 }),
      _chart: null,
      init: function () { this.recalc(); },
      recalc: function () {
        this.result = Calc.budget({
          income: Number(this.income) || 0,
          needs: Number(this.needs) || 0,
          wants: Number(this.wants) || 0,
          savings: Number(this.savings) || 0,
        });
        var self = this;
        this.$nextTick(function () { self.draw(); });
      },
      budgetCopy: function () {
        if (this.result.overspent) return "Overspent: actuals exceed take-home.";
        return "Unallocated " + Calc.formatINR(this.result.unallocated) + ".";
      },
      draw: function () {
        var el = document.getElementById("budget-chart");
        if (!el || typeof Chart === "undefined") return;
        var data = {
          labels: ["Needs", "Wants", "Savings"],
          datasets: [
            { label: "Target", data: [this.result.targetNeeds, this.result.targetWants, this.result.targetSavings], backgroundColor: "rgba(15,118,110,0.35)" },
            { label: "Actual", data: [this.needs, this.wants, this.savings], backgroundColor: "#0f766e" },
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
            plugins: { legend: { labels: { color: ink() } } },
            scales: {
              x: { ticks: { color: ink() }, grid: { display: false } },
              y: { ticks: { color: ink(), callback: function (v) { return Calc.formatINR(v); } }, grid: { color: gridColor() } },
            },
          },
        });
      },
    };
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
    var colors = ["#0f766e", "#0e7490", "#047857", "#57534e", "#b45309", "#4338ca", "#a8a29e"];
    canvas._chart = new Chart(canvas, {
      type: "doughnut",
      data: {
        labels: labels,
        datasets: [{ data: values, backgroundColor: colors.slice(0, values.length), borderWidth: 0 }],
      },
      options: {
        plugins: { legend: { position: "bottom", labels: { color: ink(), boxWidth: 10 } } },
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
})();
