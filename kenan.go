package kenan

import (
	"encoding/base64"
	"fmt"
	"github.com/louismax/kenan/core"
	"gorm.io/gorm"
	"time"
)

var settings core.Settings

func Init(svrTag string, opts ...core.KenanOption) error {
	if svrTag == "" {
		panic("Tag is not allowed to be empty !! (应用标识不允许为空!!)")
	}
	printLogo()

	settings.ServerTag = svrTag
	for _, opt := range opts {
		if err := opt.Join(&settings); err != nil {
			return fmt.Errorf("kenan框架初始化异常:%v", err)
		}
	}
	fmt.Printf("***** 应用标识:[%+v]\n", settings.ServerTag)
	settings.ValidateSetting()

	ParsingLocalConfig()

	LogInit()

	mapDb = make(map[string]*gorm.DB)

	return nil
}

func printLogo() {
	decodeBytes, _ := base64.StdEncoding.DecodeString("CuKWiOKWiOKVlyAg4paI4paI4pWX4paI4paI4paI4paI4paI4paI4paI4pWX4paI4paI4paI4pWXICAg4paI4paI4pWXIOKWiOKWiOKWiOKWiOKWiOKVlyDilojilojilojilZcgICDilojilojilZcK4paI4paI4pWRIOKWiOKWiOKVlOKVneKWiOKWiOKVlOKVkOKVkOKVkOKVkOKVneKWiOKWiOKWiOKWiOKVlyAg4paI4paI4pWR4paI4paI4pWU4pWQ4pWQ4paI4paI4pWX4paI4paI4paI4paI4pWXICDilojilojilZEK4paI4paI4paI4paI4paI4pWU4pWdIOKWiOKWiOKWiOKWiOKWiOKVlyAg4paI4paI4pWU4paI4paI4pWXIOKWiOKWiOKVkeKWiOKWiOKWiOKWiOKWiOKWiOKWiOKVkeKWiOKWiOKVlOKWiOKWiOKVlyDilojilojilZEK4paI4paI4pWU4pWQ4paI4paI4pWXIOKWiOKWiOKVlOKVkOKVkOKVnSAg4paI4paI4pWR4pWa4paI4paI4pWX4paI4paI4pWR4paI4paI4pWU4pWQ4pWQ4paI4paI4pWR4paI4paI4pWR4pWa4paI4paI4pWX4paI4paI4pWRCuKWiOKWiOKVkSAg4paI4paI4pWX4paI4paI4paI4paI4paI4paI4paI4pWX4paI4paI4pWRIOKVmuKWiOKWiOKWiOKWiOKVkeKWiOKWiOKVkSAg4paI4paI4pWR4paI4paI4pWRIOKVmuKWiOKWiOKWiOKWiOKVkQrilZrilZDilZ0gIOKVmuKVkOKVneKVmuKVkOKVkOKVkOKVkOKVkOKVkOKVneKVmuKVkOKVnSAg4pWa4pWQ4pWQ4pWQ4pWd4pWa4pWQ4pWdICDilZrilZDilZ3ilZrilZDilZ0gIOKVmuKVkOKVkOKVkOKVnQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIA==")
	fmt.Println(string(decodeBytes))
	fmt.Printf("***** Kenan框架初始化中...............\n")
	time.Sleep(2 * time.Second)
}
