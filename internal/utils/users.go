package utils

import (
	"log"
	database "mw8/database"
	"mw8/internal/models"
	"time"
)

var TempUsers []models.TempUser
var TempModifUsers []models.TempUser

// removeUser
// remove the models.User which models.User.Id is sent in argument from jsonFile.
/* func removeUser(id int) {
	users, err := database.retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
	changeUsers(users)
} */

// updateUser
// modifies the models.User in jsonFile that matches
// `updatedUser`'s Id with `updatedUser`'s content.
/* func updateUser(updatedUser models.User) {
	users, err := database.retrieveUsers()
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	}
	for i, user := range users {
		if user.Id == updatedUser.Id {
			users[i] = updatedUser
		}
	}
	changeUsers(users)
} */

func deleteTempUser(temp models.TempUser) {
	for i, user := range TempUsers {
		if user == temp {
			TempUsers = append(TempUsers[:i], TempUsers[i+1:]...)
		}
	}
}

func PushTempUser(id string) {
	log.Printf("TempUsers: %#v\n", TempUsers)
	log.Printf("id: %#v\n", id)
	for _, temp := range TempUsers {
		if temp.ConfirmID == id {
			temp.User.Id = database.GetIdNewUser()
			database.CreateUser(temp.User)
			deleteTempUser(temp)
		}
	}
}

func ManageTempUsers() {
	duration := setDailyTimer()
	for {
		for _, user := range TempUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				deleteTempUser(user)
			}
		}
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}
func ManageTempModifUsers() {
	duration := setDailyTimer()
	for {
		for _, user := range TempModifUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				deleteTempModifUsers(user)
			}
		}
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}
func deleteTempModifUsers(temp models.TempUser) {
	for i, user := range TempModifUsers {
		if user == temp {
			TempModifUsers = append(TempModifUsers[:i], TempModifUsers[i+1:]...)
		}
	}
}
func PushTempModifUser(id string) {
	log.Println(GetCurrentFuncName())
	log.Printf("PushTempModifUser id= %#v\n", id)
	for _, temp := range TempModifUsers {
		if temp.ConfirmID == id {
			database.UpdateUser(temp.User)
			deleteTempModifUsers(temp)
		}
	}
}
