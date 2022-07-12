package database

type DatabaseConnector interface {
	InsertLogItem(collection string, item map[string]interface{})
}
