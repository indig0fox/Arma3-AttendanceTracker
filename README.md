# Arma 3 Attendance Tracker

## Setup

**You will need a running MySQL or MariaDB instance.**

The following SQL commands will set up the necessary tables for the application. You can run them from the MySQL command line or from a tool like phpMyAdmin.

*In future, an ORM will be used to set this up automatically.*

```sql
CREATE DATABASE `arma3_attendance` /*!40100 DEFAULT CHARACTER SET utf8mb3 */;

USE `arma3_attendance`;

-- a3server.missions definition
CREATE TABLE `missions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `world_id` int(11) DEFAULT NULL,
  `mission_hash` varchar(100) NOT NULL DEFAULT '',
  `mission_name` varchar(100) NOT NULL,
  `mission_name_source` varchar(100) DEFAULT NULL,
  `briefing_name` varchar(100) DEFAULT NULL,
  `on_load_name` varchar(100) DEFAULT NULL,
  `author` varchar(100) DEFAULT NULL,
  `server_name` varchar(100) DEFAULT NULL,
  `server_profile` varchar(100) DEFAULT NULL,
  `mission_start` datetime DEFAULT NULL COMMENT 'In UTC',
  PRIMARY KEY (`id`),
  KEY `mission_hash` (`mission_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb3;

-- arma3_attendance.attendance definition
CREATE TABLE `attendance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `join_time` datetime DEFAULT NULL COMMENT 'Stored in UTC',
  `disconnect_time` datetime DEFAULT NULL COMMENT 'Stored in UTC',
  `mission_hash` varchar(100) DEFAULT NULL,
  `event_type` varchar(100) NOT NULL,
  `player_id` varchar(30) NOT NULL,
  `player_uid` varchar(100) NOT NULL,
  `profile_name` varchar(100) NOT NULL,
  `steam_name` varchar(100) DEFAULT NULL,
  `is_jip` tinyint(4) DEFAULT NULL,
  `role_description` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `mission_hash` (`mission_hash`),
  CONSTRAINT `attendance_ibfk_1` FOREIGN KEY (`mission_hash`) REFERENCES `missions` (`mission_hash`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb3;


-- a3server.worlds definition
CREATE TABLE `worlds` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` varchar(100) DEFAULT NULL,
  `display_name` varchar(100) DEFAULT NULL,
  `world_name` varchar(100) NOT NULL,
  `world_name_original` varchar(100) DEFAULT NULL,
  `world_size` int(11) DEFAULT NULL,
  `latitude` float DEFAULT NULL,
  `longitude` float DEFAULT NULL,
  `workshop_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `world_name` (`world_name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb3;
```

Finally, copy `config.example.json` to `config.json` and update it with your database credentials and path.

## QUERIES

### Show missions with attendance

This will retrieve a view showing all missions with attendance data, sorted by the most recent mission joins first. Mission events without a mission disconnect_time (due to server crash or in-progress mission) will be ignored.

```sql
select
    a.server_profile as Server,
    a.briefing_name as "Mission Name",
    a.mission_start as "Start Time",
    b.display_name as "World",
    c.profile_name as "Player Name",
    c.player_uid as "Player UID",
    TIMESTAMPDIFF(
        MINUTE,
        c.join_time,
        c.disconnect_time
    ) as "Play Time (m)",
    c.join_time as "Join Time",
    c.disconnect_time as "Leave Time"
from missions a
    LEFT JOIN worlds b ON a.world_id = b.id
    LEFT JOIN attendance c ON a.mission_hash = c.mission_hash
where
    c.event_type = 'Mission'
    AND c.disconnect_time IS NOT NULL
    AND TIMESTAMPDIFF(
        MINUTE,
        c.join_time,
        c.disconnect_time
    ) > 0
```
