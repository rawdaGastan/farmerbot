// Package manager provides how to manage nodes, farms and power
package manager

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var testFarm = models.Farm{
	ID:          1,
	Description: "test",
	PublicIPs:   1,
}

func TestFarmManager(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mocks.NewMockRedisManager(ctrl)

	farmManager := NewFarmManager(db, log.Logger)

	t.Run("test valid define farm", func(t *testing.T) {
		db.EXPECT().SetFarm(testFarm).Return(nil)

		err := farmManager.Define(testFarm)
		assert.NoError(t, err)
	})

	t.Run("test invalid define farm: db failed", func(t *testing.T) {
		db.EXPECT().SetFarm(testFarm).Return(fmt.Errorf("error"))

		err := farmManager.Define(testFarm)
		assert.Error(t, err)
	})
}
