package main

import (
	"testing"

	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/krogo"
)

func TestRefreshTables(t *testing.T) {
	k := krogo.New()
	seeder := datastore.NewSeeder(&k.DataStore, "./configs")
	seeder.RefreshTables(t, "book")
}
