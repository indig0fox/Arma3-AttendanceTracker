# Arma 3 Attendance Tracker

---

## Setup

### Set Up a Database Engine

**You will need a running MySQL or MariaDB instance.**

If you do not have a MySQL or MariaDB instance, you can use one of the following:

- [MySQL Community Server](https://dev.mysql.com/downloads/mysql/) (free)
- [MariaDB](https://mariadb.org/) (free)
- [Docker](https://www.docker.com/) (free)
  - [MySQL Docker Image](https://hub.docker.com/_/mysql)
  - [MariaDB Docker Image](https://hub.docker.com/_/mariadb)

The tested and recommended database engine is MariaDB version 11.0.2.

### Create a Database

You will need to create a database in your instance. You can do this using your preferred MySQL client or one of the following:

- [phpMyAdmin](https://www.phpmyadmin.net/) (web based)
- [HeidiSQL](https://www.heidisql.com/) (Windows only and simpler)
- [MySQL Workbench](https://www.mysql.com/products/workbench/) (feature rich but heavy)
- [DBeaver](https://dbeaver.io/) (feature rich but heavy)

You could also use a CLI client such as the `mysql` command line client to run the following command:

```sql
CREATE DATABASE `arma3_attendance`;
```

### Mod Installation

1. Download the latest release from the [releases page](https://github.com/indig0fox/Arma3-AttendanceTracker/releases).
1. Extract the .zip and move `@AttendanceTracker` to your Arma 3 server's root directory.
1. Inside of `@AttendanceTracker` you will find a `config.json` file. Open this file and configure it to your circumstances. See the [Configuration](#configuration) section for more information.
1. Add the mod to your server's startup parameters. For example: `-serverMod="@AttendanceTracker;"`

At next run, the Arma 3 server will launch with the mod running.

### Configuration

The configuration file example is located at `@AttendanceTracker/config.json.example`. This should be copied to `@AttendanceTracker/config.json` and edited to suit your circumstances.

The following table describes the configuration options.

| Key | Type | Description | Default |
| --- | --- | --- | --- |
| sqlConfig.mySqlHost | string | The hostname of your MySQL instance. | localhost |
| sqlConfig.mySqlPort | integer | The port of your MySQL instance. | 3306 |
| sqlConfig.mySqlUser | string | The username to use when connecting to your MySQL instance. | root |
| sqlConfig.mySqlPassword | string | The password to use when connecting to your MySQL instance. | root |
| sqlConfig.mySqlDatabase | string | The name of the database to use. | arma3_attendance |
| armaConfig.dbUpdateIntervalSeconds | integer | The number of seconds between disconnect_time updates per user. | 90 |
| armaConfig.serverEventFillNullMinutes | integer | The max session duration to fill in for missing server disconnect_time values. | 90 |
| armaConfig.missionEventFillNullMinutes | integer | The max session duration to fill in for missing mission disconnect_time values. | 15 |
| armaConfig.debug | boolean | Whether or not to enable debug logging. | false |

## Usage

The extension uses [GORM](https://gorm.io/) for database access and will automatically create the schema in the database you specify in the configuration file.

---

## Important Notes

### Logging

"debug": true:
The extension will log ERROR and WARN events to the Arma 3 server's RPT file, which can be found in the server's profile folder.

"debug": false:
The extension will log ERROR, WARN, INFO, and DEBUG events to the Arma 3 server's RPT file, which can be found in the server's profile folder.

All events will always be logged to `@AttendanceTracker/attendanceTracker.log` in log line format.

### Timezone

All times will be logged as UTC time. This is to ensure that all times are logged in a consistent manner, regardless of the timezone of the server. Because these are DATETIME fields, they will not be adjusted to your local time when viewing them in a database client.

To do so, you can use the following function in your queries:

```sql
CONVERT_TZ(<field>, 'UTC', 'US/Eastern')
```

A full list of timezones available to your database can be found this way:

```sql
SELECT * 
FROM mysql.time_zone_name
```

### Performance

The extension will update the disconnect_time field for each player every `dbUpdateIntervalSeconds` seconds. This is to ensure that the disconnect_time field is updated in the event that the server crashes or the mission ends without a disconnect event.

These calls are threaded in the Go runtime and will not block the Arma 3 server while processing. The default value of 90 seconds should be sufficient for most servers. Each period begins when a player connects to the server or connects to a mission, which provides a natural offset.

### NULL disconnect_time Values

In the event that the server crashes or a disconnect event for a mission is not sent, the next join for each will update past rows based on the following:

If the join time for a row is within [`{event_type}EventFillNullMinutes`](#configuration) minutes of the previous disconnect time, the previous disconnect time will be updated to the new join time. Otherwise, it will be set as [`{event_type}EventFillNullMinutes`](#configuration) from the join time for that row.

This is an attempt to account for missing events for individual players while not attributing large gap periods to their calculated session times. If the server crashes, the extension will update all rows with a NULL disconnect_time to the current time. See [Server Crash Time Filling](#server-crash-time-filling) for more information.

#### Server Crash Time Filling

The addon will update `@AttendanceTracker/lastServerTime.txt` with Arma 3's `diag_tickTime` every 30 seconds. This is to ensure that the server time is always available to the extension, even if the server crashes. This file is not used for any other purpose.

On each time update, the extension will check this file and compare the received value to it. If the lastServerTime < lastServerTime.txt, the extension will assume that the server has restarted and will update all event rows with a NULL disconnect_time to the current time OR the threshold specified in the configuration file, whichever produces the smaller session duration.

---

## Database Schema

| Table Name | Description |
| --- | --- |
| worlds | Stores world information. |
| missions | Stores mission information. |
| attendance_items | Stores rows that indicate player information and join/disconnect times. |

### Worlds

The worlds table will store basic info about the world. This is used to link missions to worlds.

### Missions

The missions table will store basic info about the mission. This is used to link attendance items to missions.

### Attendance Items

The attendance_items table will store rows that indicate player information and join/disconnect times. This can be used to calculate play time per player per mission. Each row is also linked to a mission, so that these records can be grouped.

---

## Useful Queries

### Show missions with attendance

This will retrieve a view showing all missions with attendance data, sorted by the most recent mission joins first. Mission events without a mission disconnect_time (due to server crash or in-progress mission) will be ignored.

See [Timezone](#timezone) for more information on converting times to your local timezone.

```sql
select
    a.server_profile as Server,
    a.briefing_name as "Mission Name",
    CONVERT_TZ(a.mission_start, 'UTC', 'US/Eastern') as "Start Time",
    b.display_name as "World",
    c.profile_name as "Player Name",
    c.player_uid as "Player UID",
    TIMESTAMPDIFF(
        MINUTE,
        c.join_time,
        c.disconnect_time
    ) as "Play Time (m)",
    CONVERT_TZ(c.join_time, 'UTC', 'US/Eastern') as "Join Time",
    CONVERT_TZ(c.disconnect_time, 'UTC', 'US/Eastern') as "Leave Time"
from missions a
    LEFT JOIN worlds b ON a.world_id = b.id
    LEFT JOIN attendance_items c ON a.mission_hash = c.mission_hash
where
    c.event_type = 'Mission'
    AND c.disconnect_time IS NOT NULL
    AND TIMESTAMPDIFF(
        MINUTE,
        c.join_time,
        c.disconnect_time
    ) > 0
```

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## Development

### Prerequisites

- [Go 1.16.4](https://golang.org/doc/install)
- [MinGW-w64](https://sourceforge.net/projects/mingw-w64/) (Windows only)
- [GCC](https://gcc.gnu.org/) (Linux only)

### Building Extension using Docker

You will need Docker Engine installed and running. This can be done on Windows or on Linux. However, you will need to use Linux containers if you're on Windows (specified in Docker Desktop settings).

Once it's built, copy the file from ./dist to the project root, then build the addon.

#### COMPILING FOR WINDOWS

```ps1
docker pull x1unix/go-mingw:1.20

# version is semantic + build date + git hash
# e.g. 1.0.0-2021-05-30-1a2b3c4d
$versionSem = '1.1.0'
$dateStr = Get-Date -Format 'yyyyMMdd'
$version = "$versionSem-$dateStr-$(git rev-parse --short HEAD)"

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o dist/AttendanceTracker_x64.dll -buildmode=c-shared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" ./cmd

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o dist/AttendanceTracker.dll -buildmode=c-shared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" ./cmd

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o dist/AttendanceTracker_x64.exe -ldflags "-w -s -X main.EXTENSION_VERSION=$version" ./cmd
```

#### COMPILING FOR LINUX

```bash
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build .

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/AttendanceTracker_x64.so -linkshared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" ./cmd

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o dist/AttendanceTracker.so -linkshared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" ./cmd
```

### Compile Addon

First, move the compiled dlls from extension/AttendanceTracker/dist to the project root.

To prepare the addon, you'll need to download the [HEMTT](https://brettmayson.github.io/HEMTT/commands/build.html) binary, place it in the project root, and run the following command:

```bash
./HEMTT.exe release
```

The PBOs and relevant files will be placed in the ./.hemmttout/build directory.

---

## Credits

- [Go](https://golang.org/)
- [Go Arma 3 Extension Template](https://github.com/code34/armago_x64/tree/master)
- [GORM](https://gorm.io/)