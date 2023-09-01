import * as radix from "@radix-ui/colors";
import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.astro"],
	// darkMode: "class",
	theme: {
		extend: {
			colors: {
				"hi-contrast": "var(--text-hi-contrast)",
				"lo-contrast": "var(--text-lo-contrast)",
				"focus-ring": "hsl(var(--jade9))",
			},
			fontFamily: {
				sans: [
					"Atkinson Hyperlegible",
					...defaultTheme.fontFamily.sans,
				],
			},
		},
	},
	plugins: [
		require("windy-radix-palette")({
			colors: {
				jade: radix.jade,
				jadeDark: radix.jadeDark,
				grass: radix.grass,
				grassDark: radix.grassDark,
				sage: radix.sage,
				sageDark: radix.sageDark,
				sand: radix.sand,
				sandDark: radix.sandDark,
			},
		}),
		require("@tailwindcss/typography"),
		require("windy-radix-typography")({
			colors: ["jade", "sage"],
		}),
	],
};
