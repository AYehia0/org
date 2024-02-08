package pkg

/*
The main package for the application.
This package is responsible for starting the application and setting up the necessary dependencies.
	- API layer
	- Database layer
	- Utils: like logger, config, etc.
*/

// export a interface type to represent the application
type App interface {
	Start() error
}

type app struct {
	Server *Server
}

// the function to create a new app
func NewApp(server *Server) App {
	return &app{Server: server}
}

// the function to start the app
func (a *app) Start() error {
	if err := a.Server.Init(); err != nil {
		return err
	}
	return a.Server.Run()
}

func StartApp() error {
	// create a new server
	server := &Server{}
	app := NewApp(server)
	return app.Start()
}
