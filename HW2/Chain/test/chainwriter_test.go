package test

import (
	"Chain/pkg/block"
	"Chain/pkg/blockchain/chainwriter"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func TestStoreOrphanBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	bl := MockedBlock()
	ub := &chainwriter.UndoBlock{}
	br := cw.StoreBlock(bl, ub, 0)
	if br.BlockFile != "data/block_0.txt" {
		t.Errorf("Expected file name: %v Actual file name: %v", "data/block_0.txt", br.BlockFile)
	}
	if br.UndoFile != "" {
		t.Errorf("Expected file name: %v Actual file name: %v", "", br.UndoFile)
	}
}

func TestStoreBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	bl := MockedBlock()
	ub := MockedUndoBlock()
	br := cw.StoreBlock(bl, ub, 0)
	if br.BlockFile != "data/block_0.txt" {
		t.Errorf("Expected file name: %v Actual file name: %v", "data/block_0", br.BlockFile)
	}
	if br.UndoFile != "data/undo_0.txt" {
		t.Errorf("Expected file name: %v Actual file name: %v", "", br.UndoFile)
	}
}

func TestWriteBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	b := MockedBlock()
	pb := block.EncodeBlock(b)
	serializedBlock, _ := proto.Marshal(pb)
	fi := cw.WriteBlock(serializedBlock)
	if fi.StartOffset != 0 {
		t.Errorf("Expected start offset: %v\nActual start offset: %v", 0, fi.StartOffset)
	}
	if int(fi.EndOffset) != len(serializedBlock) {
		t.Errorf("Expected end offset: %v\nActual end offset: %v", 0, fi.EndOffset)
	}
	if fi.FileName != "data/block_0.txt" {
		t.Errorf("Expected file name: %v Actual file name: %v", "data/block_0", fi.FileName)
	}
}

func TestWriteUndoBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	ub := MockedUndoBlock()
	pub := chainwriter.EncodeUndoBlock(ub)
	serializedUndoBlock, _ := proto.Marshal(pub)
	ufi := cw.WriteUndoBlock(serializedUndoBlock)
	if ufi.StartOffset != 0 {
		t.Errorf("Expected start offset: %v\nActual start offset: %v", 0, ufi.StartOffset)
	}
	if int(ufi.EndOffset) != len(serializedUndoBlock) {
		t.Errorf("Expected end offset: %v\nActual end offset: %v", 0, ufi.EndOffset)
	}
	if ufi.FileName != "data/undo_0.txt" {
		t.Errorf("Expected file name: %v Actual file name: %v", "data/block_0", ufi.FileName)
	}
}

func TestReadBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	b := MockedBlock()
	pb := block.EncodeBlock(b)
	serializedBlock, _ := proto.Marshal(pb)
	fi := cw.WriteBlock(serializedBlock)
	b2 := cw.ReadBlock(fi)
	if !reflect.DeepEqual(b, b2) {
		t.Errorf("Expected block: %v\nActual block: %v", b, b2)
	}
}

func TestReadUndoBlock(t *testing.T) {
	defer cleanUp()
	cw := chainwriter.New(chainwriter.DefaultConfig())
	ub := MockedUndoBlock()
	pub := chainwriter.EncodeUndoBlock(ub)
	serializedUndoBlock, _ := proto.Marshal(pub)
	ufi := cw.WriteUndoBlock(serializedUndoBlock)
	ub2 := cw.ReadUndoBlock(ufi)
	if !reflect.DeepEqual(ub, ub2) {
		t.Errorf("Expected block: %v\nActual block: %v", ub, ub2)
	}
}

func TestRead100Blocks(t *testing.T) {
	defer cleanUp()
	config := chainwriter.DefaultConfig()
	config.MaxBlockFileSize = 100
	cw := chainwriter.New(config)

	var blocks []*block.Block
	var fileInfos []*chainwriter.FileInfo

	// write blocks
	for i := 0; i < 100; i++ {
		b := MockedBlock()
		b.Header.Nonce = uint32(i)
		blocks = append(blocks, b)
		pb := block.EncodeBlock(b)
		serializedBlock, _ := proto.Marshal(pb)
		fi := cw.WriteBlock(serializedBlock)
		fileInfos = append(fileInfos, fi)
	}

	// read blocks
	for i := 0; i < 100; i++ {
		b := cw.ReadBlock(fileInfos[i])
		if !reflect.DeepEqual(blocks[i], b) {
			t.Errorf("Block: %v/99\nExpected block: %v\nActual block: %v", i, blocks[i], b)
		}
	}

}
