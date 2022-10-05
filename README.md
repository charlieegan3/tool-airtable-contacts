# tool-airtable-contacts

This is a tool for [toolbelt](https://github.com/charlieegan3/toolbelt) which wraps functionality to run a contacts
database in [Airtable](http://airtable.com). It syncs to a CardDAV endpoint and uses
[webhook-rss](https://github.com/charlieegan3/tool-webhook-rss) to send updates via RSS.

Example config:

```yaml
tools:
  ...
  airtable-contacts:
    jobs:
      day:
        schedule: "0 0 5 * * *"
      week:
        schedule: "0 0 5 * * 0"
      sync:
        schedule: "0 */15 * * * *"
    endpoint: https://...
    airtable:
      key: xxxxxxxxxxxxxxxxx
      base: xxxxxxxxxxxxxxxxx
      table: xxxxxxxxxxxxxxxxx
      view: xxx
    carddav:
      server_url: https://...
      user: name@example.com
      password: xxxxxxxxxxx
    vcard:
      use_v3: true
      photo:
        size: 100
```
