package gameservice

// 베팅 관련 처리를 한 후 프론트로 베팅처리결과 전달
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
	Winners 		[]string `json:"winners,omitempty"`
}

type GameStartResponse struct {
	ReadyPlayers []string 
	FirstPlayer string 
	SmallBlind string 
	BigBlind string 
}

func NewGameStartResponse(readyPlayers []string, firstPlayer, smallBlind, bigBlind string) *GameStartResponse {
	return &GameStartResponse{
		readyPlayers,
		firstPlayer,
		smallBlind,
		bigBlind,
	}
}