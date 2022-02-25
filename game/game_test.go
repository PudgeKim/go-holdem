package game

import (
	"errors"
	"fmt"
	"github.com/PudgeKim/player"
	"reflect"
	"testing"
)

func TestGame_Bet(t *testing.T) {
	var players []*player.Player
	var betInfo BetInfoRequest

	kim := player.New("kim", 300, 200)
	park := player.New("park", 300, 200)
	han := player.New("han", 300, 200)
	lee := player.New("lee", 300, 200)
	choi := player.New("choi", 300, 200)

	players = []*player.Player{&kim, &park, &han, &lee, &choi}

	game := New(players)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "anonymous"
	betInfo.IsDead = false

	isBetFinished, err := game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if !errors.Is(err, NoPlayerExists) {
		t.Error("NoPlayerExistsError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "kim"
	betInfo.IsDead = false

	isBetFinished, err = game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if !errors.Is(err, InvalidPlayerTurn) {
		t.Error("InvalidPlayerTurnError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.players[2].IsDead = true

	isBetFinished, err = game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if !errors.Is(err, DeadPlayer) {
		t.Error("DeadPlayerError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.players[2].IsDead = false
	game.players[2].IsLeft = true

	isBetFinished, err = game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if !errors.Is(err, PlayerLeft) {
		t.Error("PlayerLeftError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = true
	game.players[2].IsDead = false
	game.players[2].IsLeft = false

	isBetFinished, err = game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if err != nil {
		t.Error("there should be no error")
	}
	if game.players[2].IsDead != true {
		t.Error("han's IsDead field should be set to true")
	}

	// 여기부터는 정상 베팅
	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.players[2].IsDead = false

	isBetFinished, err = game.Bet(betInfo)
	if isBetFinished {
		t.Error("it should be false")
	}
	if err != nil {
		t.Error("there should be no error")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "lee"
	betInfo.IsDead = false

	game.Bet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "choi"
	betInfo.IsDead = false

	game.Bet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "kim"
	betInfo.IsDead = false

	game.Bet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "park"
	betInfo.IsDead = false

	isBetFinished, err = game.Bet(betInfo)
	if !isBetFinished {
		t.Error("bet should be finished")
	}
	if err != nil {
		t.Error("there should be no error")
	}
}

func TestGame_isValidBet(t *testing.T) {
	var p player.Player

	p = player.New("kim", 300, 200)
	game := New([]*player.Player{&p})

	game.currentBet = 20
	betType, err := game.isValidBet(&p, 20)
	if betType != Check {
		t.Error("betType should be Check")
	}
	if err != nil {
		t.Error("there should be no error")
	}

	game.currentBet = 20
	p = player.New("kim", 300, 200)
	betType, err = game.isValidBet(&p, 30)
	if betType != Raise {
		t.Error("betType should be Raise")
	}
	if err != nil {
		t.Error("there should be no error")
	}

	p.TotalBet = 30
	betType, err = game.isValidBet(&p, 170)
	if betType != AllIn {
		t.Error("betType should be AllIn")
	}
	if err != nil {
		t.Error("there should be no error")
	}

	game.currentBet = 20
	p = player.New("kim", 300, 200)
	betType, err = game.isValidBet(&p, 250)
	if !errors.Is(err, OverBalance) {
		t.Error("OverBalanceError should be occurred")
	}

	game.currentBet = 20
	p = player.New("kim", 300, 200)
	betType, err = game.isValidBet(&p, 10)
	if !errors.Is(err, LowBetting) {
		t.Error("LowBettingError should be occurred")
	}
}

func TestGame_getNextPlayerIdx(t *testing.T) {
	var players []*player.Player

	players = []*player.Player{
		{Nickname: "kim", IsDead: false, IsLeft: false},
		{Nickname: "park", IsDead: false, IsLeft: false},
		{Nickname: "han", IsDead: true, IsLeft: false},
		{Nickname: "lee", IsDead: false, IsLeft: true},
		{Nickname: "choi", IsDead: false, IsLeft: false},
	}

	game := New(players)
	nextIdx, err := game.getNextPlayerIdx()
	if err != nil {
		t.Error("there should be no error")
	}

	if nextIdx != 4 {
		fmt.Println(nextIdx)
		t.Error("next index should be 4")
	}

	players = []*player.Player{
		{Nickname: "kim", IsDead: true, IsLeft: false},
		{Nickname: "park", IsDead: true, IsLeft: false},
		{Nickname: "han", IsDead: true, IsLeft: false},
		{Nickname: "lee", IsDead: false, IsLeft: true},
		{Nickname: "choi", IsDead: true, IsLeft: false},
	}

	game = New(players)
	nextIdx, err = game.getNextPlayerIdx()
	if !errors.Is(err, NoPlayersLeft) {
		t.Error("there should be NoPlayersLeftError")
	}
}

func TestGame_AddPlayer(t *testing.T) {
	kim := player.Player{Nickname: "kim"}
	han := player.Player{Nickname: "han"}

	players := []*player.Player{
		&kim,
		&han,
	}

	game := New(players)
	park := player.Player{Nickname: "park"}
	game.AddPlayer(&park)

	if !reflect.DeepEqual(game.GetAllPlayers(), []*player.Player{
		&kim,
		&han,
		&park,
	}) {
		t.Error("players should be kim, han, park")
	}
}

func TestGame_RemovePlayer(t *testing.T) {
	kim := player.Player{Nickname: "kim"}
	han := player.Player{Nickname: "han"}
	park := player.Player{Nickname: "park"}

	players := []*player.Player{&kim, &han, &park}

	game := New(players)
	game.RemovePlayer(park)

	if !reflect.DeepEqual(game.GetAllPlayers(), []*player.Player{&kim, &han}) {
		t.Error("players should be kim, han")
	}
}
