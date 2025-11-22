package auth

import (
	"os"
	"testing"
)

func TestCreateDb(t *testing.T) {
	if _, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING"); !ok {
		t.Skip()
	} else {
		_, err := CreateNewDb()
		if err != nil {
			t.Errorf("Not expecting an error when creating a new database instance, got %s", err.Error())
		}
	}
}
