package handlers

import (
	"github.com/PudgeKim/go-holdem/game"
	"github.com/gin-gonic/gin"
	"net/http"
)

var RequestChan = make(chan game.BetInfo)

// 프론트로부터 Json 요청 받음 (방 id, game.BetInfo 정보들 필요함)
// unmarshal 후에 방을 관리하는 채널에 넘겨줌
// 방을 관리하는 채널에서 방 번호에 해당하는 게임에 정보를 넘겨줌
// 게임을 관리하는 채널(game.betChan)에서 해당 요청을 처리
// game.betChan에서 요청 처리 후 다시 프론트로 넘겨줘야하는데
// 아래처럼 맨 처음 프론트에서 요청 받는 함수에서 <-game.betChan 이렇게 blocking하면 될듯?

func GameHandler(c *gin.Context) {
	var req game.BetInfo

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
