package test

import (
	"Chain/pkg/blockchain/blockinfodatabase"
	"reflect"
	"testing"
)

func TestStoreBlockRecord(t *testing.T) {
	defer cleanUp()
	blockinfo := blockinfodatabase.New(blockinfodatabase.DefaultConfig())
	br := MockedBlockRecord()
	blockinfo.StoreBlockRecord("hash", br)
}

func TestGetSameRecord(t *testing.T) {
	defer cleanUp()
	blockinfo := blockinfodatabase.New(blockinfodatabase.DefaultConfig())
	br := MockedBlockRecord()
	blockinfo.StoreBlockRecord("hash", br)
	br2 := blockinfo.GetBlockRecord("hash")
	if !reflect.DeepEqual(br, br2) {
		t.Errorf("Block records not equal")
	}
}

func TestGetDifferentRecords(t *testing.T) {
	defer cleanUp()
	blockinfo := blockinfodatabase.New(blockinfodatabase.DefaultConfig())
	br := MockedBlockRecord()
	br2 := MockedBlockRecord()
	br2.UndoEndOffset = 20
	blockinfo.StoreBlockRecord("hash", br)
	blockinfo.StoreBlockRecord("hash2", br2)
	rbr := blockinfo.GetBlockRecord("hash")
	rbr2 := blockinfo.GetBlockRecord("hash2")
	if reflect.DeepEqual(rbr, rbr2) {
		t.Errorf("Block records should not be equal")
	}
}
