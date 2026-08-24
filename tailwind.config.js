/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.{templ,go}"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        paper: "#eef3fb",
        ink: "#0f1b33",
        night: "#071018",
        sky: "#38bdf8",
        brand: {
          DEFAULT: "#2b6cef",
          dark: "#1d4ed8",
          muted: "#2b6cef1a",
        },
      },
      fontFamily: {
        serif: ["Fraunces", "ui-serif", "Georgia", "serif"],
        sans: ["DM Sans", "ui-sans-serif", "system-ui", "sans-serif"],
      },
      boxShadow: {
        card: "0 1px 2px rgba(15,27,51,0.06), 0 14px 36px rgba(43,108,239,0.10)",
        glow: "0 8px 30px rgba(43,108,239,0.32)",
      },
    },
  },
  plugins: [],
};
