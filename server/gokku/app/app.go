package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"magnax.ca/gokku/server/gokku"
	"os"
	"os/exec"
	"path"
)

type App struct {
	Name string
}

func NewApp(name string) *App {
	app := new(App)
	app.Name = name
	return app
}

func UnmarshalApp(data []byte) *App {
	app := new(App)
	err := yaml.Unmarshal(data, struct{ App *App }{app})
	if err != nil {
		panic(err)
	}
	return app
}

func LoadApp(name string) (*App, error) {
	data, err := ioutil.ReadFile(path.Join(gokku.CurrentConfig.AppDirectory, name, "app.yml"))
	if err != nil {
		return nil, err
	}
	return UnmarshalApp(data), nil
}

func MustLoadApp(name string) *App {
	app, err := LoadApp(name)
	if err != nil {
		panic(err)
	}
	return app
}

func (a *App) Create() (err error) {
	err = os.Mkdir(a.Path(), 0770)
	if err != nil {
		return fmt.Errorf("app: could not create app directory: %v\n", err)
	}
	// from here on out, an error must make us panic, so that we can remove the directory as needed.
	defer func() {
		if r := recover(); r != nil {
			_ = os.RemoveAll(a.Path())
			err = fmt.Errorf("app: %v", r)
		}
	}()

	cmd := exec.Command("git", "init", "--quiet", "--bare", a.RepositoryPath())
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = a.Save()
	if err != nil {
		panic(err)
	}

	return nil
}

func (a *App) Marshal() []byte {
	data, err := yaml.Marshal(struct{ App *App }{a})
	if err != nil {
		panic(err)
	}
	return data
}

func (a *App) Save() error {
	data := a.Marshal()
	err := ioutil.WriteFile(a.ConfigPath(), data, 0660)
	return err
}

func (a *App) Path() string {
	return path.Join(gokku.CurrentConfig.AppDirectory, a.Name)
}

func (a *App) ConfigPath() string {
	return path.Join(gokku.CurrentConfig.AppDirectory, a.Name, "app.yml")
}

func (a *App) RepositoryPath() string {
	return path.Join(gokku.CurrentConfig.AppDirectory, a.Name, "repo.git")
}
