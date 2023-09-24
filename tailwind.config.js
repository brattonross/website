import defaultTheme from "tailwindcss/defaultTheme";
import { alias } from "windy-radix-palette";

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.astro"],
	darkMode: "class",
	theme: {
		extend: {
			colors: {
				"hi-contrast": "var(--text-hi-contrast)",
				"lo-contrast": alias("sand", 11),
				"focus-ring": alias("jade", 9),
			},
			fontFamily: {
				sans: ["Atkinson Hyperlegible", ...defaultTheme.fontFamily.sans],
			},
			// Based on https://github.com/radix-ui/themes/blob/main/packages/radix-ui-themes/src/styles/tokens/typography.css
			fontSize: {
				"5xl": [
					"60px",
					{
						letterSpacing: "-0.025em",
						lineHeight: "60px",
					},
				],
				"4xl": [
					"35px",
					{
						letterSpacing: "-0.01em",
						lineHeight: "40px",
					},
				],
				"3xl": [
					"28px",
					{
						letterSpacing: "-0.0075em",
						lineHeight: "36px",
					},
				],
				"2xl": [
					"24px",
					{
						letterSpacing: "-0.00625em",
						lineHeight: "30px",
					},
				],
				xl: [
					"20px",
					{
						letterSpacing: "-0.005em",
						lineHeight: "28px",
					},
				],
				lg: [
					"18px",
					{
						letterSpacing: "-0.0025em",
						lineHeight: "26px",
					},
				],
				base: [
					"16px",
					{
						letterSpacing: "0em",
						lineHeight: "24px",
					},
				],
				sm: [
					"14px",
					{
						letterSpacing: "0em",
						lineHeight: "20px",
					},
				],
				xs: [
					"12px",
					{
						letterSpacing: "0.0025em",
						lineHeight: "16px",
					},
				],
			},
		},
	},
	plugins: [require("windy-radix-palette"), require("@tailwindcss/typography")],
};
