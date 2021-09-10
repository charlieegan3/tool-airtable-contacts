# airtable-contacts

This is a binary with functionality to sync a custom contact database from
Airtable.

**Note:** It's not likely this will have any useful functionality for anyone
other than the author.

Current functionality:

* Download contact data
* Generate vCard file for one or more contacts
* Sync the contacts to carddav server

Example config file:

```
airtable:
  key: xxxxxxxxxxxxxxxxx
  base: xxxxxxxxxxxxxxxxx
  table: xxxxxxxxxxxxxxxxx
  view: xxxxxx
vcard:
  photo:
    size: 100
carddav:
  serverURL: https://example.com/dav/addressbooks/user/user@example.com/Default
  user: user@example.com
  password: xxxxxxxxxxxxxxxx
```
