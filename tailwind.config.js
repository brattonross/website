import defaultTheme from "tailwindcss/defaultTheme";
import radix, { alias } from "windy-radix-palette";

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.astro"],
	darkMode: "class",
	theme: {
		extend: {
			boxShadow: {
				b: "0 1px var(--tw-shadow-color)",
				t: "0 -1px var(--tw-shadow-color)",
			},
			colors: {
				"hi-contrast": "var(--text-hi-contrast)",
				"lo-contrast": alias("sand", 11),
				gray: alias("sand"),
				grayA: alias("sandA"),
			},
			fontFamily: {
				sans: ["Atkinson Hyperlegible", ...defaultTheme.fontFamily.sans],
			},
		},
	},
	plugins: [radix, require("@tailwindcss/typography")],
};
