package db

/*func TestInsert_Success(t *testing.T) {
	goose, dbmock, _ := sqlmock.New()
	defer goose.Close()

	repo := &GameRepository{goose: sqlx.NewDb(goose, "sqlmock"), logger: zaptest.NewLogger(t)}

	game := model.GameStats{
		PlayerStats: []model.PlayerStats{
			{PlayerID: "p1", Date: time.Now(), Points: 10, Rebounds: 5, Assists: 3, Steals: 2, Turnovers: 1, Blocks: 0, Fouls: 2, MinutesPlayed: 30},
			{PlayerID: "p2", Date: time.Now(), Points: 15, Rebounds: 8, Assists: 5, Steals: 1, Turnovers: 2, Blocks: 1, Fouls: 3, MinutesPlayed: 28},
		},
	}

	dbmock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO game_stats (id, game_id, date, player_id, points, rebounds, assist, steals, turnOvers, blocks, fouls, minutesPlayed) VALUES (?,?,?,?,?,?,?,?,?,?,?,?), (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "p1", 10, 5, 3, 2, 1, 0, 2, 30, // First row (11 values)
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "p2", 15, 8, 5, 1, 2, 1, 3, 28). // Second row (11 values)
		WillReturnResult(sqlmock.NewResult(1, 2))

	err := repo.Insert(game)
	assert.NoError(t, err)
	assert.NoError(t, dbmock.ExpectationsWereMet())

}*/
