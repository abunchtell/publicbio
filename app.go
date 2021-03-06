package publicbio

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/writeas/web-core/converter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	serverSoftware = "Public Bio"
	softwareURL    = "https://publicb.io"
)

var (
	// Software version can be set from git env using -ldflags
	softwareVer = "0.1.0"
)

type App struct {
	cfg *Config

	singleUser *Profile
}

func (app *App) multiUser() bool {
	return app.singleUser == nil
}

type Config struct {
	Host   string
	Port   int
	static bool

	UserFile string
}

func Serve(cfg *Config) {
	app := &App{
		cfg: cfg,
	}

	if cfg.UserFile != "" {
		f, err := ioutil.ReadFile(cfg.UserFile)
		if err != nil {
			log.Fatalf("File error: %v", err)
		}

		err = json.Unmarshal(f, &app.singleUser)
		if err != nil {
			log.Fatalf("Unable to read user config: %v", err)
		}
		fmt.Printf("Results: %v\n", app.singleUser)
	} else {
		log.Fatal("No user configuration")
	}

	if app.cfg.static {
		if err := renderTemplate(os.Stdout, "profile", app.singleUser); err != nil {
			log.Fatal(err)
		}
		return
	}

	r := mux.NewRouter()
	app.InitRoutes(r)

	http.Handle("/", r)
	log.Printf("Serving on http://localhost:%d", app.cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.Port), nil)
}

func initConverter() {
	formDecoder := schema.NewDecoder()
	formDecoder.RegisterConverter(converter.NullJSONString{}, converter.JSONNullString)
}

// FormatVersion constructs the version string for the application
func FormatVersion() string {
	return serverSoftware + " " + softwareVer
}
