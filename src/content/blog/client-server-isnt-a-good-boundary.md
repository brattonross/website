---
title: Client/server isn't a good boundary
date: 2023-10-28
description: Developers are too quick to specialise
---

In a [recent post on X](https://twitter.com/cramforce/status/1717951682049708514), Malte Ubl writes:

> I'm all for specialization in software engineering, but client/server isn't a good boundary. As a frontend engineer you need to control the entire user experience, and that requires being empowered to program the frontend server.

When I saw this post, it immediately resonated with me. I've been planning on making a blog post to follow up on my [previous one](/blog/hypermedia-is-a-natural-answer), where I mentioned that it is my feeling that in a lot of cases, having frontend developers own an API for their frontend comes with a bunch of benefits.

Malte managed to put my feelings into more concise words than I could manage, so I wanted to share this with you, and expand on it a bit with my own thoughts.

## You can't expect to specialise

There seems to be a trend in web development where people will firmly plant themselves in the "frontend" or "backend" camp. Either they build UIs with JavaScript frameworks on the client, fetching data using libraries like React Query or GraphQL, and handling JSON API responses, or they work purely in the domain of data, building "RESTful" JSON APIs, querying databases, and reluctantly adding endpoints to their APIs to satisfy the needs of their frontend colleagues.

This seems counter-intuitive to me. I get the sense that this way of thinking has come from the popularisation of aspiring to be a FAANG engineer; in these huge tech companies, there are enough engineers that specialising in a particular area makes sense. There may be an entire team dedicated to just a single section of a web app front-end. But most of us aren't FAANG engineers, and we ought not to copy the way that they do things just because it is what the big companies are doing.

I've worked most of my career at small to medium sized companies, where there are not so many engineers. In fact, at my current job, I was the only engineer on my team for a while. Therefore, I'm used to taking on multiple roles, diving into various codebases to fix issues, front-end or back-end. Whatever I'm required to do to help my team succeed.

It may be surprising then to learn that my official job title is "Front-end Developer". I don't blame the company I work for for giving me a title that misrepresents my skills; I think we as developers have become so used to identifying with this frontend/backend split that it is just a normal occurance in the industry nowadays. It occurs to me that making the desire to specialise commonplace may alienate newer developers into thinking that they have to "pick a side", and end up limiting their abilities or opportunities.

By no means am I saying that we shouldn't specialise at all however, I think if you go very deep into a particular subject or technology, you can become an invaluable asset to your team. I just don't think that you should specialise so much that you end up hurting your career in the process. Learn a decent amount about front-end and back-end topics. Build and deploy a full-stack application on a VPS. If something in particular stands out to you as interesting, dig into it further and become an expert at it. You'll be a better developer for it.

## Owning an API as a front-end developer

In that post, Malte mentions that controlling the entire user experience requires programming the front-end server. My interpretation of this is that in order to provide the best experience for your users, you often need to specialise how your back-end works, so that it collaborates well with your front-end. This can be an issue if your front-end directly talks to a data API, like a JSON API, that is owned by another set of developers. They may be unwilling, or unable, to add features to the API that you would like in order to improve your user experience. It would be much easier if we could just control the API ourselves. Well, we can in a way.

By introducing another API inbetween our front-end and whatever other services it communicates with, we have given ourselves several new opportunities, including transforming data into a format that makes sense for us before returning it, performing any other back-end work that needs to be done, and even using our preferred method of storing session details, such as cookies, which can be more secure than dealing with tokens, which your back-end APIs might require. No longer do we need to ask other developers to make changes to their API just for us, we can instead just fetch that extra bit of data ourselves, in the format that makes most sense for us.

This pattern is commonly referred to as a Backend for Frontend (BFF).

## What is a BFF anyway?

BFFs are a pattern in software where you build an API that is meant specifically for a particular front-end implementation. For example, you might be building a social media app, and you want both a web app and mobile app. Instead of sharing an API between the two, it can be beneficial to introduce a BFF for each app. The BFF will provide specialised endpoints that surface the data that the application wants, in a format that is most appropriate for it. This is useful because the web app may include slightly different features from the mobile app, or may want to lay things out in a different way, which would work better if it had the data coming back from the API in a different format.

```
Without BFF
┌────────┐       ┌───────┐
│Data API│──────>│Web App│
└────────┘   │   └───────┘
             │   ┌──────────┐
             └──>│Mobile App│
                 └──────────┘

With BFF
┌────────┐       ┌───────┐       ┌───────┐
│Data API│──────>│Web BFF│──────>│Web App│
└────────┘   │   └───────┘       └───────┘
             │   ┌──────────┐    ┌──────────┐
             └──>│Mobile BFF│───>│Mobile App│
                 └──────────┘    └──────────┘
```

If you're a front-end developer, you might actually be using a BFF already, without even realising it.

Frameworks for building front-end applications these days tend to come with the ability to run server-side code. Astro, Next.js, Remix, Nuxt, and SvelteKit are all examples of frameworks that allow you to run server-side JavaScript, i.e. run a BFF alongside your front-end. In fact, [Remix has a page in their docs about BFFs](https://remix.run/docs/en/main/guides/bff).

## Client/server isn't a good boundary

I definitely agree with Malte. I feel that developers categorising ourselves into "back-end" and "front-end" can make sense, mainly in the context of larger companies, or in cases where the developer is choosing to specifically focus on a particular technology. For most developers, however, it doesn't seem beneficial to me.

Rather than a split of client/server, we should instead create boundaries at points that still allow developers to be flexible with their work. A front-end developer that can control the API behind can produce an app that is more efficient, and provides a better experience for users.

## Bonus: Hypermedia

It is my opinion that for the vast majority of applications, if you own the API for the front-end, then you'll be doing yourself a favour by making that API a hypermedia API that serves HTML responses, rather than JSON. [I've written about this in more detail in the past](/blog/hypermedia-is-a-natural-answer), if you're interested.
