# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful. 

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission. 


## Submission

1. Fork this repository and send us a link to your fork after pushing your changes. 
2. Heroku hosting - The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

## Demo

The code in this repo is live on Heroku at https://ppai-shakesearch.herokuapp.com/.

## Changelog + Prioritized wishlist

- [X] Fix out of bound panic in handling of search index.
- [X] Highlight search query in rendered result.
- [X] Retain whitespace in search results, as in the original corpus.
- [X] Avoid overlapping index in search results to avoid repitition of matches.
- [X] Feature: Case insensitive search.
- [X] Feature: Search within a specific work.
- [X] Feature: Added pagination (client-side, only for usablity, not for perf).
- [X] Show results starting at a sentence (decent effort, not perfect)
- [X] Avoid abrupt word cuts in search result ending. Add ellipsis just like Google results.
- [X] Remove "," separator in search results and add a clear separator in search results.
- [ ] Unit tests.
- [ ] Display gracefull error when no match is found.
- [ ] Support finding a phrase that is split across multiple lines.
- [ ] Add context to search results such as title of work, act/section, line number.
- [ ] Add link to rendered work with page anchored at the spot of occurence. Like github's #L feature.
- [ ] Autocomplete for phrases.
- [ ] Fuzzy search ("did you mean" feature)
- [ ] Spellcheck
- [ ] Performance: server side pagination, caching etc

Code has been refactored to be modular. Search has been moved to a reusable package.
