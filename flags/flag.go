package flags

import (
	"github.com/Brukhan/project/option"
)

type Options struct {
	HttpAddress string
	DbPath string
	Db string
}

func Register(place *Options) {
	set := option.NewOptionSet("CayleyOptions")
	set.String(&place.HttpAddress, "http-address", "HttpAddress", ":9000", "Address to listen")
	set.String(&place.DbPath, "dbpath", "DbPath", "cayley_test", "Path to the database")
	set.String(&place.Db, "db", "Db", "leveldb", "Database Backend")
}
