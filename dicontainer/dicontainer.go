package dicontainer

import (
	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/controller"
	"github.com/hackathon-20260110/api/driver"
	"github.com/hackathon-20260110/api/service"
	"go.uber.org/dig"
)

func GetContainer() *dig.Container {
	container := dig.New()
	err := container.Provide(driver.NewPsql)
	if err != nil {
		panic(err)
	}
	err = container.Provide(driver.NewFirestore)
	if err != nil {
		panic(err)
	}
	err = container.Provide(driver.NewGenAIClient)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewR2ClientFromEnv)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewR2Adapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewUserAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewLLMAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewOnboardingAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewUserInfoAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewAvatarAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewProfileAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewDiagnosisAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(service.NewDiagnosisService)
	if err != nil {
		panic(err)
	}
	err = container.Provide(controller.NewDiagnosisController)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewAvatarChatAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewMissionAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewMatchingAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewUserChatAdapter)
	if err != nil {
		panic(err)
	}
	err = container.Provide(adapter.NewNotificationAdapter)
	if err != nil {
		panic(err)
	}
	return container
}
