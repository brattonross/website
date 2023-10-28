window.customElements.define(
	"theme-select",
	class extends HTMLFieldSetElement {
		public connectedCallback() {
			this.addEventListener("change", this.#handleChange);
		}

		public disconnectedCallback() {
			this.removeEventListener("change", this.#handleChange);
		}

		#handleChange(event: Event) {
			if (event.target === null || !("value" in event.target)) {
				return;
			}

			const isDarkMode =
				event.target.value === "dark" ||
				(event.target.value === "auto" && this.#prefersDarkMode);
			document.documentElement.classList.toggle("dark", isDarkMode);
			document.documentElement.style.colorScheme = isDarkMode
				? "dark light"
				: "light dark";
		}

		get #prefersDarkMode() {
			return window.matchMedia("(prefers-color-scheme: dark)").matches;
		}
	},
	{ extends: "fieldset" },
);
