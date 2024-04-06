package main

import (
	"github.com/Swechhya/isol8r-backend/router"
	"github.com/Swechhya/isol8r-backend/setup"
)

func main() {
	// TestRun()
	setup.Setup()

	//setup router
	r := router.SetupRouter()
	r.Run()
}

// func TestRun() error {

// 	ecr := "654654451390.dkr.ecr.us-east-1.amazonaws.com/test:"

// 	repoName := strings.Split("Swechhya/Contact-Manager-Frontend", "/")[1]
// 	dest := fmt.Sprintf("%s%s-%s", ecr, "feature-new-test", repoName)
// 	err := services.GenerateBuildManifest("feature-new-test", "Swechhya/Contact-Manager-Frontend", "main", dest)
// 	if err != nil {
// 		return err
// 	}
// 	err = services.GenerateDeployManifest("feature-new-test", dest, nil)
// 	if err != nil {
// 		return err
// 	}
// 	err = services.DeployEnvironment("feature-new-test")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
