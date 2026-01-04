package migrations

import "github.com/go-gormigrate/gormigrate/v2"

func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		Migration_1765528367(),
		Migration_202512291126(),
	}
}
