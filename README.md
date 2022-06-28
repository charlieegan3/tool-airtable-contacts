# airtable-contacts

This is a binary with functionality to sync a custom contact database from
Airtable.

**Note:** It's not likely this will have any useful functionality for anyone
other than the author.

Current functionality:

* Download contact data
* Generate vCard file for one or more contacts
* Sync the contacts to carddav server
* Send notifications about birthdays and special events

Example config file:

```
airtable:
  key: xxxxxxxxxxxxxxxxx
  base: xxxxxxxxxxxxxxxxx
  table: xxxxxxxxxxxxxxxxx
  view: Active
vcard:
  use_v3: true
  photo:
    size: 100
dropbox:
  path: /contacts.vcard
  token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
carddav:
  serverURL: https://carddav.fastmail.com/dav/addressbooks/user/xxxxxxxxxxxxxxxxxxxxxxxxx/Default
  user: xxxxxxxxxxxxxxxxxxxxxxxxx
  password: xxxxxxxxxxxxxxxx
pushover:
  user_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  app_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
webhook:
  endpoint: https://hooks.zapier.com/hooks/catch/xxxxxxxxxxxxxx
```
