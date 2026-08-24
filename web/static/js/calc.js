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
    if (n >= 1e3) return sign + "₹" + trimZeros(n / 1e3) + " k";
    return sign + "₹" + indianGroup(n);
  }

  function westernGroup(n) {
    n = Math.round(Math.abs(n));
    var s = String(n);
    var out = [];
    while (s.length > 3) {
      out.unshift(s.slice(-3));
      s = s.slice(0, -3);
    }
    if (s) out.unshift(s);
    return out.join(",");
  }

  function formatUSD(n) {
    var sign = n < 0 ? "-" : "";
    n = Math.abs(n);
    if (n >= 1e6) return sign + "$" + trimZeros(n / 1e6) + "M";
    if (n >= 1e3) return sign + "$" + trimZeros(n / 1e3) + "k";
    return sign + "$" + westernGroup(n);
  }

  function formatMoney(n, region) {
    if (region === "us" || region === "usd") return formatUSD(n);
    return formatINR(n);
  }

  function isUSD(region) {
    return region === "us" || region === "usd";
  }

  function fireNumber(annualExpenses, swr) {
    if (swr <= 0) return 0;
    return annualExpenses / (swr / 100);
  }

  var RATE_NPS = 9;
  var RATE_PPF = 7.1;
  var RATE_EPF = 8.25;
  var RATE_LIQUID = 6;
  var RATE_GOLD = 8;
  var HOUSE_IN = { 1: 20000000, 2: 9000000, 3: 4500000 };
  var HOUSE_US = { 1: 800000, 2: 400000, 3: 220000 };

  function nz(v) {
    v = Number(v);
    if (!isFinite(v) || v < 0) return 0;
    return v;
  }

  function houseCost(tier, region) {
    var table = isUSD(region) ? HOUSE_US : HOUSE_IN;
    if (tier === 3) return table[3];
    if (tier === 2) return table[2];
    return table[1];
  }

  function houseAdd(input) {
    if (input.housing !== "buy") return 0;
    return houseCost(Number(input.cityTier) || 1, input.region);
  }

  function startPots(input) {
    return {
      general: nz(input.currentSavings),
      nps: nz(input.npsNow),
      ppf: nz(input.ppfNow),
      epf: nz(input.epfNow),
      foreign: nz(input.foreignNow),
      stopped: nz(input.stoppedNow),
      gold: nz(input.goldNow),
      jewellery: nz(input.jewelleryNow),
    };
  }

  function potsSpendable(p) {
    return p.general + p.nps + p.ppf + p.epf + p.foreign + p.stopped + p.gold;
  }

  function potsTotal(p) {
    return potsSpendable(p) + p.jewellery;
  }

  function monthlyIn(input) {
    return nz(input.monthlySavings) + nz(input.npsMonthly) + nz(input.ppfMonthly) + nz(input.epfMonthly) + nz(input.foreignMonthly) + nz(input.goldMonthly);
  }

  function steppedMonthly(base, stepUpPct, month) {
    base = nz(base);
    if (base === 0 || stepUpPct <= 0 || month <= 12) return base;
    var years = Math.floor((month - 1) / 12);
    return base * Math.pow(1 + stepUpPct / 100, years);
  }

  function ppfAt(input, month) {
    var v = steppedMonthly(input.ppfMonthly, input.stepUp, month);
    if (!isUSD(input.region) && v > 12500) return 12500;
    return v;
  }

  function stepPots(p, input, month, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold) {
    p.general = p.general * (1 + rmLiq);
    p.nps = p.nps * (1 + rmNPS) + steppedMonthly(input.npsMonthly, input.stepUp, month);
    p.ppf = p.ppf * (1 + rmPPF) + ppfAt(input, month);
    p.epf = p.epf * (1 + rmEPF) + steppedMonthly(input.epfMonthly, input.stepUp, month);
    p.foreign = p.foreign * (1 + rmEq) + steppedMonthly(input.foreignMonthly, input.stepUp, month);
    p.stopped = p.stopped * (1 + rmEq) + steppedMonthly(input.monthlySavings, input.stepUp, month);
    p.gold = p.gold * (1 + rmGold) + steppedMonthly(input.goldMonthly, input.stepUp, month);
    p.jewellery = p.jewellery * (1 + rmGold);
    return p;
  }

  function fireChart(input, swr, years, reaches) {
    var horizon = 40;
    if (reaches) {
      horizon = Math.ceil(years) + 5;
      if (horizon < 20) horizon = 20;
      if (horizon > 50) horizon = 50;
    }
    var rmLiq = monthlyEffective(RATE_LIQUID);
    var rmEq = monthlyEffective(input.expectedReturn);
    var rmNPS = monthlyEffective(isUSD(input.region) ? input.expectedReturn : RATE_NPS);
    var rmPPF = monthlyEffective(RATE_PPF);
    var rmEPF = monthlyEffective(RATE_EPF);
    var rmGold = monthlyEffective(RATE_GOLD);
    var inf = input.inflation / 100;
    var pots = startPots(input);
    var expenses = nz(input.annualExpenses);
    var house = houseAdd(input);
    var pts = [];
    for (var y = 0; y <= horizon; y++) {
      pts.push({
        year: y,
        age: input.age + y,
        corpus: potsSpendable(pots),
        netWorth: potsTotal(pots),
        target: fireNumber(expenses, swr) + house,
        parked: pots.general,
        nps: pots.nps,
        ppf: pots.ppf,
        epf: pots.epf,
        foreign: pots.foreign,
        invested: pots.stopped,
        gold: pots.gold,
        jewellery: pots.jewellery,
      });
      for (var m = 0; m < 12; m++) {
        pots = stepPots(pots, input, y * 12 + m + 1, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold);
      }
      expenses *= 1 + inf;
      house *= 1 + inf;
    }
    return pts;
  }

  function fire(input) {
    var swr = input.swr > 0 ? input.swr : 4;
    var lifestyle = fireNumber(nz(input.annualExpenses), swr);
    var house = houseAdd(input);
    var n = lifestyle + house;
    var grow = Math.pow(1 + (Number(input.inflation) || 0) / 100, 20);
    var pots = startPots(input);
    var out = {
      fireNumber: n,
      fireNumberLater: n * grow,
      lifestyle: lifestyle,
      lifestyleLater: lifestyle * grow,
      houseAdd: house,
      startingCorpus: potsTotal(pots),
      jewellery: pots.jewellery,
      jewelleryLater: pots.jewellery,
      monthlyIn: monthlyIn(input),
      lean: lifestyle * 0.5 + house,
      regular: n,
      fat: lifestyle * 2 + house,
      years: 0,
      reachesFire: false,
      fiAge: input.age,
      chart: [],
    };
    if (potsSpendable(pots) >= n) {
      out.reachesFire = true;
      out.chart = fireChart(input, swr, 0, true);
      return out;
    }
    var rmLiq = monthlyEffective(RATE_LIQUID);
    var rmEq = monthlyEffective(input.expectedReturn);
    var rmNPS = monthlyEffective(isUSD(input.region) ? input.expectedReturn : RATE_NPS);
    var rmPPF = monthlyEffective(RATE_PPF);
    var rmEPF = monthlyEffective(RATE_EPF);
    var rmGold = monthlyEffective(RATE_GOLD);
    var im = monthlyEffective(input.inflation);
    var expenses = nz(input.annualExpenses);
    var houseNow = house;
    var maxMonths = 80 * 12;
    for (var m = 1; m <= maxMonths; m++) {
      pots = stepPots(pots, input, m, rmLiq, rmEq, rmNPS, rmPPF, rmEPF, rmGold);
      expenses = expenses * (1 + im);
      houseNow = houseNow * (1 + im);
      if (potsSpendable(pots) >= fireNumber(expenses, swr) + houseNow) {
        out.reachesFire = true;
        out.years = m / 12;
        out.fiAge = input.age + Math.round(out.years);
        out.jewelleryLater = pots.jewellery;
        out.chart = fireChart(input, swr, out.years, true);
        return out;
      }
    }
    out.jewelleryLater = pots.jewellery;
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
    formatUSD: formatUSD,
    formatMoney: formatMoney,
    fire: fire,
    sip: sip,
    emi: emi,
    emergency: emergency,
    budget: budget,
  };
})(window);
