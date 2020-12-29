## Unique articles API
Simple service to store articles and check for duplicates.

##### To start service on localhost:8080 run:
`docker-compose up`

### API:
To add article `POST /articles` with json `{"content": "..."}`, which returns `{"id": 5, "content": "...", "duplicate_article_ids": [1, 3]}`

To get article by id `GET /articles/:id`, which returns article if exists
`{"id": 3, "content": "...", "duplicate_article_ids": [1, 5]}`

`GET /articles` returns all saved articles
```
{"articles": [
  {"id": 1, "content": "...", "duplicate_article_ids": [3, 5]},
  {"id": 2, "content": "...", "duplicate_article_ids": []},
  {"id": 3, "content": "...", "duplicate_article_ids": [1, 5]},
  {"id": 4, "content": "...", "duplicate_article_ids": []},
  {"id": 5, "content": "...", "duplicate_article_ids": [1, 3]}
]}
```

#### Duplicates methodology:
##### Articles are filtered before comparison:
- removed punctuation, frequent words, transitional expressions
- words stemming (process of reducing inflected (or sometimes derived) words to their word stem, base or root form)
- convert synonyms to base form

##### Articles comparison:
Modified Levenshtein distance for measuring the difference between two articles.
Distance between two articles is the minimum number of single-word edits (insertions, deletions or substitutions) required to change one article into the other

Default duplicate level: **95%**, could be changed in `configs/config.json`

#### Database
In db stored two versions of article - content and filtered content.
Used in-memory database for simplicity.
