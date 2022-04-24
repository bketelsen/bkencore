module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      typography: {
        DEFAULT: {
          css: {
            code: {
              color: "#A21AAF",
            },
            a: {
              color: "#6D26DA",
              '&:hover': {
                color: "#FFFFFF",
                backgroundColor: "#6D26DA",
                borderRadius: "0.2rem",
                textDecoration: "none",
              },
            },
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
}
