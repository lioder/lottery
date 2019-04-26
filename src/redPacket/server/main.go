package server

// 压力测试: wrk -t10 -c10 -d5 http://localhost:8080/send?uid=1&money=1000&num=10
// wrk -t100 -c100 -d1 http://localhost:8080/get?id=1000001825
import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"sync"
	"time"
)

type packetController struct {
	Ctx iris.Context
}

var packetList = new(sync.Map)

// http://localhost:8080
// 返回所有的红包个数和金额
func (c *packetController) Get() map[uint32][2]int {
	rs := make(map[uint32][2]int)
	packetList.Range(func(key, value interface{}) bool {

		id := key.(uint32)
		list := value.([]uint)
		var total int
		for _, v := range list {
			total += int(v)
		}
		rs[id] = [2]int{len(list), total}
		return true
	})
	return rs
}

// http://localhost:8080/send?uid=1&money=100&num=100
// 发送红包
func (c *packetController) GetSend() string {
	uid, errUid := c.Ctx.URLParamInt("uid")
	money, errMoney := c.Ctx.URLParamFloat64("money")
	num, errNum := c.Ctx.URLParamInt("num")

	if errUid != nil || errMoney != nil || errNum != nil {
		return fmt.Sprintf("参数格式错误, uid: %d, money: %f, num: %d", uid, money, num)
	}

	moneyInCent := int(money * 100)
	if num < 1 || moneyInCent < 1 || uid < 1 {
		return fmt.Sprintf("参数数值错误, uid: %d, money: %f, num: %d", uid, money, num)
	}

	leftNum := num
	leftMoneyInCent := moneyInCent
	list := make([]uint, num)

	for leftNum > 0 {
		if leftNum == 1 {
			list[num-1] = uint(leftMoneyInCent)
			break
		}

		if leftMoneyInCent == leftNum {
			for ; leftNum > 0; leftNum-- {
				list[num-leftNum] = 1
			}
			break
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		rMax := 0.2
		if num >= 100 {
			rMax = 0.5
		} else if num >= 50 {
			rMax = 0.4
		} else if num >= 10 {
			rMax = 0.3
		}
		money := r.Intn(int(rMax * float64(leftMoneyInCent-leftNum)))
		if money < 1 {
			money = 1
		}
		list[num-leftNum] = uint(money)
		leftMoneyInCent -= int(list[num-leftNum])
		leftNum--
	}
	id := rand.Uint32()
	packetList.Store(id, list)
	return fmt.Sprintf("发送成功，分享链接: http://localhost:8080/get?id=%d", id)
}

func (c *packetController) GetGet() string {
	id, errId := c.Ctx.URLParamInt("id")
	if errId != nil {
		return fmt.Sprintf("参数错误，id: %d", id)
	}
	_, ok := packetList.Load(uint32(id))
	if !ok {
		return fmt.Sprintf("红包不存在，id: %d", id)
	}

	callback := make(chan uint)
	task := task{id: uint32(id), callback: callback}
	chTasks := tasks[id%taskNum]
	chTasks <- task
	money := <-callback
	if money <= 0 {
		return fmt.Sprintf("很遗憾没有抢到红包")
	} else {
		return fmt.Sprintf("恭喜抢到一个红包，金额为 %.2f", float64(money)*0.01)
	}
}

func newApp() *iris.Application {
	app := iris.Default()
	mvc.New(app.Party("/")).Handle(&packetController{})
	for i := 0; i < len(tasks); i++ {
		tasks[i] = make(chan task)
		go fetchPacketService(tasks[i])
	}
	return app
}

func Run() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
