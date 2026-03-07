package interfaces

type IDatabase interface {
	Create(value interface{}) IDatabase
	First(dest interface{}, conds ...interface{}) IDatabase
	Find(dest interface{}, conds ...interface{}) IDatabase
	Save(value interface{}) IDatabase
	Delete(value interface{}, conds ...interface{}) IDatabase
	Where(query interface{}, args ...interface{}) IDatabase
	Preload(query string, args ...interface{}) IDatabase
	Offset(offset int) IDatabase
	Limit(limit int) IDatabase
	Model(value interface{}) IDatabase
	Count(count *int64) IDatabase
	Joins(query string, args ...interface{}) IDatabase
	Group(name string) IDatabase
	Distinct(args ...interface{}) IDatabase
	GetError() error
}
