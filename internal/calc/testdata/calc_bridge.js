"use strict";

const fs = require("fs");
const vm = require("vm");

const calcPath = process.argv[2];
if (!calcPath) {
  process.stderr.write("usage: node calc_bridge.js path/to/calc.js < jobs.json\n");
  process.exit(2);
}

const src = fs.readFileSync(calcPath, "utf8");
const sandbox = {};
sandbox.globalThis = sandbox;
vm.createContext(sandbox);
vm.runInContext(src, sandbox);
const Calc = sandbox.Calc;
if (!Calc) {
  process.stderr.write("Calc missing after loading " + calcPath + "\n");
  process.exit(1);
}

const jobs = JSON.parse(fs.readFileSync(0, "utf8"));
const out = [];
for (const job of jobs) {
  switch (job.op) {
    case "fire":
      out.push(Calc.fire(job.input));
      break;
    case "sip":
      out.push(Calc.sip(job.input));
      break;
    case "emi":
      out.push(Calc.emi(job.input));
      break;
    case "emergency":
      out.push(Calc.emergency(job.input));
      break;
    case "budget":
      out.push(Calc.budget(job.input));
      break;
    case "formatINR":
      out.push(Calc.formatINR(job.n));
      break;
    case "formatUSD":
      out.push(Calc.formatUSD(job.n));
      break;
    default:
      process.stderr.write("unknown op " + job.op + "\n");
      process.exit(1);
  }
}
process.stdout.write(JSON.stringify(out));
