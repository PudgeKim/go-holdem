package game

import (
	"errors"
	"fmt"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/player"
	"testing"
)

func TestGame_Bet(t *testing.T) {
	var players []*player.Player
	var betInfo BetInfo

	kim := player.New("kim", 300, 200)
	park := player.New("park", 300, 200)
	han := player.New("han", 300, 200)
	lee := player.New("lee", 300, 200)
	choi := player.New("choi", 300, 200)

	kim.IsReady = true
	park.IsReady = true
	han.IsReady = true
	lee.IsReady = true
	choi.IsReady = true

	players = []*player.Player{&kim, &park, &han, &lee, &choi}

	game := New()
	game.SetPlayers(players)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "anonymous"
	betInfo.IsDead = false

	nextPlayer, isBetEnd, err := game.HandleBet(betInfo)
	if nextPlayer != "" {
		t.Error("nextPlayer should be empty")
	}

	if isBetEnd {
		t.Error("it should be false")
	}
	if !errors.Is(err, gameerror.NoPlayerExists) {
		t.Error("NoPlayerExistsError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "kim"
	betInfo.IsDead = false

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if isBetEnd {
		t.Error("it should be false")
	}
	if !errors.Is(err, gameerror.InvalidPlayerTurn) {
		t.Error("InvalidPlayerTurnError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.Players[2].IsDead = true

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if isBetEnd {
		t.Error("it should be false")
	}
	if !errors.Is(err, gameerror.DeadPlayer) {
		t.Error("DeadPlayerError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.Players[2].IsDead = false
	game.Players[2].IsLeft = true

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if isBetEnd {
		t.Error("it should be false")
	}
	if !errors.Is(err, gameerror.PlayerLeft) {
		t.Error("PlayerLeftError should be occurred")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = true
	game.Players[2].IsDead = false
	game.Players[2].IsLeft = false

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if nextPlayer != "lee" {
		t.Error("nextPlayer should be lee")
	}
	if isBetEnd {
		t.Error("it should be false")
	}
	if err != nil {
		t.Error("there should be no error")
	}
	if game.Players[2].IsDead != true {
		t.Error("han's IsPlayerDead field should be set to true")
	}

	// 여기부터는 정상 베팅
	betInfo.BetAmount = 20
	betInfo.PlayerName = "han"
	betInfo.IsDead = false
	game.Players[2].IsDead = false

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if isBetEnd {
		t.Error("it should be false")
	}
	if err != nil {
		t.Error("there should be no error")
	}

	betInfo.BetAmount = 20
	betInfo.PlayerName = "lee"
	betInfo.IsDead = false

	game.HandleBet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "choi"
	betInfo.IsDead = false

	game.HandleBet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "kim"
	betInfo.IsDead = false

	game.HandleBet(betInfo)

	betInfo.BetAmount = 20
	betInfo.PlayerName = "park"
	betInfo.IsDead = false

	nextPlayer, isBetEnd, err = game.HandleBet(betInfo)
	if !isBetEnd {
		t.Error("bet should be finished")
	}
	if err != nil {
		t.Error("there should be no error")
	}
}

func TestGame_isValidBet(t *testing.T) {
	var p player.Player

	p = player.New("kim", 300, 200)
	game := New()
	game.SetPlayers([]*player.Player{&p})

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
	if !errors.Is(err, gameerror.OverBalance) {
		t.Error("OverBalanceError should be occurred")
	}

	game.currentBet = 20
	p = player.New("kim", 300, 200)
	betType, err = game.isValidBet(&p, 10)
	if !errors.Is(err, gameerror.LowBetting) {
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

	game := New()
	game.SetPlayers(players)
	nextIdx, err := game.getNextPlayerIdx()
	if err != nil {
		t.Error("there should be no error")
	}

	if nextIdx != 4 {
		fmt.Println(nextIdx)
		t.Error("next index should be 4")
	}

	players = []*player.Player{
		{Nickname: "kim", IsReady: true, IsDead: true, IsLeft: false},
		{Nickname: "park", IsReady: true, IsDead: true, IsLeft: false},
		{Nickname: "han", IsReady: true, IsDead: true, IsLeft: false},
		{Nickname: "lee", IsReady: true, IsDead: false, IsLeft: true},
		{Nickname: "choi", IsReady: true, IsDead: true, IsLeft: false},
	}

	game = New()
	game.SetPlayers(players)
	nextIdx, err = game.getNextPlayerIdx()
	if !errors.Is(err, gameerror.NoPlayersLeft) {
		t.Error("there should be NoPlayersLeftError")
	}
}
