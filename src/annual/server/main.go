package server

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type lotteryController struct {
	Ctx iris.Context
}

var usrList []string
var mu sync.Mutex

func (*lotteryController) Get() string {
	count := len(usrList)
	return fmt.Sprintf("当前参与的抽奖人数：%d", count)
}

// 在导入用户时会出现并发
func (c *lotteryController) PostImport() string {
	strUsrs := c.Ctx.FormValue("users")
	usrs := strings.Split(strUsrs, ",")

	mu.Lock()
	defer mu.Unlock()

	for _, u := range usrs {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			usrList = append(usrList, u)
		}
	}
	return fmt.Sprintf("导入抽奖人数：%d, 当前参与的抽奖人数：%d", len(usrs), len(usrList))
}

func (c *lotteryController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()

	count := len(usrList)
	if count <= 0 {
		return fmt.Sprintf("当前没有参与的用户，请先导入用户")
	} else if count == 1 {
		return fmt.Sprintf("恭喜用户%s中奖", usrList[0])
	} else {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		return fmt.Sprintf("恭喜用户%s中奖", usrList[index])
	}
}

func Run() {
	app := NewApp()
	app.Run(iris.Addr(":8080"))
}

func NewApp() (app *iris.Application) {
	app = iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}
