package blockinfodatabase

import (
	"Chain/pkg/pro"
	"Chain/pkg/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

// BlockInfoDatabase is a wrapper for a levelDB
type BlockInfoDatabase struct {
	db *leveldb.DB
}

// New returns a BlockInfoDatabase given a Config
func New(config *Config) *BlockInfoDatabase {
	db, err := leveldb.OpenFile(config.DatabasePath, nil)
	if err != nil {
		utils.Debug.Printf("Unable to initialize BlockInfoDatabase with path {%v}", config.DatabasePath)
	}
	return &BlockInfoDatabase{db: db}
}

// TASK_1
// StoreBlockRecord stores a BlockRecord in the BlockInfoDatabase.
// hash is the hash of the block, and the key for the blockRecord.
// blockRecord is the value we're storing in the database
//
// At a high level, here's what this function is doing:
// (1) converting a blockRecord to a protobuf version (for more effective storage)
// (2) converting the protobuf to bytes
// (3) storing the byte version of the blockRecord in our database
func (blockInfoDB *BlockInfoDatabase) StoreBlockRecord(hash string, blockRecord *BlockRecord) {
	// 1. convert the blockRecord to a proto version
	
	// 2. marshaling (serializing) the proto record to bytes
	
	// 3. checking that the marshalling process didn't throw an error
	
	// 4. attempting to store the bytes in our database AND checking to make
	// sure that the storing process doesn't fail. The Put(key, value, writeOptions)
	// function is levelDB's.
	
}

// TASK_2
// GetBlockRecord returns a BlockRecord from the BlockInfoDatabase given
// the relevant block's hash.
// hash is the hash of the block, and the key for the blockRecord.
//
// At a high level, here's what this function is doing:
// (1) retrieving the byte version of the protobuf record.
// (2) converting the bytes to protobuf
// (3) converting the protobuf to blockRecord and returning that.
func (blockInfoDB *BlockInfoDatabase) GetBlockRecord(hash string) *BlockRecord {
	// 1. attempting to retrieve the byte-version of the protobuf record
	// from our database AND checking that the value is retrieved successfully.
	// The Get(key, writeOptions) function is levelDB's.
	
	// 2. creating a protobuf blockRecord object to fill
	
	// 3. unmarshalling (deserializing) the bytes stored in the database into the
	// protobuf object created on line 66. Checking that the conversion process
	// from bytes to protobuf object succeeds.
	
	// 4. convert the protobuf record to a normal blockRecord and returning that.
	
}
