import type { Config } from "tailwindcss";

const ink = "rgb(40, 40, 37)";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        jost: ["var(--font-jost)", "sans-serif"],
      },
      colors: {
        ink,
        lime:   "#e7fe05",
        yellow: "#f6cb44",
        pink:   "#fe91e6",
        green:  "#38ed81",
        teal:   "#20a291",
        purple: "#9723c9",
      },
      borderColor: {
        DEFAULT: ink,
      },
      boxShadow: {
        sm: `2px 2px 0 0 ${ink}`,
        md: `3px 3px 0 0 ${ink}`,
        lg: `4px 4px 0 0 ${ink}`,
        xl: `6px 6px 0 0 ${ink}`,
      },
      borderRadius: {
        DEFAULT: "4px",
      },
    },
  },
  plugins: [],
};

export default config;
