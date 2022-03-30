// This is auto-generated file using 'krogo migrate' tool. DO NOT EDIT.
package migrations

import (
	dbmigration "github.com/krogertechnology/krogo/cmd/krogo/migration/dbMigration"
)

func All() map[string]dbmigration.Migrator {
	return map[string]dbmigration.Migrator{
		"20220328222746": K20220328222746{},
	}
}
