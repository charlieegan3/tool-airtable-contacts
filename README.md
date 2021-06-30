# airtable-contacts

This is a binary with functionality to sync a contact database from Airtable.

It's not likely this will have any useful functionality for anyone other than
the author.

Current functionality:

* Download contact data
* Generate vCard file
* Upload card to location in Dropbox

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
dropbox:
  path: /contacts.vcard
  token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
