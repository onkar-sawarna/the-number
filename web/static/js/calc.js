(function (global) {
  function nearlyZero(v) {
    return Math.abs(v) < 1e-12;
  }

  function monthlyEffective(annualPct) {
    return Math.pow(1 + annualPct / 100, 1 / 12) - 1;
  }

  function trimZeros(v) {
    var s = v.toFixed(2).replace(/0+$/, "").replace(/\.$/, "");
    return s || "0";
  }

  function indianGroup(n) {
    n = Math.round(Math.abs(n));
    var s = String(n);
    if (s.length <= 3) return s;
    var last3 = s.slice(-3);
    var rest = s.slice(0, -3);
    var groups = [];
    while (rest.length > 2) {
      groups.unshift(rest.slice(-2));
      rest = rest.slice(0, -2);
    }
    if (rest) groups.unshift(rest);
    return groups.join(",") + "," + last3;
  }

  function formatINR(n) {
    var sign = n < 0 ? "-" : "";
    n = Math.abs(n);
    if (n >= 1e7) return sign + "₹" + trimZeros(n / 1e7) + " Cr";
    if (n >= 1e5) return sign + "₹" + trimZeros(n / 1e5) + " L";
    return sign + "₹" + indianGroup(n);
  }

  function fireNumber(annualExpenses, swr) {
    if (swr <= 0) return 0;
    return annualExpenses / (swr / 100);
  }

  function fireChart(input, swr, years, reaches) {
    var horizon = 40;
    if (reaches) {
      horizon = Math.ceil(years) + 5;
      if (horizon < 15) horizon = 15;
      if (horizon > 50) horizon = 50;
    }
    var rm = monthlyEffective(input.expectedReturn);
    var inf = input.inflation / 100;
    var corpus = input.currentSavings;
    var expenses = input.annualExpenses;
    var pts = [];
    for (var y = 0; y <= horizon; y++) {
      pts.push({
        year: y,
        age: input.age + y,
        corpus: corpus,
        target: fireNumber(expenses, swr),
      });
      for (var m = 0; m < 12; m++) {
        corpus = corpus * (1 + rm) + input.monthlySavings;
      }
      expenses *= 1 + inf;
    }
    return pts;
  }

  function fire(input) {
    var swr = input.swr > 0 ? input.swr : 4;
    var n = fireNumber(input.annualExpenses, swr);
    var out = {
      fireNumber: n,
      lean: n * 0.5,
      regular: n,
      fat: n * 2,
      years: 0,
      reachesFire: false,
      fiAge: input.age,
      chart: [],
    };
    if (input.currentSavings >= n) {
      out.reachesFire = true;
      out.chart = fireChart(input, swr, 0, true);
      return out;
    }
    var rm = monthlyEffective(input.expectedReturn);
    var im = monthlyEffective(input.inflation);
    var corpus = input.currentSavings;
    var expenses = input.annualExpenses;
    var maxMonths = 80 * 12;
    for (var m = 1; m <= maxMonths; m++) {
      corpus = corpus * (1 + rm) + input.monthlySavings;
      expenses = expenses * (1 + im);
      if (corpus >= fireNumber(expenses, swr)) {
        out.reachesFire = true;
        out.years = m / 12;
        out.fiAge = input.age + Math.round(out.years);
        out.chart = fireChart(input, swr, out.years, true);
        return out;
      }
    }
    out.chart = fireChart(input, swr, 0, false);
    return out;
  }

  function sipFV(monthly, existing, rm, n) {
    if (n <= 0) return existing;
    if (nearlyZero(rm)) return existing + monthly * n;
    var fvSip = (monthly * (Math.pow(1 + rm, n) - 1)) / rm;
    return fvSip + existing * Math.pow(1 + rm, n);
  }

  function sip(input) {
    var years = Math.max(0, input.years | 0);
    var n = years * 12;
    var rm = monthlyEffective(input.expectedReturn);
    var fv = sipFV(input.monthly, input.existing, rm, n);
    var invested = input.monthly * n + input.existing;
    var chart = [];
    for (var y = 0; y <= years; y++) {
      var nm = y * 12;
      chart.push({
        year: y,
        invested: input.monthly * nm + input.existing,
        fv: sipFV(input.monthly, input.existing, rm, nm),
      });
    }
    return { invested: invested, fv: fv, gain: fv - invested, chart: chart };
  }

  function emiAmount(principal, annualRate, years) {
    var n = years * 12;
    if (n <= 0) return 0;
    var r = annualRate / 12 / 100;
    if (nearlyZero(r)) return principal / n;
    var pow = Math.pow(1 + r, n);
    return (principal * r * pow) / (pow - 1);
  }

  function simulateLoan(principal, r, payment, maxMonths) {
    var sim = { opening: principal, months: 0, interest: 0, monthly: [] };
    if (principal <= 0.01) return sim;
    if (maxMonths < 1) maxMonths = 1;
    var bal = principal;
    var eps = 0.01;
    for (var m = 1; m <= maxMonths; m++) {
      var intM = bal * r;
      var due = bal + intM;
      var pay = payment;
      if (pay <= 0) {
        sim.interest += intM;
        bal = due;
        sim.monthly.push(bal);
        sim.months = m;
        continue;
      }
      if (pay > due) pay = due;
      sim.interest += intM;
      bal = due - pay;
      if (bal <= eps) {
        sim.monthly.push(0);
        sim.months = m;
        return sim;
      }
      sim.monthly.push(bal);
      sim.months = m;
    }
    return sim;
  }

  function balanceAtYear(sim, year) {
    if (year <= 0) return sim.opening;
    var idx = year * 12 - 1;
    if (idx >= sim.monthly.length) {
      if (!sim.monthly.length) return 0;
      var last = sim.monthly[sim.monthly.length - 1];
      return last <= 0.01 ? 0 : last;
    }
    return sim.monthly[idx];
  }

  function emi(input) {
    var years = Math.max(1, input.years | 0);
    var n = years * 12;
    var r = input.annualRate / 12 / 100;
    var emiVal = emiAmount(input.principal, input.annualRate, years);
    var base = simulateLoan(input.principal, r, emiVal, n);
    var prepaidPrincipal = Math.max(0, input.principal - input.lump);
    var prepaid = simulateLoan(prepaidPrincipal, r, emiVal + input.extraMonthly, n);
    var interestSaved = base.interest - prepaid.interest;
    var monthsSaved = base.months - prepaid.months;
    var chart = [];
    for (var y = 0; y <= years; y++) {
      chart.push({
        year: y,
        balance: balanceAtYear(base, y),
        prepaidBalance: balanceAtYear(prepaid, y),
      });
    }
    return {
      emi: emiVal,
      months: base.months,
      totalInterest: base.interest,
      totalPayment: input.principal + base.interest,
      prepaidMonths: prepaid.months,
      prepaidInterest: prepaid.interest,
      interestSaved: interestSaved < 0 ? 0 : interestSaved,
      monthsSaved: monthsSaved < 0 ? 0 : monthsSaved,
      chart: chart,
    };
  }

  function emergency(input) {
    var target = input.monthlyEssentials * input.monthsCover;
    var gap = Math.max(0, target - input.currentBuffer);
    var coverageNow = input.monthlyEssentials > 0 ? input.currentBuffer / input.monthlyEssentials : 0;
    var out = {
      target: target,
      gap: gap,
      coverageNow: coverageNow,
      monthsToFill: 0,
      reaches: gap === 0,
      chart: [],
    };
    var rm = input.parkingReturn / 12 / 100;
    function chart(fillMonths) {
      var horizon = fillMonths > 0 ? fillMonths : 12;
      if (horizon < 12) horizon = 12;
      if (horizon > 40 * 12) horizon = 40 * 12;
      var bal = input.currentBuffer;
      var pts = [{ month: 0, balance: bal, target: target }];
      var step = horizon > 36 ? 12 : 1;
      for (var m = 1; m <= horizon; m++) {
        bal = bal * (1 + rm) + input.monthlyTopup;
        if (m % step === 0 || m === horizon || (fillMonths > 0 && m === fillMonths)) {
          pts.push({ month: m, balance: bal, target: target });
        }
      }
      return pts;
    }
    if (gap === 0) {
      out.chart = chart(0);
      return out;
    }
    var bal = input.currentBuffer;
    var maxM = 40 * 12;
    for (var m = 1; m <= maxM; m++) {
      bal = bal * (1 + rm) + input.monthlyTopup;
      if (bal >= target) {
        out.monthsToFill = m;
        out.reaches = true;
        out.chart = chart(m);
        return out;
      }
    }
    out.chart = chart(maxM);
    return out;
  }

  function budget(input) {
    var targetNeeds = input.income * 0.5;
    var targetWants = input.income * 0.3;
    var targetSavings = input.income * 0.2;
    var unallocated = input.income - input.needs - input.wants - input.savings;
    return {
      targetNeeds: targetNeeds,
      targetWants: targetWants,
      targetSavings: targetSavings,
      deltaNeeds: input.needs - targetNeeds,
      deltaWants: input.wants - targetWants,
      deltaSavings: input.savings - targetSavings,
      unallocated: unallocated,
      overspent: input.needs + input.wants + input.savings > input.income + 0.005,
      savingsRate: input.income > 0 ? (input.savings / input.income) * 100 : 0,
    };
  }

  global.Calc = {
    formatINR: formatINR,
    fire: fire,
    sip: sip,
    emi: emi,
    emergency: emergency,
    budget: budget,
  };
})(window);
