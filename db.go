package main

import (
	"github.com/hashicorp/go-memdb"
	"log"
	"sync/atomic"
)

type counter int32

func (c *counter) next() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

var db *memdb.MemDB
var id counter = 0

const ArticleTable = "article"
const ArticleId = "id"

func init() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			ArticleTable: {
				Name: ArticleTable,
				Indexes: map[string]*memdb.IndexSchema{
					ArticleId: {
						Name:    ArticleId,
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}
	var err error
	db, err = memdb.NewMemDB(schema)
	if err != nil {
		log.Fatal(err)
	}
}

func NextId() int {
	return int(id.next())
}
