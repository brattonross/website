export const THEMES = ["light", "dark", "auto"] as const;

export type Theme = (typeof THEMES)[number];

export function isTheme(theme: unknown): theme is Theme {
	return THEMES.includes(theme as Theme);
}

export const THEME_KEY = "brattonross_theme" as const;

export function humanizeTheme(theme: Theme): string {
	switch (theme) {
		case "light":
			return "Light";
		case "dark":
			return "Dark";
		case "auto":
			return "Auto";
	}
}

export function getThemeOrDefault(theme: unknown): Theme {
	if (isTheme(theme)) {
		return theme;
	}
	return "auto";
}
