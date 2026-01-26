# 2GIS Friends Tracker

A real-time location tracker that monitors and stores the data of your 2GIS friends.

## Overview

2GIS is an offline map and navigation tool widely used in Russia, Kazakhstan, Uzbekistan, and other countries ([see the list here](https://en.wikipedia.org/wiki/2GIS)).

This application connects to the 2GIS API and tracks the locations of your friends in real-time by establishing a websocket connection. Other data includes battery status, speed, and whether a friend is charging their phone. These data are then stored in a local SQLite database for analysis (not yet implemented).

## Features

- **Real-time Location Tracking**: Connects to 2GIS API for live updates
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
ACCESS_TOKEN=your_2gis_access_token
APP_VERSION=your_app_version # 6.31.0 by default
USER_AGENT=your_user_agent
SITE_DOMAIN=https://2gis.ru # or kz/ae
DB_PATH=your_db_path # "./tracker.db" by default
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

To run this tracker, you will need your personal 2GIS access token. This token allows the application to authenticate as "you" and fetch the locations of your friends.

> **⚠️ Security Warning:** Your access token (like a password) is considered sensitive data. **Do not** share it publicly or commit it to GitHub. If you share this code, use a `.env` file or environment variables to keep it secret.

#### Step 1: Log in to 2GIS
1. Open your web browser (Chrome, Edge, or Firefox recommended)
2. Go to [2gis.kz](https://2gis.kz) (or 2gis.ru)
3. Log in to your account if you haven't already
4. Go to the **Friends** tab

#### Step 2: Open Developer Tools
1. Right-click anywhere on the page and select **Inspect** (or press `F12` / `Ctrl+Shift+I`)
2. In the developer window that opens, click on the **Network** tab at the top

#### Step 3: Capture the Token
1. In the **Network** tab, locate the **WS** tab (under the Disable cache checkbox)
2. Refresh the page (`F5` or `Ctrl+R`)
3. You will see many GET requests with 101 status appear

#### Step 4: Extract the Token
1. Double click on any of the GET requests
2. This should open a new tab. The link of the tab should end with `&token=<your_token>`
3. Copy only the `<your_token>` and paste it in `.env` file inside the`ACCESS_TOKEN` variable

> **TODO**: Document the process for generating 2GIS access tokens

## Usage

Run the tracker from the root directory:

```bash
go run ./cmd/tracker
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

### Features

- [ ] Display connection status in UI
- [ ] Add tabs to UI

### Dev related
- [ ] Document access token generation process
- [ ] Rewrite database logic to use goose and sqlc (if needed)
- [ ] Add license
- [ ] Implement graceful shutdown
- [ ] Add unit tests
- [ ] Add error recovery and reconnection logic
- [ ] Log metrics and monitoring to a file
- [ ] Save location data once in 10 minutes (custom)

## Privacy & Legal

**This project is for educational and research purposes only**

1.  **Unofficial:** This project is **not** affiliated with, endorsed by, or connected to 2GIS.
2.  **Terms of Service:** Using this tool may violate 2GIS's Terms of Service. Use it at your own risk. The author is not responsible if your account gets banned or restricted.
3.  **Privacy:** This tool allows tracking location data. It is the user's responsibility to ensure they have the consent of the individuals they are tracking and to comply with local privacy laws.
4.  **No Liability:** The software is provided "as is", without warranty of any kind. The author is not responsible for any damage, data loss, or legal consequences resulting from the use of this tool.

## License

[MIT license](https://github.com/rangaroo/2gis-friends/blob/main/LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## Author

[rangaroo](https://github.com/rangaroo)

