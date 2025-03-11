# Project Title

A brief description of what your project does and who it's for.

## Table of Contents

- [About](#about)
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## About

This serves a statistics of NBA players 

## Getting Started

### Prerequisites

Before running the project, make sure you have the following installed:

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- Mysql if not using docker

### Installation

To install the project, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/ireuven89/skyhawk.git
    ```
2. If running on local :
   if running on local
   set 2 envs
   MYSQL_ROOT_PASSWORD
   MYSQL_HOST
   REDIS_HOST

    navigate to the project directory
    ```bash
    cd skyhawk/backend
    run go build .
    run /.backend
    ```


3. If using Docker,run the docker compose command:

    ```bash
   cd skyhawk
   run docker compose --env-file app.env up -d app
    ```

## Usage

The project serves 4 main APIs
  1. log game stats post - POST /api/v1/games/log
    example request is:
      
         \\json
               {     
         "id": "Warriors vs lakers", - can be empty
         "date": "2025-03-08T19:30:00Z",
         "teams": [
               {
               "name": "Lakers",
               "players": [
               {
               "player_name": "LeBron James",
               "points": 30,
               "rebounds": 12,
               "assists": 8,
               "steals": 2,
               "blocks": 1,
               "fouls": 2,
               "turnovers": 3,
               "minutes_played": 38.5
               }
               ]
               },
         {
               "name": "Warriors",
               "players": [
               {
               "player_name": "Stephen Curry",
               "points": 35,
               "rebounds": 5,
                "assists": 6,
                 "steals": 1,
             "blocks": 0,
             "fouls": 1,
                   "turnovers": 2,
               "minutes_played": 40
               }
               ]
               }
         ]
         }
  2. fetch player season stats GET players/season/:player_id
  3. fetch team season stats GET /teams/stats/season/:team_id
  4. fetch game stats GET /games/:game_id
     returns all the player stats of a game

  5. Deployment on AWS:
  open an ecr with the project name
  install aws cli on your local machine
  build the image docker build -t skyhawk .
  tag the image with the image name and aws id
  push the image to the remote registry 
   