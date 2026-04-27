package db
import ("sync"; "github.com/dong-labs/think/internal/core/db")
const NAME = "dida"
type DidaDatabase struct { *db.Database }
func NewDidaDatabase() *DidaDatabase { return &DidaDatabase{Database: db.NewDatabase(NAME)} }
var dbInstance *DidaDatabase
var once sync.Once
func GetDB() *DidaDatabase { once.Do(func() { dbInstance = NewDidaDatabase() }); return dbInstance }
