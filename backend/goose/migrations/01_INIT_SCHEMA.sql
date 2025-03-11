-- +goose up
CREATE TABLE IF NOT EXISTS teams (
                                     id VARCHAR(36) PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     CONSTRAINT unique_team_name UNIQUE (name)
);

-- Players table
CREATE TABLE IF NOT EXISTS players (
                                       id VARCHAR(36) PRIMARY KEY,
                                       name VARCHAR(100) NOT NULL,
                                       team_id VARCHAR(36) NOT NULL,
                                       FOREIGN KEY (team_id) REFERENCES teams(id),
                                       CONSTRAINT unique_player_name_team_id UNIQUE (name, team_id)
);

-- Game stats table
CREATE TABLE IF NOT EXISTS game_stats (
                                          id varchar(36) PRIMARY KEY,
                                          game_id VARCHAR(36) NOT NULL,
                                          player_id VARCHAR(36) NOT NULL,
                                          date timestamp NOT NULL default current_timestamp,
                                          points INT NOT NULL,
                                          rebounds INT NOT NULL,
                                          assists INT NOT NULL,
                                          steals INT NOT NULL,
                                          blocks INT NOT NULL,
                                          fouls INT NOT NULL CHECK (fouls <= 6),
                                          turnovers INT NOT NULL,
                                          minutes_played FLOAT NOT NULL CHECK (minutes_played >= 0 AND minutes_played <= 48.0),
                                          FOREIGN KEY (player_id) REFERENCES players(id),
                                          INDEX idx_game_id (game_id),
                                          INDEX idx_player_id (player_id),
                                          INDEX idx_date (date)
);

CREATE OR REPLACE VIEW player_season_stats AS
SELECT
    p.id AS player_id,
    p.name AS player_name,
    p.team_id,
    t.name AS team_name,
    COUNT(DISTINCT gs.game_id) AS games_played,
    COALESCE(AVG(gs.points), 0) AS avg_points,
    COALESCE(AVG(gs.rebounds), 0) AS avg_rebounds,
    COALESCE(AVG(gs.assists), 0) AS avg_assists,
    COALESCE(AVG(gs.steals), 0) AS avg_steals,
    COALESCE(AVG(gs.blocks), 0) AS avg_blocks,
    COALESCE(AVG(gs.fouls), 0) AS avg_fouls,
    COALESCE(AVG(gs.turnovers), 0) AS avg_turnovers,
    COALESCE(AVG(gs.minutes_played), 0) AS avg_minutes_played
FROM
    players p
        LEFT JOIN
    game_stats gs ON p.id = gs.player_id
        JOIN
    teams t ON p.team_id = t.id
GROUP BY
    p.id, p.name, p.team_id, t.name;

-- Team season stats view
CREATE OR REPLACE VIEW team_season_stats AS
SELECT
    t.id AS team_id,
    t.name AS team_name,
    COUNT(DISTINCT gs.game_id) AS games_played,
    COALESCE(AVG(gs.points), 0) AS avg_points,
    COALESCE(AVG(gs.rebounds), 0) AS avg_rebounds,
    COALESCE(AVG(gs.assists), 0) AS avg_assists,
    COALESCE(AVG(gs.steals), 0) AS avg_steals,
    COALESCE(AVG(gs.blocks), 0) AS avg_blocks,
    COALESCE(AVG(gs.fouls), 0) AS avg_fouls,
    COALESCE(AVG(gs.turnovers), 0) AS avg_turnovers,
    COALESCE(AVG(gs.minutes_played), 0) AS avg_minutes_played
FROM
    teams t

        LEFT JOIN
    players p ON t.id = p.team_id
        LEFT JOIN
    game_stats gs ON p.id = gs.player_id
GROUP BY
    t.id, t.name;

