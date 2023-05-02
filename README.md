# 17th-attendanceTracker

## Setup

**You will need a running MySQL or MariaDB instance.**

Create a database with a name of your choosing. Then, run the following SQL command against it to create a table.

```sql
CREATE TABLE `attendancelog` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `timestamp` DATETIME NOT NULL,
    `event_hash` VARCHAR(100) NOT NULL DEFAULT md5(concat(`server_name`,`mission_name`,`author`,`mission_start`)) COLLATE 'utf8mb3_general_ci',
    `event_type` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
    `player_id` VARCHAR(30) NOT NULL COLLATE 'utf8mb3_general_ci',
    `player_uid` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
    `profile_name` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
    `steam_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `is_jip` TINYINT(4) NULL DEFAULT NULL,
    `role_description` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `mission_start` DATETIME NOT NULL,
    `mission_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `briefing_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `mission_name_source` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `on_load_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `author` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb3_general_ci',
    `server_name` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
    `server_profile` VARCHAR(100) NOT NULL COLLATE 'utf8mb3_general_ci',
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb3_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=383
;
```

Finally, copy `config.example.json` to `config.json` and update it with your database credentials.