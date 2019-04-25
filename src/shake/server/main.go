package server

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type shakeController struct {
	Ctx iris.Context
}

type gift struct {
	id    int
	name  string
	total int
	left  int
}

var gifts *gift = &gift{
	id:    1,
	name:  "现金",
	total: 1000,
	left:  1000,
}

var logger *log.Logger
var mu sync.Mutex

func Run() {
	app := NewApp()
	app.Run(iris.Addr(":8080"))
}

func NewApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&shakeController{})
	initLog()
	return app
}

func (c *shakeController) Get() string {
	return fmt.Sprintf("当前奖品总额：%d, 还剩%d", gifts.total, gifts.left)
}

func (c *shakeController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()
	code := luckyCode()
	if code <= 3 {
		ok, res := sendGift()
		if ok {
			saveLuckyData(code, gifts)
		}
		return res
	} else {
		return "对不起，没有中奖"
	}
}

func saveLuckyData(code int32, gift *gift) {
	logger.Printf("中奖号码：%d，还剩%d份", code, gift.left)
}

func initLog() {
	f, err := os.Create("./src/shake/shake.log")
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		logger = log.New(f, "", log.Lmicroseconds|log.Ldate)
	}
}

func luckyCode() int32 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10)
}

func sendGift() (bool, string) {
	if gifts.left > 0 {
		gifts.left--
		return true, "恭喜中奖"
	} else {
		return false, "奖品已发完"
	}
}
