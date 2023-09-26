package main

import (
	"fmt"
	"vkpdeveloper/statuspages-cli/ui"
	"vkpdeveloper/statuspages-cli/utils"
)

func main() {

	config, err := utils.InitAppConfig()

	if err != nil {
		fmt.Println(err)
	}

	config.ReadConfig()

	client := utils.NewStatusPageClient(config)

	ui.Run(config, client)
}
