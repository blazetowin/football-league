# Go Football League Simulator

A full-featured football league simulation backend built with Go and SQLite.
This project supports both CLI-based and REST API-based interaction for simulating matches, managing fixtures, viewing league tables, and generating championship predictions.

---

## Overview

* Simulates a 4-team league over 6 weeks (home and away)
* Team powers affect match scores (editable)
* Auto-generates fixtures with home/away balance
* View week-by-week progress in CLI or via HTTP endpoints
* SQLite database with schema auto-loaded on start
* Easily reset and customize league structure and results

---

## Project Structure

```
go-football-league/
├── cmd/
│   └── server.go            # API server entrypoint
├── internal/
│   ├── api/routes/          # HTTP route handlers
│   ├── domain/              # Data models
│   ├── league/              # Core simulation logic
│   │   ├── match.go
│   │   ├── predictor.go
│   │   ├── printer.go
│   │   ├── simulator.go
│   │   └── standings.go
│   ├── migration/
│   │   └── schema.sql       # SQL schema + initial teams
│   └── repository/
│       └── database.go      # DB connection and schema execution
├── league.db                # Auto-created SQLite database
├── main.go                  # CLI simulation runner
├── go.mod / go.sum
```

---

##  Running the Project

###  Option 1: CLI Simulation

Simulates each week step-by-step in terminal.

```bash
go run main.go
```

* Press `Enter` to go to the next week
* After Week 4, title predictions are printed below the table
* Match results and league table are printed every week

---

### Option 2: REST API Server

Launches HTTP server with API endpoints.

```bash
go run ./cmd/server.go
```

Server runs at: `http://localhost:8080`

---

## Database Reset / Customization

* Delete existing database:

  ```bash
  rm league.db
  ```
* Edit team powers in:
  `internal/migration/schema.sql` → bottom section
* Schema and initial data are auto-applied at startup

---

## Fixture Generation Logic

* Automatically generates 6-week fixtures using `CreateFixture()`
* Each team plays 3 home + 3 away matches
* You can adjust scoring advantage logic in `match.go` (line 70–71)

---

## API Endpoints

| Method | Endpoint                               | Description                                       |
| ------ | -------------------------------------- | ------------------------------------------------- |
| GET    | `/api/matches/{week}`                  | Simulate matches for a given week                 |
| GET    | `/api/league-table?week=3`             | Get league standings up to week 3                 |
| PUT    | `/api/match/{id}`                      | Manually update a match score                     |
| GET    | `/api/play-all-weeks`                  | Simulate and return all weeks at once             |
| GET    | `/api/week-summary?week=4`             | Summary of matches, table & predictions (week 4+) |
| GET    | `/api/championship-predictions/{week}` | Title probabilities (enabled after week 4)        |

---

## Example CLI Output

```bash
===== WEEK 3 =====

Match Results:
  Chelsea 2-1 Arsenal
  Liverpool 0-0 Manchester City

Standings:
Team           MP  W  D  L  GF  GA  GD  Pts
------------------------------------------
Chelsea        3   2  0  1   5   3  +2   6
...
```

---

## Tech Stack

* **Language:** Go (Golang)
* **Database:** SQLite (via `github.com/mattn/go-sqlite3`)
* **Architecture:** Clean folder structure (`cmd`, `internal`)
* **No external frameworks** — fully built on Go stdlib

---

## License

MIT © 2025 [Arda Olgun](https://github.com/blazetowin)

```}
```
