package channels

// 베팅 관련 처리를 한 후 게임으로부터 온 응답을 받아서
// 프론트로 전달해주기 위해 필요한 채널
type BetResponse struct {
	Error            error  `json:"error"`
	IsBetEnd         bool   `json:"is_bet_end"` // true면 플레이어들의 베팅이 모두 끝나서 다음 턴으로 넘어감
	IsPlayerDead     bool   `json:"is_player_dead"`
	PlayerCurrentBet uint64 `json:"player_current_bet"`
	PlayerTotalBet   uint64 `json:"player_total_bet"`
	GameCurrentBet   uint64 `json:"game_current_bet"`
	GameTotalBet     uint64 `json:"game_total_bet"`
	NextPlayerName   string `json:"next_player_name"`
	GameStatus       string `json:"game_status"` // FreeFlop, Flop, Turn, River
}
