/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
      content: ["*.templ"],
      theme: {
            maxWidth: {
                  "1/4": "25%",
                  "1/2": "300px",
                  "3/4": "75%",
            },
      },
      theme: {},
      plugins: [],
};
