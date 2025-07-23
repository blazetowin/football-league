-- ===================================================
-- Football League Database Schema
-- ===================================================
-- This SQL script creates the necessary tables for a football league simulation system.
-- It includes tables for teams, match results, and calculated championship predictions.

-- ============================
-- Teams Table
-- ============================
CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,        -- Team name must be unique
    power INTEGER NOT NULL CHECK (power BETWEEN 1 AND 100) -- Power rating from 1 to 100
);

-- ============================
-- Matches Table
-- ============================
CREATE TABLE IF NOT EXISTS matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER NOT NULL CHECK (week >= 1 AND week <= 6), -- Week must be within valid range
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    home_goals INTEGER DEFAULT NULL CHECK (home_goals >= 0),
    away_goals INTEGER DEFAULT NULL CHECK (away_goals >= 0),
    FOREIGN KEY (home_team_id) REFERENCES teams(id),
    FOREIGN KEY (away_team_id) REFERENCES teams(id),
    CONSTRAINT unique_match UNIQUE (week, home_team_id, away_team_id) -- Prevent duplicate fixtures
);

-- ============================
-- Championship Predictions Table
-- ============================
CREATE TABLE IF NOT EXISTS championship_predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    chance REAL NOT NULL CHECK (chance >= 0 AND chance <= 100), -- Chance must be a percentage
    FOREIGN KEY (team_id) REFERENCES teams(id),
    CONSTRAINT unique_team_prediction UNIQUE (team_id)
);

-- ============================
-- Initial Data: Teams
-- ============================
INSERT INTO teams (name, power) VALUES 
('Chelsea', 90),
('Arsenal', 85),
('Manchester City', 88),
('Liverpool', 83)
ON CONFLICT(name) DO NOTHING; -- Prevent duplicate inserts if rerun
