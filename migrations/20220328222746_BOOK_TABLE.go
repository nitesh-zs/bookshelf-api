package migrations

import (
	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/log"
)

type K20220328222746 struct {
}

func (k K20220328222746) Up(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method UP of migration K20220328222746: BOOK TABLE")

	query := `create table if not exists book (
				id char(36) primary key,
				title varchar(255) not null,
				author varchar(255) not null,
				summary text not null,
				genre varchar(255) not null,
				year integer,
				reg_num varchar(50) not null unique,
				publisher varchar(255),
				language varchar(50) not null,
				image_uri varchar(255)
			);`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (k K20220328222746) Down(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method DOWN of migration K20220328222746: BOOK TABLE")

	query := `DROP TABLE IF EXISTS book;`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}

	return err
}
