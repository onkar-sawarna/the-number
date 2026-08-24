/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.{templ,go}"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        paper: "#f7f6f3",
        ink: "#14181f",
        night: "#0c0f12",
        brand: {
          DEFAULT: "#0f766e",
          muted: "#0f766e1a",
        },
      },
      fontFamily: {
        serif: ["Fraunces", "ui-serif", "Georgia", "serif"],
        sans: ["DM Sans", "ui-sans-serif", "system-ui", "sans-serif"],
      },
      boxShadow: {
        card: "0 1px 2px rgba(20,24,31,0.05), 0 12px 32px rgba(20,24,31,0.04)",
      },
    },
  },
  plugins: [],
};
