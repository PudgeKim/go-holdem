package channels

// BetInfo 게임이 시작되면 BetInfo 타입을 받는 채널로부터
// 요청이 올 때마다 베팅 관련 처리를 함
// 아래 같은 순서로 처리 됨
// 프론트 요청 -> BetInfo channel -> game에서 베팅 처리 -> Response 채널로 에러 또는 응답 반환 -> 프론트
type BetInfo struct {
	PlayerName string
	BetAmount  uint64
	IsDead     bool
}
