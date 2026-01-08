# 2GIS Friends Tracker

A real-time location tracker for 2GIS friends that monitors and stores your friends' location data via WebSocket connection.

## Overview

This application connects to the 2GIS WebSocket API to track your friends' locations in real-time. It receives location updates, battery status, and other state information, displaying them in the console and storing them in a SQLite database for historical analysis.

## Features

- **Real-time Location Tracking**: Connects to 2GIS WebSocket API for live updates
- **Friend Profile Management**: Caches friend profiles for quick lookups
- **Location History**: Stores all location updates in a SQLite database


## Prerequisites

- Go 1.25.4 or higher
- 2GIS account with access token

## Installation

1. Clone the repository:
```bash
git clone https://github.com/rangaroo/2gis-friends.git
cd 2gis-friends
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the `.env.example` file to `.env` and fill in the values.

```env
ACCESS_TOKEN= your_2gis_access_token
APP_VERSION=  your_app_version # 6.31.0 by default
USER_AGENT=   your_user_agent
SITE_DOMAIN=  https://2gis.ru # or kz/ae
DB_PATH=      your_db_path # "./tracker.db" by default
```

## Configuration

The application requires the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `ACCESS_TOKEN` | Your 2GIS API access token | `eyJhbGc...` |
| `APP_VERSION` | 2GIS app version | `6.31.0` |
| `USER_AGENT` | User agent string for WebSocket connection | `Mozilla/5.0...` |
| `SITE_DOMAIN` | 2GIS site domain | `https://2gis.ru` |
| `DB_PATH` | Path to SQLite database file | `./tracker.db` |

### Obtaining Access Token

To run this tracker, you need your personal 2GIS access token. This token allows the application to authenticate as "you" and fetch your friends' locations.

> **⚠️ Security Warning:** Your access token is sensitive data (like a password). **Do not** share it publicly or commit it to GitHub. If you share this code, use a `.env` file or environment variables to keep it secret.

#### Step 1: Log in to 2GIS
1. Open your web browser (Chrome, Edge, or Firefox recommended).
2. Go to [2gis.kz](https://2gis.kz) (or 2gis.ru).
3. Log in to your account if you haven't already.
4. Go to the **friends tab**

#### Step 2: Open Developer Tools
1. Right-click anywhere on the page and select **Inspect** (or press `F12` / `Ctrl+Shift+I`).
2. In the developer window that opens, click on the **Network** tab at the top.

#### Step 3: Capture the Token
1. In the **Network** tab, locate the **WS** tab (under the disable cache checkbox)
2. Refresh the page (`F5` or `Ctrl+R`)
3. You will see many GET requests with 101 status appear

#### Step 4: Extract the Token
1. Double click on any of the GET requests
2. This should open a new tab. The link of the tab should end with `&token=<your_token>`
3. Copy only the `<your_token>` and paste it in `.env` file inside the`ACCESS_TOKEN` variable

> **TODO**: Document the process for generating 2GIS access tokens

Currently, you need to obtain an access token from the 2GIS API.
1. Authenticating with your 2GIS account
2. Extracting the access token from the authenticated session

## Usage

Run the tracker from the `cmd/tracker` directory:

```bash
cd cmd/tracker
go run .
```

Or build and run the binary:

```bash
go build -o tracker ./cmd/tracker
./tracker
```

## Database Schema

The application creates a `locations` table with the following structure:

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER | Primary key (auto-increment) |
| `user_id` | TEXT | Friend's userID on 2GIS |
| `lat` | REAL | Latitude |
| `lon` | REAL | Longitude |
| `accuracy` | REAL | GPS accuracy in meters |
| `speed` | REAL | Movement speed |
| `battery_level` | REAL | Battery level (0.0-1.0) |
| `is_charging` | BOOLEAN | Charging status |
| `timestamp` | INTEGER | Unix timestamp (milliseconds) |


## Known Issues & TODOs

- [ ] Document access token generation process
- [ ] Rewrite database logic to use goose and sqlc
- [ ] Add license
- [ ] Implement graceful shutdown
- [ ] Add unit tests
- [ ] Add error recovery and reconnection logic
- [ ] Add metrics and monitoring
- [ ] Add CLI flags for runtime configuration

## Privacy & Legal

⚠️ **Important**: This tool tracks location data of other users. Ensure you have:
- Proper authorization to track friends' locations
- Compliance with local privacy laws and regulations
- Consent from tracked individuals where required
- Understanding of 2GIS Terms of Service

Use this tool responsibly and ethically.

## License

[TODO]

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## Author

[rangaroo](https://github.com/rangaroo)
