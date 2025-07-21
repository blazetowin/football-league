-- TEAMS
-- This SQL script creates the necessary tables for a football league database.
-- It includes tables for teams, matches, and championship predictions.
CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    power INTEGER NOT NULL
);

-- MATCHES
-- This table records the matches played between teams, including the week of the match, the teams involved, and the goals scored by each team.
-- The home_team_id and away_team_id are foreign keys referencing the teams table.
-- The home_goals and away_goals fields store the number of goals scored by each team in the match.
CREATE TABLE IF NOT EXISTS matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER NOT NULL,
    home_team_id INTEGER,
    away_team_id INTEGER,
    home_goals INTEGER DEFAULT NULL,
    away_goals INTEGER DEFAULT NULL,
    FOREIGN KEY (home_team_id) REFERENCES teams(id),
    FOREIGN KEY (away_team_id) REFERENCES teams(id)
);

-- CHAMPIONSHIP PREDICTIONS
-- This table stores predictions for the championship chances of each team.
-- Each prediction is associated with a team and includes a chance value representing the likelihood of that team winning the championship.
-- The team_id is a foreign key referencing the teams table.
CREATE TABLE IF NOT EXISTS championship_predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER,
    chance REAL,
    FOREIGN KEY (team_id) REFERENCES teams(id)
);
-- This section inserts sample data into the teams table.
INSERT INTO teams (name, power) VALUES 
('Chelsea', 90),
('Arsenal', 85),
('Manchester City', 88),
('Liverpool', 83);
