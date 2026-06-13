import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        ink: "#17181f",
        sand: "#f4efe7",
        coral: "#f26a4b",
        gold: "#e9b949",
        sea: "#4aa3a2",
      },
      boxShadow: {
        card: "0 24px 64px rgba(23, 24, 31, 0.14)",
      },
    },
  },
  plugins: [],
};

export default config;

