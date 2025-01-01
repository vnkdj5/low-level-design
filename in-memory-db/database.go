package inmemorydb

type Database interface {
	CreateTable(name string, schema map[string]string) error
	Insert(tableName string, key string, record Record) error
	Select(tableName, attribute, whereKey string, whereValue interface{}) ([]interface{}, error)
	SelectWithConditions(
		tableName string,
		selectAttributes []string,
		conditions []Condition,
		logicalOperator string,
	) ([]map[string]interface{}, error)
	Get(tableName string, key string) (Record, error)
	//Update(tableName string, key string, updates Record) error
	CreateIndex(tableName, column string) error
	Delete(tableName string, key string) error
}
