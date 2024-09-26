package db

import (
	"os"
	"time"

	"bms.dse/src/utils"
	"bms.dse/src/utils/gatekeeper"
	"bms.dse/src/utils/logutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Log struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp time.Time `gorm:"type:timestamptz"         json:"timestamp"`
	User      string    `gorm:"type:text"                json:"token"`
	Level     string    `gorm:"type:text"                json:"level"`
	Message   string    `gorm:"type:text"                json:"message"`
}

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("DB")

var db  *gorm.DB
var gkReady     = gatekeeper.NewGateKeeper(true)
var gkConnected = gatekeeper.NewGateKeeper(true)
// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func Create(v interface{}) *gorm.DB {
	gkReady.Wait()
	return db.Create(v)
}

func Table(name string) *gorm.DB {
	gkReady.Wait()
	return db.Table(name)
}

func AutoMigrate(models ...interface{}) error {
	gkReady.Wait()
	return db.AutoMigrate(models...)
}

// ------------------------------------------------------------
// : Debug
// ------------------------------------------------------------
func Debug() {
	if !utils.IsDebugMode() { return }
}

// ------------------------------------------------------------
// : Postgres
// ------------------------------------------------------------
func Init() {
	var err error

	envDSN := os.Getenv("DB_DSN")
	
	logger.Info().
		Str("dsn", envDSN).
		Msg("Connecting")
	 
	db, err = gorm.Open(postgres.Open(envDSN), &gorm.Config{})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Postgres")
		return
	}
	db.Logger = db.Logger.LogMode(0)

	gkReady.Unlock()

	logger.Info().Msg("Connected")
	go Debug()
}
