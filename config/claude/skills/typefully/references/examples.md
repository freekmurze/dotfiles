
## Examples

### Set up default social set
```bash
# Check current config
./scripts/typefully.js config:show

# Set default (interactive - lists available social sets)
./scripts/typefully.js config:set-default

# Set default (non-interactive)
./scripts/typefully.js config:set-default 123 --location global
```

### Create a tweet (using default social set)
```bash
./scripts/typefully.js drafts:create --text "Hello, world!"
```

### Create a tweet with explicit social_set_id
```bash
./scripts/typefully.js drafts:create 123 --text "Hello, world!"
```

### Create a cross-platform post (specific platforms)
```bash
./scripts/typefully.js drafts:create --platform x,linkedin,threads --text "Big announcement!"
```

### Resolve LinkedIn mention syntax from a company URL
```bash
./scripts/typefully.js linkedin:organizations:resolve --organization-url "https://www.linkedin.com/company/typefullycom/"
```

### Create a LinkedIn draft with a mention
```bash
./scripts/typefully.js drafts:create --platform linkedin --text "Thanks @[Typefully](urn:li:organization:86779668) for the support."
```

### Create a post on all connected platforms
```bash
./scripts/typefully.js drafts:create --all --text "Posting everywhere!"
```

### Create and schedule for next slot
```bash
./scripts/typefully.js drafts:create --text "Scheduled post" --schedule next-free-slot
```

### Create with tags
```bash
./scripts/typefully.js drafts:create --text "Marketing post" --tags marketing,product
```

### List scheduled posts sorted by date
```bash
./scripts/typefully.js drafts:list --status scheduled --sort scheduled_date
```

### Get queue view for a date range
```bash
./scripts/typefully.js queue:get --start-date 2026-02-01 --end-date 2026-02-29
```

### Get X post analytics for a date range
```bash
./scripts/typefully.js analytics:posts:list --start-date 2026-03-01 --end-date 2026-03-07
```

### Get X post analytics including replies
```bash
./scripts/typefully.js analytics:posts:list --start-date 2026-03-01 --end-date 2026-03-07 --include-replies
```

### Paginate through X analytics results
```bash
./scripts/typefully.js analytics:posts:list --start-date 2026-03-01 --end-date 2026-03-31 --limit 100 --offset 100
```

### Get queue schedule
```bash
./scripts/typefully.js queue:schedule:get
```

### Replace queue schedule rules
```bash
./scripts/typefully.js queue:schedule:put --rules '[{"h":9,"m":30,"days":["mon","wed","fri"]}]'
```

### Reply to a tweet
```bash
./scripts/typefully.js drafts:create --platform x --text "Great thread!" --reply-to "https://x.com/user/status/123456"
```

### Post to an X community
```bash
./scripts/typefully.js drafts:create --platform x --text "Community update" --community 1493446837214187523
```

### Create an X quote post
```bash
./scripts/typefully.js drafts:create --platform x --text "My take on this" --quote-post-url "https://x.com/user/status/1234567890123456789"
```

### Update a draft to quote an X post
```bash
./scripts/typefully.js drafts:update 456 --platform x --quote-post-url "https://x.com/user/status/1234567890123456789" --use-default
```

### Create draft with share URL
```bash
./scripts/typefully.js drafts:create --text "Check this out" --share
```

### Create draft with scratchpad notes
```bash
./scripts/typefully.js drafts:create --text "Launching next week!" --scratchpad "Draft for product launch. Coordinate with marketing team before publishing."
```

### Upload media and create post with it
```bash
# Single command handles upload + polling - returns when ready!
./scripts/typefully.js media:upload ./image.jpg
# Returns: {"media_id": "abc-123-def", "status": "ready", "message": "Media uploaded and ready to use"}

# Create post with the media attached
./scripts/typefully.js drafts:create --text "Check out this image!" --media abc-123-def
```

### Upload multiple media files
```bash
# Upload each file (each waits for processing)
./scripts/typefully.js media:upload ./photo1.jpg  # Returns media_id: id1
./scripts/typefully.js media:upload ./photo2.jpg  # Returns media_id: id2

# Create post with multiple media on one post (comma-separated)
./scripts/typefully.js drafts:create --text "Photo dump!" --media id1,id2
```

### Per-post media in a thread

Use `|` to delimit per-post media groups, mirroring how `---` delimits text into a thread. Each `|`-separated group is attached to the post at the matching index. Within a group, `,` still attaches multiple media to the same post.

```bash
# Thread of 3 posts, one image per post:
./scripts/typefully.js drafts:create --platform x --media id1,id2,id3 --text "Intro post

---

Middle post

---

Final post"
# Wait — that comma form attaches all three to post 0. Use `|` instead:

./scripts/typefully.js drafts:create --platform x --media "id1|id2|id3" --text "Intro post

---

Middle post

---

Final post"

# Mix multiple media on one post and skip another post:
./scripts/typefully.js drafts:create --media "a,b||c" --text "Two images here

---

This post has no media

---

One image on the last post"
```

A single comma-only spec (no `|`) keeps the historical behavior of attaching to post 0 only.

### Add media to an existing draft
```bash
# Upload media
./scripts/typefully.js media:upload ./new-image.jpg  # Returns media_id: xyz

# Update draft with media (456 is the draft_id)
./scripts/typefully.js drafts:update 456 --text "Updated post with image" --media xyz --use-default
```

### Setup (interactive)
```bash
./scripts/typefully.js setup
```

### Setup (non-interactive, for scripts/CI)
```bash
# Auto-selects default social set if only one exists
./scripts/typefully.js setup --key typ_xxx --location global

# With explicit default social set
./scripts/typefully.js setup --key typ_xxx --location global --default-social-set 123

# Skip default social set selection entirely
./scripts/typefully.js setup --key typ_xxx --no-default
```
