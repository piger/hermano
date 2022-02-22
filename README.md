# hermano

Scraper for used Herman Miller chairs sold by [Griffin Office](https://usedaeronireland.ie/used-herman-miller-aeron-chairs/).

## Configuration file

The configuration file is used to specify the credentials to send Push notifications via Pushover and to list
all product names that will be ignored by the program; the structure of the configuration file is:

```
api_token = "AABBBCCDDEEFF"
user_key = "FFFEEEDDCCBBAAA"
ignored = [
	"title of a product that I want to ignore",
    "title of another product that I don't want to show"
]
```

`api_token` and `user_key` can be omitted to prevent sending Push notifications.

