# 17th-attendanceTracker

## Setup

**You will need a running MySQL or MariaDB instance.**

Create a database with a name of your choosing. Then, run the following SQL command against it to create a table.

```sql
-- a3server.attendancelog definition

CREATE TABLE `attendancelog` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `timestamp` datetime NOT NULL,
  `event_hash` varchar(100) NOT NULL,
  `event_type` varchar(100) NOT NULL,
  `player_id` varchar(30) NOT NULL,
  `player_uid` varchar(100) NOT NULL,
  `profile_name` varchar(100) NOT NULL,
  `steam_name` varchar(100) DEFAULT NULL,
  `is_jip` tinyint(4) DEFAULT NULL,
  `role_description` varchar(100) DEFAULT NULL,
  `mission_start` datetime NOT NULL,
  `mission_name` varchar(100) DEFAULT NULL,
  `briefing_name` varchar(100) DEFAULT NULL,
  `mission_name_source` varchar(100) DEFAULT NULL,
  `on_load_name` varchar(100) DEFAULT NULL,
  `author` varchar(100) DEFAULT NULL,
  `server_name` varchar(100) NOT NULL,
  `server_profile` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2713 DEFAULT CHARSET=utf8mb3;

-- a3server.`missions` definition

CREATE TABLE `missions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mission_name` varchar(100) NOT NULL,
  `mission_name_source` varchar(100) DEFAULT NULL,
  `briefing_name` varchar(100) DEFAULT NULL,
  `on_load_name` varchar(100) DEFAULT NULL,
  `author` varchar(100) DEFAULT NULL,
  `server_name` varchar(100) DEFAULT NULL,
  `server_profile` varchar(100) DEFAULT NULL,
  `mission_start` datetime DEFAULT NULL,
  `mission_hash` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;


-- a3server.`worlds` definition

CREATE TABLE `worlds` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` varchar(100) DEFAULT NULL,
  `display_name` varchar(100) DEFAULT NULL,
  `world_name` varchar(100) NOT NULL,
  `world_name_original` varchar(100) DEFAULT NULL,
  `world_size` int(11) DEFAULT NULL,
  `latitude` float DEFAULT NULL,
  `longitude` float DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3;

```

Finally, copy `config.example.json` to `config.json` and update it with your database credentials.
