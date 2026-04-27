package db
import ("sync"; "github.com/dong-labs/think/internal/core/db")
const NAME = "expire"
type ExpireDatabase struct { *db.Database }
func NewExpireDatabase() *ExpireDatabase { return &ExpireDatabase{Database: db.NewDatabase(NAME)} }
var dbInstance *ExpireDatabase
var once sync.Once
func GetDB() *ExpireDatabase { once.Do(func() { dbInstance = NewExpireDatabase() }); return dbInstance }
