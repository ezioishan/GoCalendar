package main

import (
	"encoding/json"
	"net/http"
	"time"

	"my-api/models"

	"my-api/db"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

/*
  - API Endpoints:
  - ViewEvents :
    parameters : start_time, end_time
    url: /api/start
  - ScheduleEvent :
  - CancelEvent
  - UpdateEvent

st1:2024-03-29T22:30:00.000Z
et1:2024-03-29T22:45:00.000Z

st2:2024-03-29T22:40:00.000Z
et2:2024-03-29T22:50:00.000Z

st3:2024-03-29T23:50:00.000Z
et3:2024-03-30T00:00:00.000Z

CreateUser:
Example:

	{
		"name" : "Ishan Sourav",
		"email" : "ishansourav7@gmail.com",
		"id" : "03145",
		"hash" : ""
	}
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	Logger.Println(db.DB.Stats().Idle)
	//get password hash
	newUser.Hash = newUser.HashPassword("password")
	Logger.Printf("Hash : %s", newUser.Hash)
	// db.DB.Query(query)
	// 	err := db.DB.Create(&newUser).Error
	// if err != nil {
	// 	// customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Error: "+err.Error())
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&newUser)
	Logger.Printf("New user created : %s", newUser.ID)
	UserMap[newUser.ID] = newUser
	query := `
		INSERT INTO users (id, name, email, hash)
		VALUES ($1, $2, $3, $4)
		`

	_, err := db.DB.Exec(query, newUser.ID, newUser.Name, newUser.Email, newUser.Hash)
	if err != nil {
		Logger.Panicf("insert query failed %s", query)
		return
	}
}

/*
Login:
Example:
user_id = "03145"


*/
// POST Method
func Login(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	//fetch user details from database
	query := ` SELECT name, email, id, hash
		FROM users where id = $1
	`
	result := db.DB.QueryRow(query, userId)
	if result.Err() != nil {
		Logger.Panicf("select query failed %s", query)
		return
	}
	var user models.User
	result.Scan(&user.Name, &user.Email, &user.ID, &user.Hash)

	// user, ok := UserMap[userId]
	// if !ok {
	// 	Logger.Printf("UserId : %s not found", userId)
	// 	http.Error(w, "UserId not found", http.StatusNotFound)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	// if password matches - generate a token
	if user.CheckPassword("password") {
		token, err := user.GenerateJWT()
		if err != nil {
			error := ErrorResponse{
				true,
				"Generate JWT failure : " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusFailedDependency)
			json.NewEncoder(w).Encode(&error)
			return
		}
		Logger.Printf("Login successful %s", userId)
		json.NewEncoder(w).Encode(&token)
		// TODO: store token in a map : {user_id : token}
	} else {
		// customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Password incorrect")
		Logger.Panic("Password does not match")
	}
}

// POST Method
func Logout(w http.ResponseWriter, r *http.Request) {
	var newEvent models.Events
	_ = json.NewDecoder(r.Body).Decode(&newEvent)
}

/*
TODO: Create a table Events

Create table Events
(
	ID PRIMARY_KEY,
	Title VARCHAR(),
	Description,
	Start_Date,
	End_Date,
	CreatorID FOREIGN_KEY--same as user id
)


during view events:
	res = select * from Events where creatorID = userID
	print res

schedule event:

only when its not overlapping with other events
so for that we need to query all the existing events and check for overlapping

insert into Events (...)
	values (...)
*/
// GET Method
func ViewEvents(w http.ResponseWriter, r *http.Request) {
	Logger.Print("ViewEvents - START")
	// username, ok := r.Context().Value("username").(string)
	userId := r.Header.Get("userId")

	if userId == "" {
		Logger.Print("User Token not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse start and end time parameters from the query string
	startStr := r.URL.Query().Get("start_time")
	endStr := r.URL.Query().Get("end_time")
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time", http.StatusBadRequest)
		return
	}

	// Filter bookings within the specified time interval
	var filteredEvents []models.Events
	query := `
		SELECT id, title, description, start_time, end_time, audience_list
		FROM events
		WHERE user_id = $1	
	`
	result, err := db.DB.Query(query, userId)
	if err != nil {
		Logger.Panicf("insert query failed %s", query)
		return
	}
	defer result.Close()
	for result.Next() {
		var event models.Events
		err = result.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.AudienceList)
		event.UserId = userId
		if err != nil {
			Logger.Panic("Error while querying EVENTS table")
		}
		if event.StartTime.After(startTime) && event.EndTime.Before(endTime) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	// for _, event := range UserEventList[userId] {
	// 	if event.StartTime.After(startTime) && event.EndTime.Before(endTime) {
	// 		filteredEvents = append(filteredEvents, event)
	// 	}
	// }

	json.NewEncoder(w).Encode(filteredEvents)
	Logger.Print("ViewEvents - END")
}

// POST Method
func ScheduleEvent(w http.ResponseWriter, r *http.Request) {
	Logger.Print("ScheduleEvent - START")

	// username, ok := r.Context().Value("username").(string)
	userId := r.Header.Get("userId")
	if userId == "" {
		Logger.Print("User Token not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var newEvent models.Events
	_ = json.NewDecoder(r.Body).Decode(&newEvent)

	eventAllUserList := newEvent.AudienceList
	eventAllUserList = append(eventAllUserList, userId)

	// check for all the users whether the new event will have any conflicts with existing ones
	for _, id := range eventAllUserList {
		query := `
			SELECT id, title, description, start_time, end_time, audience_list
			FROM events
			WHERE user_id = $1
			AND end_time >= $2
		`
		result, err := db.DB.Query(query, id)
		if err != nil {
			Logger.Panicf("select query failed %s", query)
			return
		}
		defer result.Close()
		for result.Next() {
			var event models.Events
			err = result.Scan(&event.ID,
				&event.Title,
				&event.Description,
				&event.StartTime,
				&event.EndTime,
				&event.AudienceList)
			event.UserId = userId
			if err != nil {
				Logger.Panic("Error while querying EVENTS table")
			}
			if event.StartTime.After(newEvent.StartTime) && event.EndTime.Before(newEvent.EndTime) {
				// overlap
				Logger.Panicf("new event overlaps for user %s", id)
				return
			}
		}
	}

	// does not overlap
	query := `
		INSERT INTO events (id, title, description, start_time, end_time, user_id, audience_list)
		VALUES ($1, $2, $3, $4, $5)
		`

	_, err := db.DB.Exec(query, newEvent.ID, newEvent.Title, newEvent.Description, newEvent.StartTime, newEvent.EndTime, userId, newEvent.AudienceList)
	if err != nil {
		Logger.Panicf("insert query failed %s", query)
		return
	}
	// UserEventList[userId] = append(UserEventList[userId], newEvent)
	Logger.Printf("New event added : %s", newEvent.ID)
	Logger.Print("ScheduleEvent - END")
}

func CancelEvent(w http.ResponseWriter, r *http.Request) {
	Logger.Print("CancelEvent - START")
	// username, ok := r.Context().Value("username").(string)
	// if !ok {
	// 	Logger.Print("User Token not found")
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	// Parse event ID from URL parameters
	params := mux.Vars(r)
	eventId := params["eventId"]
	userId := params["userId"]

	// Find the event with the specified ID
	query := `
			SELECT id, title, description, start_time, end_time, user_id, audience_list
			FROM events
			WHERE id = $1
		`
	result := db.DB.QueryRow(query, eventId)
	if result.Err() != nil {
		Logger.Panicf("select query failed %s", query)
		return
	}

	var event models.Events
	result.Scan(&event.ID, &event.Description, &event.StartTime, &event.EndTime, &event.UserId, &event.AudienceList)

	// if cancel req is sent from the creator of the event then delete the event
	// otherwise remove the user id from the audient_list and update the new list

	if userId == event.UserId {
		// creator
		// delete query
	} else {
		// remove user from the audience list
		// update query
	}
	// var found bool
	// for i, event := range UserEventList[username] {
	// 	if event.ID == eventID {
	// 		// Remove the event from the slice
	// 		UserEventList[username] = append(UserEventList[username][:i], UserEventList[username][i+1:]...)
	// 		found = true
	// 		break
	// 	}
	// }

	// Check if event with the specified ID was found
	// if !found {
	// 	http.Error(w, "Event not found", http.StatusNotFound)
	// 	return
	// }

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Booking cancelled successfully"))
	Logger.Print("CancelEvent - END")

}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	Logger.Print("UpdateEvent - START")
	Logger.Print("UpdateEvent - END")

}
