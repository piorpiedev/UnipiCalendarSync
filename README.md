# UnipiCalendarSync

Project is based on [University Planner](https://unipi.prod.up.cineca.it/)

This project allowes you to sync any Cineca calendar with your local calendar via [CalDav](https://en.wikipedia.org/wiki/CalDAV)

It does **NOT** include the CalDav server

This project has been made strictly for strict personal usage, and published as open source material under Apache v2. 
I highly encurage you to self host your own instance for your specific year and course

---

## Example urls
The ones i use

### IT-L31 | 1A
🌍 [WebView](https://nextcloud.piorpie.com/apps/calendar/p/2LBCcnJF4wE92zqt) 
|
🔗 CalDav: `webcal://nextcloud.piorpie.com/remote.php/dav/public-calendars/2LBCcnJF4wE92zqt?export`

### IT-L31 | 1B
🌍 [WebView](https://nextcloud.piorpie.com/apps/calendar/p/3tngiqDGMxdr5yji) 
|
🔗 CalDav: `webcal://nextcloud.piorpie.com/remote.php/dav/public-calendars/3tngiqDGMxdr5yji?export`

### IT-L31 | 1C
🌍 [WebView](https://nextcloud.piorpie.com/apps/calendar/p/Zq9q9B4c6dC2QY6X) 
|
🔗 CalDav: `webcal://nextcloud.piorpie.com/remote.php/dav/public-calendars/Zq9q9B4c6dC2QY6X?export`

---

## Config

The program requires a valid config in the same path

eg:
```json
[
    {
        "cinecaUrl": "https://unipi.prod.up.cineca.it/api/Impegni/getImpegniCalendarioPubblico",
        "cinecaCalendarId": "a2937c8086cefe9555d2c2fa",
        "caldavUrl": "https://nextcloud.yourwebsite.com/remote.php/dav",
        "caldavPath": "calendars/username/calendarname",
        "year": 2,
        "partition": "A",
        "ignored_classes": [
            "the lowercase name of an ignored class"
        ]
    },
    ...
]
```