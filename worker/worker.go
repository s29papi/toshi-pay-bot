package worker

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lib/pq"
	"github.com/s29papi/wag3r-bot/worker/client"
	"github.com/s29papi/wag3r-bot/worker/db"
	"github.com/s29papi/wag3r-bot/worker/env"
)

type DepositRequestData struct {
	Fid     int64
	Address string
	TxHash  string
	Amount  *big.Int
}

type Worker struct {
	ctx                         context.Context
	T                           *time.Ticker
	db                          *db.DB // our own client
	stopped                     bool
	StoppedFrameEvents          bool
	StoppedDepositAndWithdrawal bool
	lastProcReqTime             *int64
	s                           *client.HTTPService
	Req                         chan struct{}
	pauseTickerFn               chan struct{}
	depositRequests             []DepositRequestData

	ethclient *ethclient.Client

	done chan struct{}
}

func NewWorker() *Worker {
	ethclient, err := ethclient.Dial("https://base-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Printf("New Ethclient Error: %v\n", err)
		return nil
	}
	ctx := context.Background()
	// db for bot
	psqlInfo, err := pq.ParseURL(env.RENDER_POSTGRES_URL)
	if err != nil {
		log.Fatalln(err)
	}
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	db := db.NewDB(sdb)

	// time duration for bot
	val, err := strconv.Atoi(env.DURATION_STR)
	if err != nil {
		log.Fatalln("Error: conversion of DURATION_STR to int")
	}
	dur := time.Duration(val) * time.Second
	t := time.NewTicker(dur)
	s := client.NewHTTPService()

	return &Worker{
		ctx:             ctx,
		T:               t,
		db:              db,
		s:               s,
		Req:             make(chan struct{}),
		depositRequests: make([]DepositRequestData, 0),
		lastProcReqTime: new(int64),
		pauseTickerFn:   make(chan struct{}),
		ethclient:       ethclient,
		done:            make(chan struct{}),
	}
}

func (w *Worker) Start() {
	// go w.tick(w.T.C)
	go w.workloop()
	<-w.done

	log.Println("Bot Stopping")
	w.stopped = true

}

// // update tick to start after processing is don
// func (w *Worker) tick(t <-chan time.Time) {
// 	for {
// 		<-t
// 		log.Println("Tick go-routine:, new request initiated")
// 		w.Req <- struct{}{}
// 		log.Println("Tick go-routine: paused.")
// 		<-w.pauseTickerFn
// 	}
// }

// starts a process, then completes it by starting a tick
// stopping workloop mean stop processing requests
func (w *Worker) workloop() {
	txs := w.buildUserMentionToTx()
	fmt.Println(txs[0].CastHash)
	// for {
	// 	if w.stopped {
	// 		break
	// 	}

	// 	if !w.StoppedFrameEvents {
	// 		fmt.Println("Hello world")
	// 	}

	// 	if !w.StoppedDepositAndWithdrawal {
	// 		fmt.Println(3944884)
	// 	}

	// 	fmt.Println("Sleeping for 2 seconds...")
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("Awake after 2 seconds!")
	// }
}

func fetchEthDepositTx() {

}

func (w *Worker) Stop() {
	if w.stopped {
		log.Println("Can't stop Bot, Bot was just stopped")
		return
	}
	w.T.Stop()
	w.db.Close()
	log.Println("DB closed.")

	log.Println("Tick go-routine: stopped.")
	log.Println("Exiting...")
	w.done <- struct{}{}
}

// m := filterUserMentionsByLastUpdate(fetchUserMentions(w.s), 1708841908)
// fmt.Println(len(m.Notifications))
// change this values, they should be gotten from db
// var lastEthDepositBlock *big.Int
// var lastEthDepositTime int64
// fetchEthDepositsFromLastUpdate(w.ethclient, w.ctx, lastEthDepositBlock, lastEthDepositTime)
// fmt.Println(w.db.GetLastProcReqTime())
// create a table called last user metionm
// fetchUserMentions(w.s)
// for {

// }

// <-w.Req
// log.Println("Initiating a new request")
// mentions := GetMentions(w.s)
// w.process(mentions)
// w.pauseTickerFn <- struct{}{}
// log.Println("Tick go-routine: un-paused.")

func (w *Worker) SendDepositRequest(d DepositRequestData) {
	w.depositRequests = append(w.depositRequests, d)
	log.Println("SendDepositRequest: New Pending Deposit")
}
