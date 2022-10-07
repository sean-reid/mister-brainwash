# Mr. Brainwash

Automating ['Life Remote Control'](https://www.youtube.com/channel/UCAWt5HjfXuJG4_7j3sSd9_Q).

## About

Ever since watching [Exit Through the Gift Shop](https://www.youtube.com/watch?v=IqVXThss1z4), I've been fascinated by Thierry Guetta (a.k.a. Mr. Brainwash) and his confusing compilation of seemingly random footage entitled ['Life Remote Control'](https://youtu.be/602RM3uFc_I). Naturally, as an engineer, the question arose: can I automate this? I turned to the world's largest collection of random footage for inspiration: YouTube. This project is an attemmpt to automate the random sampling, generation, and uploading of new content to YouTube, all in the spirit of Thierry's original work: chaos!

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
* You will need to follow the steps to authenticate your account (getting a token, etc) for Oath2 to work properly. Just once- after that initial authentication, the cronjob should work smoothly.
* More on the above: you'll need to download a JSON file from Google Cloud (save it in the `auth` folder here) that contains the client id and secret.
