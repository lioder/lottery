package main

import (
	"fmt"
	"github.com/kataras/iris/httptest"
	"lottery/src/main/server"
	"sync"
	"testing"
)

func TestMVC(t *testing.T) {
	e := httptest.New(t, server.NewApp())
	e.GET("/").Expect().Status(httptest.StatusOK).Body().
		Equal("当前参与的抽奖人数：0")

	var wg sync.WaitGroup
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users", fmt.Sprintf("test_u%d", i)).Expect().Status(httptest.StatusOK)
		}(i)
	}
	wg.Wait()
	e.GET("/").Expect().Status(httptest.StatusOK).Body().
		Equal("当前参与的抽奖人数：50")
}
