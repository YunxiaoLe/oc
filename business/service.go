package business

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/dao"
	"github.com/gin-gonic/gin"
)

type Service struct {
	dao *dao.Dao
	Ctx *gin.Context
}

func NewService(ctx *gin.Context) *Service {
	return &Service{dao: dao.GetDaoInstance(), Ctx: ctx}
}
