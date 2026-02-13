/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
      "./internal/views/**/*.templ", // components
      "./internal/**/*.templ", // root
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

