package services

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/key"
	"github.com/google/go-github/github"
)

func AuthenticateGitHub() {

	fmt.Println(os.Getwd())

	// load from a file
	key, err := key.FromFile("./key.pem")
	fmt.Println(err)

	install, err := inst.NewConfig("app_id", "install_id", key)

	ctx := context.Background()
	client1 := install.Client(ctx)

	// 49284279

	r, err := client1.Get("https://api.github.com/installation/repositories")
	b, err := io.ReadAll(r.Body)
	fmt.Println(b)
	fmt.Println(string(b))

	client := github.NewClient(client1)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Printf("Error fetching user: %v\n", err)
		os.Exit(1)
	}
	repos, _, err := client.Repositories.ListBranches(ctx, *user.Login, "panik-backend", nil)
	fmt.Println(repos)

}
