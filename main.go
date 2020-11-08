package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pthum/stripcontrol-golang/config"
	"github.com/pthum/stripcontrol-golang/database"
	"github.com/pthum/stripcontrol-golang/mappings"
)

func main() {
	config.InitConfig()
	mappings.CreateURLMappings()
	database.ConnectDataBase()
	// set release mode if set
	if config.CONFIG.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Listen and serve on 0.0.0.0:8080
	serve := fmt.Sprintf("%s:%s", config.CONFIG.Server.Host, config.CONFIG.Server.Port)
	mappings.Router.Run(serve)
}

// func maind() {
// 	notSet := `{}`
// 	setNull := `{"misoPin": null}`
// 	setValid := `{"misoPin": 123}`
// 	setValidString := `{"misoPin": "123"}`
// 	setValidStringZero := `{"misoPin": 0}`

// 	parseAndPrint(notSet)
// 	parseAndPrint(setNull)
// 	parseAndPrint(setValid)
// 	parseAndPrint(setValidString)
// 	parseAndPrint(setValidStringZero)
// }
// func parseAndPrint(str string) {
// 	var b models.LedStrip

// 	json.Unmarshal([]byte(str), &b)
// 	fmt.Printf("<Value:%d> <Set:%t> <Valid:%t>\n", b.MisoPin.Value, b.MisoPin.Set, b.MisoPin.Valid)
// }
