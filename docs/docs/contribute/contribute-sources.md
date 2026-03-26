---
id: contribute-sources
sidebar_position: 2
title: Contribute Sources
description: Draft playbook for adding or improving RSS feeds in Colibri.
---

Colibri thrives on high-quality feeds. Use this quick checklist to add or tweak sources.

## 1. Edit the CSV

1. Open `cmd/fetcher/sources/sources.csv`.
2. Each row follows `id,name,url,category`.
   - `id`: slugged, lowercase, no spaces (`hacker-news`).
   - `name`: readable title.
   - `url`: fully-qualified RSS/Atom URL (`https://...`).
   - `category`: Title Case bucket (`Technology`, `Humour`, etc.). Reuse existing categories when possible.
3. Append your new row. Keep the header unchanged and avoid trailing commas.

## 2. Run the validation tests

1. From the repo root execute:
   ```bash
   go test ./internal/sources
   ```
2. The suite verifies CSV parsing plus the `ValidateSource` rules (slug format, URL, category casing). Fix any failures before committing.

## 3. Open a pull request

- Describe the feed(s) you added and why they are useful.
- Mention the test command above so reviewers know it passed.

That’s it, adding sources is intentionally lightweight: update the CSV, run the validation tests, and send a PR.

> We eventually plan to implement a more straightforward way to add sources via a web form. For the time being, adding a new source requires submitting a Pull Request.We plan on adding a more straight forward way to add sources from a web from. 