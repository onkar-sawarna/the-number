(function () {
  var m = document.cookie.match(/(?:^|; )theme=([^;]*)/);
  var t = m ? decodeURIComponent(m[1]) : "";
  if (t !== "dark" && t !== "light") {
    t = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  document.documentElement.classList.toggle("dark", t === "dark");
})();
