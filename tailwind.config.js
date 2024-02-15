import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.astro"],
	darkMode: "class",
	theme: {
		container: {
			center: true,
			padding: {
				DEFAULT: "1rem",
				md: "1.5rem",
				lg: "2rem",
			},
		},
		extend: {
			colors: {
				gray: {
					1: "var(--gray-1)",
					2: "var(--gray-2)",
					3: "var(--gray-3)",
					4: "var(--gray-4)",
					5: "var(--gray-5)",
					6: "var(--gray-6)",
					7: "var(--gray-7)",
					8: "var(--gray-8)",
					9: "var(--gray-9)",
					10: "var(--gray-10)",
					11: "var(--gray-11)",
					12: "var(--gray-12)",
					a1: "var(--gray-a1)",
					a2: "var(--gray-a2)",
					a3: "var(--gray-a3)",
					a4: "var(--gray-a4)",
					a5: "var(--gray-a5)",
					a6: "var(--gray-a6)",
					a7: "var(--gray-a7)",
					a8: "var(--gray-a8)",
					a9: "var(--gray-a9)",
					a10: "var(--gray-a10)",
					a11: "var(--gray-a11)",
					a12: "var(--gray-a12)",
				},
				"hi-contrast": "var(--gray-12)",
				"lo-contrast": "var(--gray-11)",
			},
			fontFamily: {
				sans: [
					"Atkinson Hyperlegible",
					...defaultTheme.fontFamily.sans,
				],
			},
			fontSize: {
				// https://github.com/radix-ui/themes/blob/main/packages/radix-ui-themes/src/styles/tokens/typography.css
				xs: [
					"12px",
					{
						lineHeight: "16px",
						fontWeight: "400",
						letterSpacing: "0.0025em",
					},
				],
				sm: [
					"14px",
					{
						lineHeight: "20px",
						fontWeight: "400",
						letterSpacing: "0em",
					},
				],
				base: [
					"16px",
					{
						lineHeight: "24px",
						fontWeight: "400",
						letterSpacing: "0em",
					},
				],
				lg: [
					"18px",
					{
						lineHeight: "26px",
						fontWeight: "400",
						letterSpacing: "-0.0025em",
					},
				],
				xl: [
					"20px",
					{
						lineHeight: "28px",
						fontWeight: "400",
						letterSpacing: "-0.005em",
					},
				],
				"2xl": [
					"24px",
					{
						lineHeight: "30px",
						fontWeight: "400",
						letterSpacing: "-0.00625em",
					},
				],
				"3xl": [
					"28px",
					{
						lineHeight: "36px",
						fontWeight: "400",
						letterSpacing: "-0.0075em",
					},
				],
				"4xl": [
					"35px",
					{
						lineHeight: "40px",
						fontWeight: "400",
						letterSpacing: "-0.01em",
					},
				],
				"5xl": [
					"60px",
					{
						lineHeight: "60px",
						fontWeight: "400",
						letterSpacing: "-0.025em",
					},
				],
			},
		},
	},
	plugins: [require("@tailwindcss/typography")],
};
