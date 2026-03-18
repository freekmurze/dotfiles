---
name: schema-markup
version: 1.0.0
description: >
  Generates, validates, and fixes JSON-LD structured data for schema.org markup.
  Implements Article, Product, FAQPage, BreadcrumbList, Organization, and other
  schema types. Debugs rich snippet errors and validates against Google Rich
  Results Test. Use when adding schema markup, structured data, JSON-LD, rich
  snippets, FAQ schema, product schema, review schema, or breadcrumb schema.
  For broader SEO issues, see seo-audit.
---

# Schema Markup

## Workflow

1. **Audit** — check for existing schema and errors (search page source for `application/ld+json`)
2. **Generate** — create JSON-LD markup matching the page content
3. **Validate** — test with [Rich Results Test](https://search.google.com/test/rich-results) and [Schema.org Validator](https://validator.schema.org/)
4. **Fix** — resolve any errors or warnings until clean
5. **Deploy** — place in `<head>` or end of `<body>`, only when validation passes

If `.claude/product-marketing-context.md` exists, read it first for business context.

---

## Common Schema Types

| Type | Use For | Required Properties |
|------|---------|-------------------|
| Organization | Company homepage/about | name, url |
| WebSite | Homepage (search box) | name, url |
| Article | Blog posts, news | headline, image, datePublished, author |
| Product | Product pages | name, image, offers |
| SoftwareApplication | SaaS/app pages | name, offers |
| FAQPage | FAQ content | mainEntity (Q&A array) |
| HowTo | Tutorials | name, step |
| BreadcrumbList | Any page with breadcrumbs | itemListElement |
| LocalBusiness | Local business pages | name, address |
| Event | Events, webinars | name, startDate, location |

**For all JSON-LD examples**: See [references/schema-examples.md](references/schema-examples.md)

---

## Key Examples

### Article (blog posts, news)

```json
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "How to Implement Schema Markup",
  "image": "https://example.com/image.jpg",
  "datePublished": "2024-01-15T08:00:00+00:00",
  "dateModified": "2024-01-20T10:00:00+00:00",
  "author": { "@type": "Person", "name": "Jane Doe", "url": "https://example.com/authors/jane" },
  "publisher": { "@type": "Organization", "name": "Example Company" }
}
```

### FAQPage

```json
{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [
    { "@type": "Question", "name": "What is schema markup?",
      "acceptedAnswer": { "@type": "Answer", "text": "Schema markup is structured data..." } },
    { "@type": "Question", "name": "How do I implement it?",
      "acceptedAnswer": { "@type": "Answer", "text": "Use JSON-LD format in your page head..." } }
  ]
}
```

### Combining types with `@graph`

```json
{
  "@context": "https://schema.org",
  "@graph": [
    { "@type": "Organization", "@id": "https://example.com/#org", "name": "Example", "url": "https://example.com" },
    { "@type": "WebSite", "url": "https://example.com", "publisher": { "@id": "https://example.com/#org" } },
    { "@type": "BreadcrumbList", "itemListElement": [...] }
  ]
}
```

---

## Common Errors

| Error | Fix |
|-------|-----|
| Missing required properties | Check Google's docs for required fields per type |
| Invalid date values | Use ISO 8601 format: `2024-01-15T08:00:00+00:00` |
| URLs not fully qualified | Use `https://example.com/page`, not `/page` |
| Schema doesn't match content | Markup must reflect visible page content |

---

## Validation Checklist

- [ ] Validates in [Rich Results Test](https://search.google.com/test/rich-results)
- [ ] No errors or warnings
- [ ] Matches visible page content
- [ ] All required properties included
- [ ] Monitor Search Console enhancements reports after deploy

---

## Related Skills

- **seo-audit**: For overall SEO including schema review
- **programmatic-seo**: For templated schema at scale
