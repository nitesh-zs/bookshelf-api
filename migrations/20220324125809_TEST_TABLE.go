package migrations

import (
	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/log"
)

type K20220324125809 struct {
}

func (k K20220324125809) Up(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method UP of migration K20220324125809: TEST TABLE")
	query := `CREATE TABLE IF NOT EXISTS test(
				id SERIAL PRIMARY KEY);`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (k K20220324125809) Down(d *datastore.DataStore, logger log.Logger) error {
	logger.Info("Running method DOWN of migration K20220324125809: TEST TABLE")
	query := `DROP TABLE IF EXISTS test;`

	_, err := d.DB().Exec(query)
	if err != nil {
		logger.Error(err)
	}
	return err
}
