import { defineMiddleware } from "astro:middleware";

export const onRequest = defineMiddleware(async (context, next) => {
	const isHtmxRequest = context.request.headers.get("Hx-Request") === "true";
	if (!isHtmxRequest) {
		return await next();
	}

	const response = await next();
	let html = await response.text();

	html = html.replace(/<style type\=\"text\/css\"(.*)<\/style>/s, "");
	html = html.replace(/<script(.*)<\/script>/gm, "");

	return new Response(html, response);
});
