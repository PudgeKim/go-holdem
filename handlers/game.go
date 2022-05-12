package handlers

import (
	"context"
	"net/http"

	"github.com/PudgeKim/go-holdem/channels"
	"github.com/PudgeKim/go-holdem/game"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var betResponseChanMap = make(map[uuid.UUID]chan channels.BetResponse)

type GameHandler struct {
	rooms       gamerooms.GameRooms
}

func NewGameHandler(rooms gamerooms.GameRooms) *GameHandler {
	return &GameHandler{
		rooms:       rooms,
	}
}

type GameStartReq struct {
	RoomId string `json:"room_id"`
	HostId string `json:"host_id"`
}

func (g GameHandler) Start(c *gin.Context) {
	var req GameStartReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := g.rooms.GetGameRoomAfterParse(req.RoomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// grpc 통신을 통해 유효한 방장 Id인지 확인
	host, err := g.grpcHandler.GetUser(context.Background(), req.HostId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if host.Id != req.HostId {
		c.JSON(http.StatusBadRequest, gin.H{"error": gameerror.InvalidHostId})
		return
	}

	if room.Game.IsStarted {
		c.JSON(http.StatusBadRequest, gin.H{"error": gameerror.AlreadyStarted})
		return
	}

	if err := room.Game.SetPlayers(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// blocking되면 안되므로 고루틴으로 실행 (Start함수는 무한 for loop)
	go room.Game.Start()

	statusOkWithSuccess(c, nil, nil)
}

type BetReq struct {
	RoomId    string `json:"room_id"`
	PlayerId  string `json:"player_id"`
	BetAmount uint64 `json:"bet_amount"`
	IsDead    bool   `json:"is_dead"`
}

func (g GameHandler) Bet(c *gin.Context) {
	var req BetReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestWithError(c, err)
		return
	}
	
	room, err := g.rooms.GetGameRoomAfterParse(req.RoomId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	user, err := g.grpcHandler.GetUser(context.Background(), req.PlayerId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	betInfo := channels.BetInfo{
		PlayerName: user.Name,
		BetAmount:  req.BetAmount,
		IsDead:     req.IsDead,
	}

	// Game은 무한 for loop을 돌며 베팅 정보가 들어올 때마다 처리함
	room.Game.BetChan <- betInfo

	// 게임 내에서 베팅처리를 한 후 돌아오는 응답을 기다림
	betResponse := <-betResponseChanMap[room.Id]

	if betResponse.Error != nil {
		badRequestWithError(c, betResponse.Error)
		return
	}

	// 게임 종료시 승자들 잔고 올려주고 패자들 잔고 내려야함
	if betResponse.GameStatus == game.GameEnd {
		// TODO
	}

	statusOkWithSuccess(c, nil, betResponse)
}

func SetBetResponseChan(roomId uuid.UUID) chan channels.BetResponse {
	betResponseChanMap[roomId] = make(chan channels.BetResponse)
	return betResponseChanMap[roomId]
}
