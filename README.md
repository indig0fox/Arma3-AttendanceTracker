# 17th-attendanceTracker

## Setup

**You will need a running MySQL or MariaDB instance.**

Create a database with a name of your choosing. Then, run the following SQL command against it to create a table.

```sql
-- a3server.attendancelog definition

CREATE TABLE `attendance` (
 `id` INT(11) NOT NULL AUTO_INCREMENT,
 `join_time` DATETIME NULL DEFAULT NULL,
 `disconnect_time` DATETIME NULL DEFAULT NULL,
 `mission_hash` VARCHAR(100) NULL DEFAULT '' COLLATE 'utf8mb3_general_ci',
 `event_type` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
 `player_id` VARCHAR(30) NOT NULL COLLATE 'utf8mb3_general_ci',
 `player_uid` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
 `profile_name` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
 `steam_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `is_jip` TINYINT(4) NULL DEFAULT NULL,
 `role_description` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb3_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=5868
;


-- a3server.`missions` definition

CREATE TABLE `missions` (
 `id` INT(11) NOT NULL AUTO_INCREMENT,
 `mission_name` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
 `mission_name_source` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `briefing_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `on_load_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `author` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `server_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `server_profile` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `mission_start` DATETIME NULL DEFAULT NULL,
 `mission_hash` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb3_general_ci'
ENGINE=InnoDB
;



-- a3server.`worlds` definition

CREATE TABLE `worlds` (
 `id` INT(11) NOT NULL AUTO_INCREMENT,
 `author` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `display_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `world_name` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
 `world_name_original` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 `world_size` INT(11) NULL DEFAULT NULL,
 `latitude` FLOAT NULL DEFAULT NULL,
 `longitude` FLOAT NULL DEFAULT NULL,
 `workshop_id` VARCHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
 PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb3_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=2
;

```

Finally, copy `config.example.json` to `config.json` and update it with your database credentials.
