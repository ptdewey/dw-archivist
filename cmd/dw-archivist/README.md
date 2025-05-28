# Discover Weekly Archivist

Spotify's Discover Weekly playlist is a great way to discover new music, but it only lasts for a week. This tool allows you to copy songs from your Discover Weekly playlist to another playlist, so you can keep the songs you like and listen to them later.

## Installation

With Go installed, you can install the archivist with the following command:
```bash
go install github.com/ptdewey/dw-archivist@latest
```

You can also build it from source:
```bash
git clone https://github.com/ptdewey/dw-archivist.git
cd dw-archivist
go build
```

## Usage

<!-- Runs in a cron job once a week on mondays (after new songs are added) -->

