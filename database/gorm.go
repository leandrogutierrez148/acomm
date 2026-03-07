package database

import (
	"fmt"
	"log"

	"github.com/lgutierrez148/acomm/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func New(user, password, dbname, host, port string) (db interfaces.IDatabase, close func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %s", err)
	}

	return &GormDB{db: gormDB}, sqlDB.Close
}

func (g *GormDB) Create(value interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Create(value)}
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.First(dest, conds...)}
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Find(dest, conds...)}
}

func (g *GormDB) Save(value interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Save(value)}
}

func (g *GormDB) Delete(value interface{}, conds ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Delete(value, conds...)}
}

func (g *GormDB) Where(query interface{}, args ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Where(query, args...)}
}

func (g *GormDB) Preload(query string, args ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Preload(query, args...)}
}

func (g *GormDB) Offset(offset int) interfaces.IDatabase {
	return &GormDB{db: g.db.Offset(offset)}
}

func (g *GormDB) Limit(limit int) interfaces.IDatabase {
	return &GormDB{db: g.db.Limit(limit)}
}

func (g *GormDB) Model(value interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Model(value)}
}

func (g *GormDB) Count(count *int64) interfaces.IDatabase {
	return &GormDB{db: g.db.Count(count)}
}

func (g *GormDB) Joins(query string, args ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Joins(query, args...)}
}

func (g *GormDB) Group(name string) interfaces.IDatabase {
	return &GormDB{db: g.db.Group(name)}
}

func (g *GormDB) Distinct(args ...interface{}) interfaces.IDatabase {
	return &GormDB{db: g.db.Distinct(args...)}
}

func (g *GormDB) GetError() error {
	return g.db.Error
}
