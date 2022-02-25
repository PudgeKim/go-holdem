package game

type BetInfoRequest struct {
	PlayerName string `json:"player_name"`
	BetAmount  uint64 `json:"bet_amount"`
	IsDead     bool   `json:"is_dead"`
}
