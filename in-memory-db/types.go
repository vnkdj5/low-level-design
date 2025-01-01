package inmemorydb

import "sync"

type Record map[string]interface{}

type Table struct {
	name      string
	schema    map[string]string                   // Column name -> Data type
	data      map[string]Record                   // Row ID -> Record (row data)
	indexes   map[string]map[interface{}][]string // Column -> Value -> List of Row IDs
	dataLock  sync.RWMutex
	indexLock sync.RWMutex // Lock for index operations
	persisted bool         // Flag for persistence support
}

type Condition struct {
	Attribute   string
	Operator    string
	Value       interface{}
	SecondValue interface{} // Used for BETWEEN
}
