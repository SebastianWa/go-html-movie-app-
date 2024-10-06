/** @type {import('tailwindcss').Config} */
const colors = require("tailwindcss/colors");

module.exports = {
      content: ["*.templ"],
      theme: {
            // colors: {
            //       transparent: "transparent",
            //       current: "currentColor",
            //       indigo: colors.indigo,
            // },
            maxWidth: {
                  "1/4": "25%",
                  "1/2": "300px",
                  "3/4": "75%",
            },
      },
      theme: {},
      plugins: [],
};
