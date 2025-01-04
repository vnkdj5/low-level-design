package inmemorydb

import (
	"fmt"
	"sync"
)

type InMemoryDB struct {
	tables map[string]*Table //Table Name --> table Instance
	dbLock sync.RWMutex
}

func NewInMemoryDB() Database {
	return &InMemoryDB{
		tables: make(map[string]*Table),
	}
}

func (db *InMemoryDB) CreateTable(name string, schema map[string]string) error {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	if _, exists := db.tables[name]; exists {
		return fmt.Errorf("Table %s already exists", name)
	}
	db.tables[name] = &Table{
		name:    name,
		schema:  schema,
		data:    make(map[string]Record),
		indexes: make(map[string]map[interface{}][]string),
	}
	return nil
}

func (db *InMemoryDB) Insert(tableName string, key string, record Record) error {
	db.dbLock.RLock()
	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	table.dataLock.Lock()
	defer table.dataLock.Unlock()

	//Validate schema
	for column, dataType := range table.schema {
		if value, ok := record[column]; ok {
			if fmt.Sprintf("%T", value) != dataType {
				return fmt.Errorf("invalid data type for column %s, expected %s", column, dataType)
			}
		} else {
			//Ignore this if we want optional fields
			return fmt.Errorf("missing value for column %s", column)
		}
	}

	//Insert Data
	table.data[key] = record

	//Update indexes

	table.indexLock.Lock()

	for column, value := range record {
		if index, exists := table.indexes[column]; exists {
			// If index exists for this column, update it
			index[value] = append(index[value], key)
		}
	}
	table.indexLock.Unlock()
	return nil

}

func (db *InMemoryDB) Select(tableName, attribute, whereKey string, whereValue interface{}) ([]interface{}, error) {
	db.dbLock.RLock()
	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	table.dataLock.RLock()
	defer table.dataLock.RUnlock()

	result := []interface{}{}
	if index, ok := table.indexes[whereKey]; ok { // Use index if available
		if recordIDs, exists := index[whereValue]; exists {
			for _, id := range recordIDs {
				result = append(result, table.data[id][attribute])
			}
		}
	} else { // Fallback: scan all records
		for _, record := range table.data {
			if record[whereKey] == whereValue {
				result = append(result, record[attribute])
			}
		}
	}

	return result, nil
}

func (db *InMemoryDB) Get(tableName string, id string) (Record, error) {
	db.dbLock.RLock()

	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}
	table.dataLock.RLock()
	record, found := table.data[id]
	table.dataLock.RUnlock()
	if !found {
		return nil, fmt.Errorf("record with ID %s not found", id)
	}
	return record, nil
}

func (db *InMemoryDB) Delete(tableName string, id string) error {
	db.dbLock.RLock()

	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	table.dataLock.Lock()
	delete(table.data, id)
	table.dataLock.Unlock()

	return nil
}

func (db *InMemoryDB) CreateIndex(tableName, column string) error {
	db.dbLock.RLock()
	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	table.indexLock.RLock()
	_, exists = table.indexes[column]
	table.indexLock.RUnlock()
	if exists {
		return fmt.Errorf("index on column %s already exists", column)
	}
	table.indexLock.Lock()
	defer table.indexLock.Unlock()
	table.indexes[column] = make(map[interface{}][]string)
	for id, record := range table.data {
		value := record[column]
		table.indexes[column][value] = append(table.indexes[column][value], id)
	}

	return nil
}

// func (db *InMemoryDB) saveToDisk() error {
// 	db.dbLock.RLock()
// 	defer db.dbLock.RUnlock()

// 	// Serialize and save to persistencePath
// 	// Use JSON or another format
// 	return nil
// }

func (db *InMemoryDB) SelectWithConditions(
	tableName string,
	selectAttributes []string,
	conditions []Condition,
	logicalOperator string,
) ([]map[string]interface{}, error) {
	db.dbLock.RLock()
	table, exists := db.tables[tableName]
	db.dbLock.RUnlock()
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	table.dataLock.RLock()
	defer table.dataLock.RUnlock()

	var result []map[string]interface{}

	// Iterate over all records in the table
	for _, record := range table.data {
		if evaluateConditions(record, conditions, logicalOperator) {
			// Prepare the selected attributes for the result
			selectedRecord := make(map[string]interface{})
			for _, attr := range selectAttributes {
				if value, exists := record[attr]; exists {
					selectedRecord[attr] = value
				}
			}
			result = append(result, selectedRecord)
		}
	}

	return result, nil
}

func evaluateConditions(record Record, conditions []Condition, logicalOperator string) bool {
	results := make([]bool, len(conditions))

	for i, condition := range conditions {
		value, exists := record[condition.Attribute]
		if !exists {
			results[i] = false
			continue
		}

		switch condition.Operator {
		case "=":
			results[i] = value == condition.Value
		case "!=":
			results[i] = value != condition.Value
		case "<":
			results[i] = compareNumeric(value, condition.Value, func(a, b float64) bool { return a < b })
		case ">":
			results[i] = compareNumeric(value, condition.Value, func(a, b float64) bool { return a > b })
		case "<=":
			results[i] = compareNumeric(value, condition.Value, func(a, b float64) bool { return a <= b })
		case ">=":
			results[i] = compareNumeric(value, condition.Value, func(a, b float64) bool { return a >= b })
		case "BETWEEN":
			results[i] = compareNumeric(value, condition.Value, func(a, b float64) bool { return a >= b }) &&
				compareNumeric(value, condition.SecondValue, func(a, b float64) bool { return a <= b })
		default:
			results[i] = false
		}
	}

	if logicalOperator == "AND" {
		for _, res := range results {
			if !res {
				return false
			}
		}
		return true
	} else if logicalOperator == "OR" {
		for _, res := range results {
			if res {
				return true
			}
		}
		return false
	}

	return false
}
