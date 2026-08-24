(function () {
  var m = document.cookie.match(/(?:^|; )theme=([^;]*)/);
  var t = m ? decodeURIComponent(m[1]) : "";
  if (t !== "dark" && t !== "light") {
    t = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  var dark = t === "dark";
  document.documentElement.classList.toggle("dark", dark);
  var meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute("content", dark ? "#071018" : "#eef3fb");

  var rm = document.cookie.match(/(?:^|; )region=([^;]*)/);
  var region = rm ? decodeURIComponent(rm[1]) : "";
  if (region !== "us" && region !== "in") region = "in";
  document.documentElement.setAttribute("data-region", region);
  document.documentElement.setAttribute("lang", "en");
  document.cookie = "lang=; path=/; max-age=0; samesite=lax";
})();
