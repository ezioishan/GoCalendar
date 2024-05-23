package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"my-api/db"
	middleware "my-api/middlewares"
	"my-api/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Logger *log.Logger
var UserEventList map[string][]models.Events
var UserMap map[string]models.User

func main() {

	Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	godotenv.Load()
	conStr := os.Getenv("DB_CONN_STRING")
	dbInst, err := db.NewDatabase(conStr)
	if err != nil {
		Logger.Fatalf("db connection failed %s", conStr)
	}
	db.DB = dbInst
	Logger.Printf("db connection successful :: status : %d\n", db.DB.Stats().Idle)
	defer db.DB.Close()
	Logger.Println("Main - START")

	UserEventList = make(map[string][]models.Events)
	UserMap = make(map[string]models.User)
	r := mux.NewRouter()
	// TODO: r.HandleFunc("/api/signin", Login).Methods("POST")

	r.HandleFunc("/api/create_user", CreateUser).Methods("POST")
	r.HandleFunc("/api/login", Login).Methods("POST")

	r.Handle("/api/logout", middleware.Authenticate(http.HandlerFunc(Logout))).Methods("POST")
	r.Handle("/api/view_events", middleware.Authenticate(http.HandlerFunc(ViewEvents))).Methods("GET")
	r.Handle("/api/schedule_event", middleware.Authenticate(http.HandlerFunc(ScheduleEvent))).Methods("POST")
	r.Handle("/api/cancel_event", middleware.Authenticate(http.HandlerFunc(CancelEvent))).Methods("POST")
	r.Handle("/api/update_event", middleware.Authenticate(http.HandlerFunc(UpdateEvent))).Methods("POST")

	port := os.Getenv("PORT")
	fmt.Printf("Server listening at port %s\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, r))
	Logger.Println("Main - END")

}
