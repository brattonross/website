---
title: Keeping it simple
date: 2023-07-24
description: Simplicity can be enjoyable.
---

I realised recently that I tend to like to keep things in my life simple, and people that I am inspired by tend to follow this way of thinking too. Whether it is writing code, picking clothes for the day, or cooking food; in my book simple wins out more often than not.

Simple can seem boring, but it can be exciting and enjoyable if done well. Take these two things for example—the way that Gordon Ramsay redesigns menus for restaurants in his Kitchen Nightmares series, and the Go programming language. Two things that don't seem related, but both I feel are decent examples of simplicity being effective.

In Kitchen Nightmares, Ramsay will often find that restaurants have a menu that is too bloated and/or complex, forcing the chefs to have to resort to preparing food far in advance and freezing it, or buying in pre-made dishes. Inevitably, Gordon will cut the size of the menu down, and give it a clear focus that suits the theme of the restaurant, with dishes that are simple enough for the chefs to prepare quickly, yet exciting and tasty enough to wow customers.

Go is a language that is known for its simplicity. There is generally one way to do things in Go—everyone who has used it knows the old `if err != nil` line, it is famously easy to pick up, and the language designers have a good idea in mind about what they want the language to be, refusing to add new features, or taking their time to implement them correctly (see Go generics).

When we write a solution to a problem as developers, we can do ourselves, and our users, a favour by writing less code overall to solve the problem. I'm not talking about writing a sweet one-liner to make yourself look cool to your colleagues—save that for your Advent of Code solutions so that you can show off to everyone on GitHub. I'm talking about choosing a good tool for the job; picking a language that you are experienced with, and that does the job to a decent standard. Not choosing to put state into a global store with reducers, actions, and selectors, when it could just be a few lines in a component. Write your entire project in a single file until it warrants being split up in some way.

I've heard people say that they think that _"You Ain't Gonna Need It"_ (YAGNI) is trending again at the moment, and I have to say that I think its not a bad trend to be a part of.
