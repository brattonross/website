---
title: Hypermedia is a natural answer
date: 2023-10-01
description: A natural solution to front-end communication with APIs.
---

If you've done front-end development where you are building a client-side app that consumes an API that is owned by some other set of people, you've likely come across this issue; as your front-end design and requirements change, you either end up repeatedly asking the owners of the API to make changes to better accommodate the front-end, or you end up with a bunch of JavaScript code that transforms that API responses into a format that your UI can more easily consume.

These solutions are not ideal. Either you are pestering people for changes, which can be slow, not to mention there are probably more valuable things those people could be doing, or you ship a bunch of JS to the browser that shouldn't really be necessary. It feels like a lose-lose to me.

Other solutions do exist of course. GraphQL is a way for the API devs to give the front-end devs the flexibility to request whatever they want from the database, so that they don't have to make changes to accommodate them all the time. But still, it is JS that has to be included on the client, not to mention it can [open up new attack vectors](https://intercoolerjs.org/2016/02/17/api-churn-vs-security.html#the-problem-with-the-solution).

## Front-end devs ought to be API devs

A solution that appeals to me is to have front-end developers own a Backend-for-Frontend (BFF) API. This is a pattern that puts an API in front of a particular implementation of an application, and allows it to send only exactly what is required down to the client app, leading to a much leaner front-end code base. Plus, we get to avoid any pesky CORS issues and the like as a bonus.

This pattern is already popular in the front-end and React world, with frameworks like Next.js and Remix allowing you to do data fetching and actions on the server-side, away from the browser. These frameworks build a BFF directly into the same codebase as the front-end. This can feel really powerful as a front-end developer, personally it feels a lot more freeing to know that I can write whatever I want to run server-side. I don't have to bother other teams to ask them to make changes just for my front-end, I can use server-only APIs, etc.

But, some things still irk me about this approach; we are communicating between the client and server with JSON, most likely duplicating logic between the two parts, and we are shipping JS for the framework and meta-framework to the browser. We can do better.

## JSON -> HTML

Due to the popularity of the SPA and micro-services, if you ask a developer today to describe what an API is, they will likely tell you that it is a server that serves JSON responses and uses JWT for authentication. Whilst this sort of approach may be appropriate for micro-service APIs that communicate amongst one another, for a web page, it is definitely not ideal.

The kind of API that the browser is best at consuming isn't JSON--it's HTML. Browsers are very good at dealing with HTML, so why don't we make our API serve HTML to the browser? Since HTML is a form of Hypermedia, our BFF becomes a Hypermedia API! Since we're serving HTML, we now don't need to ship JS for a front-end framework to the browser.

_But what about muh interactivity?_ I hear you cry.

Fret not, chances are you probably don't even need it. A lot of us are building simple applications that do things like CRUD on some objects. There is no need for reactivity here; just plain HTML forms and a little sprinkle of JS can get you a long way.

_Ok but what if I really do need reactivity?_

I get you, even in CRUD applications like an admin dashboard, we probably do want some small bits of interactivity like dropdown menus, modals, conditionally displayed form fields, etc. We don't need something as complex as React for this though. For a lot of applications, we can probably do this stuff in "vanilla" JS. Browser APIs are pretty good nowadays. If you want to reach for a library, something like [Alpine.js](https://alpinejs.dev) is suitable enough.

In the rare case that you are actually building something that really needs lots of interactivity on the front-end, go ahead, reach for that framework. I'm not going to tell anybody.

### JavaScript on the back-end

If you're a developer who only knows JS, and isn't yet comfortable picking up another language to write a server in, I'd recommend [Astro](https://astro.build). It is perfectly capable of serving Hypermedia APIs in its SSR mode, in fact, this site runs with that exact setup. The best choice for HTTP/Hypermedia servers in my opinion is Go, and it is quite a simple language to pick up and get running with (hint: that's part of what it was designed for). Give it a go if you haven't already.

## Extending HTML

Once you've worked enough with plain HTML APIs, you start to realise that there are some limitations. Part of the reason why writing a bunch of JS on the front-end got popular was because of these limitations. Imagine, for example, a form that takes a little while to submit; a few seconds or more. If you are using plain HTML and CSS then the only indication to your user that something is happening is the little spinner in the browser tab, if they even notice that. They might click the button again, report an issue, or even just give up and leave because they didn't get good feedback about what was happening. This can be remedied with some JS: disable the button when the form is submitted, show a spinner inline, or maybe even a progress bar depending on the situation.

This is an example of how some JavaScript on the page can be a benefit to your users. The ideal situation is keeping the amount of JS we ship slim, whilst using [Progressive Enhancement](https://en.wikipedia.org/wiki/Progressive_enhancement) to make sure that our users get at least basic functionality, with a more rich experience if their internet connection can handle it.

One library that makes this sort of stuff easier, that I have enjoyed so far, is [htmx](https://htmx.org). It extends HTML attributes to allow any element to become a "hypermedia control". Now, instead of being limited to forms and anchor tags as ways to communicate with the server, we can make any elements send requests with any HTTP verb we want. We can still fall back to forms and anchors for users with a slow connection, but if users are able to load (a small amount of) JS, then we can deliver them a SPA-like experience.

Htmx also has the bonus of working with any server-side that you want, which the creator refers to as [HOWL](https://htmx.org/essays/hypermedia-on-whatever-youd-like/) (Hypermedia On Whatever you'd Like). This is a big plus in my eyes, as I can bring my htmx knowledge with me between projects and jobs, even if they use different languages, and still be just as productive. Personally, I also feel like HTML will eventually become somewhat htmx-like, with more elements than just forms and anchors being hypermedia controls, so being able to think in terms of hypermedia seems like a decent investment.

## A natural answer

When I think about building web applications, I naturally come to the conclusion that I should be using hypermedia. I've described my thought processes above, but to sum it up:

- I believe that the devs building the front-end should own a BFF for their front-end
- If you own the API for your application;
  - Why would you want to duplicate logic between the client and server? Keep the logic on the server and use hypermedia as the engine of application state (HATEOAS)
  - Why would we want to communicate between the client and server with JSON, when the browser is better at consuming HTML?
- In the majority of cases, we probably don't need complex client-side JS, most likely we can get by with forms and progressive enhancement
- There are limitations to plain HTML, so enhancing our abilities using libraries like htmx and Alpine can give users a much better experience
