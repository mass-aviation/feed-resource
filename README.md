# Concourse feed resource

This resource supports various feed types, starting from RSS and ending with JSON.  
Check [gofeed](https://github.com/mmcdole/gofeed#supported-feed-types) readme for more details about supported formats.

This resource exists because [concourse-rss-resource](https://github.com/suhlig/concourse-rss-resource) does not support other formats than RSS.

## Usage

`check`, `in` and `out` behaviour are supposed to be almost identical with concourse-rss-resource.

### Resource type

```yaml
resource_types:
- name: feed-resource
  type: registry-image
  source:
    repository: docker.zentria.ee/mass-aviation/feed-resource
```

### Resource & Example usage

```yaml
resources:
- name: cloudflare-incidents
  type: feed-resource
  icon: bomb
  public: true
  source:
    url: 'https://www.cloudflarestatus.com/history.atom'

jobs:
- name: notify-cloudflare-incident
  public: true
  plan:
  - get: cloudflare-incidents
    trigger: true

  - load_var: incident-title
    file: cloudflare-incidents/title

  - load_var: incident-date
    file: cloudflare-incidents/pubDate

  - load_var: incident-content
    file: cloudflare-incidents/content

  - put: email-notify
    params:
      subject_text: |
        [Auto] New Cloudflare incident reported - ((.:incident-title))
      body_text: |
        Report found at: ((.:incident-date))

        ((.:incident-content))
```
