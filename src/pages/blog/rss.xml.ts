import rss from "@astrojs/rss";
import type { APIContext } from "astro";
import { getCollection } from "astro:content";

export async function GET(context: APIContext) {
	const articles = await getCollection("blog");

	return rss({
		title: "Ross Bratton",
		description: "Ross Bratton's personal blog",
		site: context.site ?? "https://brattonross.xyz",
		items: articles.map((article) => {
			return {
				title: article.data.title,
				description: article.data.description,
				link: `/blog/${article.slug}`,
				pubDate: article.data.date,
			};
		}),
		customData: `<language>en</language>`,
	});
}
