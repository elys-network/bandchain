package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type BandDB struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewDB(dialect, path string, metadata map[string]string) (*BandDB, error) {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&Metadata{},
		&Event{},
		&Validator{},
		&ValidatorVote{},
	)

	db.Model(&ValidatorVote{}).AddForeignKey(
		"consensus_address",
		"validators(consensus_address)",
		"RESTRICT",
		"RESTRICT",
	)

	for key, value := range metadata {
		err := db.Where(Metadata{Key: key}).
			Assign(Metadata{Value: value}).
			FirstOrCreate(&Metadata{}).Error
		if err != nil {
			panic(err)
		}
	}

	return &BandDB{db: db}, nil
}

func (b *BandDB) BeginTransaction() {
	if b.tx != nil {
		panic("BeginTransaction: Cannot begin a new transaction without closing the pending one.")
	}
	b.tx = b.db.Begin()
	if b.tx.Error != nil {
		panic(b.tx.Error)
	}
}

func (b *BandDB) Commit() {
	err := b.tx.Commit().Error
	if err != nil {
		panic(err)
	}
	b.tx = nil
}

func (b *BandDB) RollBack() {
	err := b.tx.Rollback()
	if err != nil {
		panic(err)
	}
	b.tx = nil
}

func (b *BandDB) HandleTransaction(tx auth.StdTx, txHash []byte, logs sdk.ABCIMessageLogs) {
	msgs := tx.GetMsgs()
	if len(msgs) != len(logs) {
		panic("Inconsistent size of msgs and logs.")
	}

	for idx, msg := range msgs {
		events := logs[idx].Events
		kvMap := make(map[string]string)
		for _, event := range events {
			for _, kv := range event.Attributes {
				kvMap[event.Type+"."+kv.Key] = kv.Value
			}
		}
		b.HandleMessage(msg, kvMap)
	}
}

func (b *BandDB) HandleMessage(msg sdk.Msg, attributes map[string]string) error {
	switch msg := msg.(type) {
	// Just proof of concept
	case zoracle.MsgCreateDataSource:
		{
			_ = msg
			// Event message split events on report event eg.
			// message map[action:report]
			// message map[sender:band17xpfvakm2amg962yls6f84z3kell8c5lfkrzn4]
			// action, ok := attributes["action"]
			// if ok {
			// 	b.handleMessageEvent(action)
			// }
			return nil
		}
	default:
		// TODO: Better logging
		return errors.New("HandleMessage: There isn't event handler for this type")
	}
}