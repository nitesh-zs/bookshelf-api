// nolint:dupl //Cannot use same test for different migrations
package migrations

import (
	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/log"
)

type K20220406123607 struct {
}

func (k K20220406123607) Up(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method UP of migration K20220406123607: USER TABLE")

	query := `	drop type if exists "user_types";
				create type user_types as enum ('general', 'admin');
				create table if not exists "user"(
					id char(36) primary key,
					email varchar(255) unique not null,
					name varchar(255) not null,
					type user_types,
					c_time timestamp default now()
				);`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (k K20220406123607) Down(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method DOWN of migration K20220406123607: USER TABLE")

	query := `	drop table if exists "user";
				drop type if exists user_types;`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}

	return err
}
