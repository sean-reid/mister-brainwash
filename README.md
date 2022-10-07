# Mr. Brainwash

Automating 'Life Remote Control'

## Build

Use make: `make build`

## Run

Also use make: `make run`

## Automate

Set up a cronjob on your desktop to upload a video every 6 hours: `sudo crontab -e`. Add the following:

```
0 5,11,17,23 * * * su -l sean -c /usr/local/bin/brainwash.sh sean >> /var/log/brainwash.log 2>&1
```

Replace `sean` with the name of your user, and modify the path to point to `brainwash.sh` in this repo (or move the shell script to `/usr/local/bin`, your choice).

## Troubleshooting

* If you run out of credits prior to auto-upload, the video is still saved in `videos/output`. Upload this manually until your credits reset the next day.
