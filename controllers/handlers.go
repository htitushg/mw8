package controllers

import (
	database "Mw7/database"
	"Mw7/internal/middlewares"
	"Mw7/internal/models"
	"Mw7/internal/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func indexHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	err := models.Tmpl["index"].ExecuteTemplate(w, "base", "indexHandlerGet --- Restricted area! ---")
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerPut(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	sessionID, _ := r.Cookie("updatedCookie")
	err := models.Tmpl["index"].ExecuteTemplate(w, "base", "indexHandlerPut"+sessionID.Value+"\nUsername: "+utils.SessionsData[sessionID.Value].Pseudo+"\nIP address: "+utils.SessionsData[sessionID.Value].IpAddress)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	sessionID, _ := r.Cookie("updatedCookie")
	err := models.Tmpl["index"].ExecuteTemplate(w, "base", "indexHandlerPut"+sessionID.Value+"\nUsername: "+utils.SessionsData[sessionID.Value].Pseudo+"\nIP address: "+utils.SessionsData[sessionID.Value].IpAddress)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerNoMeth(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	utils.Logger.Warn("indexHandlerNoMeth", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusMethodNotAllowed))
	message := fmt.Sprintf("HTTP Error %v", http.StatusMethodNotAllowed)
	err := models.Tmpl["index"].ExecuteTemplate(w, "base", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	utils.Logger.Warn("indexHandlerOther", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))
	message := fmt.Sprintf("This adress is not valid : Error %v", http.StatusNotFound)
	err := models.Tmpl["error404"].ExecuteTemplate(w, "base", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "login":
			message = "<div class=\"message\">Wrong Pseudo or password!</div>"
		case "restricted":
			message = "<div class=\"message\">You need to login to access that area!</div>"
		}
	}
	err := models.Tmpl["login"].ExecuteTemplate(w, "base", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	credentials := models.Credentials{
		Username: r.FormValue("Pseudo"),
		Password: r.FormValue("password"),
	}
	if utils.CheckPwd(credentials) {
		utils.OpenSession(&w, credentials.Username, r)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login?err=login", http.StatusSeeOther)
	}
}

func registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "pseudo":
			message = "<div class=\"message\">Username must be at least 8 characters long! with first 1 upper or lower case letter, can have digits!</div>"
		case "name":
			message = "<div class=\"message\">Pseudo already used!</div>"
		case "password":
			message = "<div class=\"message\">Both passwords need to be equal!</div>"
		case "email":
			message = "<div class=\"message\">Wrong email value!</div>"
		case "pass":
			message = "<div class=\"message\">Your password must contain at least one lowercase letter, uppercase letter, a number and a special character and have a minimum length of 8 characters!</div>"
		}
	}
	err := models.Tmpl["register"].ExecuteTemplate(w, "base", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	formValues := struct {
		username  string
		email     string
		password1 string
		password2 string
	}{
		username:  r.FormValue("pseudo"),
		email:     strings.TrimSpace(strings.ToLower(r.FormValue("email"))),
		password1: r.FormValue("password1"),
		password2: r.FormValue("password2"),
	}
	switch {
	case !utils.CheckPseudo(formValues.username):
		http.Redirect(w, r, "register?err=pseudo", http.StatusSeeOther)
		return
	case database.UserOrEmailExist(formValues.username, formValues.email):
		http.Redirect(w, r, "register?err=name", http.StatusSeeOther)
		return
	case formValues.password1 != formValues.password2:
		http.Redirect(w, r, "register?err=password", http.StatusSeeOther)
		return
	case !utils.CheckEmail(formValues.email):
		http.Redirect(w, r, "register?err=email", http.StatusSeeOther)
		return
	case !utils.CheckPasswd(formValues.password1):
		http.Redirect(w, r, "register?err=pass", http.StatusSeeOther)
		return
	}
	hash, salt := utils.NewPwd(formValues.password1)
	newTempUser := models.TempUser{
		ConfirmID:    "",
		CreationTime: time.Now(),
		User: models.User{
			Id:        0,
			Pseudo:    formValues.username,
			HashedPwd: hash,
			Salt:      salt,
			Email:     formValues.email,
		},
	}
	utils.SendMail(&newTempUser, "mail.gohtml")
	utils.TempUsers = append(utils.TempUsers, newTempUser)
	log.Printf("newTempUser: %#v\n", newTempUser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Printf("homeHandlerGet avant appel Tmpl vers base")
	err := models.Tmpl["index"].ExecuteTemplate(w, "base", "homeHandlerGet --- Restricted area! ---")
	if err != nil {
		log.Fatalln(err)
	}
}

func logHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("level") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("level", r.URL.Query().Get("level")))
		return
	} else if r.URL.Query().Has("user") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("user", r.URL.Query().Get("user")))
		return
	}
	json.NewEncoder(w).Encode(utils.RetrieveLogs())
}

func confirmHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	if r.URL.Query().Has("id") {
		id := r.URL.Query().Get("id")
		// Traiter la confirmation
		utils.PushTempUser(id)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func confirmupdateHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	if r.URL.Query().Has("id") {
		id := r.URL.Query().Get("id")
		// Traiter la mise à jour de la confirmation
		log.Printf("confirmupdateHandlerGet,id= '%v'\n", id)
		utils.PushTempModifUser(id)
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
func logoutHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	utils.Logout(&w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func ModifUserHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	// Il faut vérifier si l'utilisateur est bien connecté
	// Il faut ici collecter les infos de l'utilisateur
	// créons une structure pour cela
	session, _ := utils.GetSession(r)
	log.Println("début de r.URL.Query().Has('id')")
	// Récupération des infos utilisateur
	user, ok := database.SelectUser(session.Pseudo)

	formValues := struct {
		Username string
		Email    string
		Message  string
	}{
		Username: user.Pseudo,
		Email:    user.Email,
		Message:  "",
	}
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "pseudo":
			formValues.Message = "Username must be at least 8 characters long! with first 1 upper or lower case letter, can have digits!"
		case "name":
			formValues.Message = "Pseudo already used!"
		case "pass1":
			formValues.Message = "Both new passwords need to be equal!"
		case "pass2":
			formValues.Message = "Actual password value is Wrong !"
		}
	}
	log.Println("suite de r.URL.Query().Has('id') avant ok, formValues.Username=", formValues.Username)
	if ok {
		err := models.Tmpl["modifuser"].ExecuteTemplate(w, "base", formValues)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
func ModifUserHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	// Il faut ici récupérer les valeurs du formulaire
	Id, _ := strconv.Atoi(r.FormValue("id"))
	session, _ := utils.GetSession(r)
	formValues := struct {
		id             int
		username       string
		email          string
		actualpassword string
		password1      string
		password2      string
	}{
		id:             Id,
		username:       r.FormValue("pseudo"),
		email:          strings.TrimSpace(strings.ToLower(r.FormValue("email"))),
		actualpassword: strings.TrimSpace(r.FormValue("actualpassword")),
		password1:      strings.TrimSpace(r.FormValue("password1")),
		password2:      strings.TrimSpace(r.FormValue("password2")),
	}
	user, _ := database.SelectUser(session.Pseudo)
	var Utilisateur models.Credentials
	var newTempModifUser models.TempUser
	var emailchanged bool
	var passchanged bool
	var pseudochanged bool
	Utilisateur.Password = formValues.actualpassword
	Utilisateur.Username = session.Pseudo
	log.Printf("ActuelPassword='%v', Password1='%v', Password2='%v'\n", formValues.actualpassword, formValues.password1, formValues.password2)
	if utils.CheckPwd(Utilisateur) {
		log.Printf("Le mot de passe actuel est le bon\n")
		log.Printf("user.Pseudo ='%v', formValues.username ='%v'\n", user.Pseudo, formValues.username)
		log.Printf("user.Email ='%v', formValues.email ='%v'\n", user.Email, formValues.email)
		log.Printf("ActuelPassword='%v', Password1='%v', Password2='%v'\n", formValues.actualpassword, formValues.password1, formValues.password2)

		if user.Pseudo != formValues.username {
			log.Printf("Le nom de l'utilisateur à changé: %#v, au lieu de %v\n", formValues.username, user.Pseudo)
			// Il faut vérifier que le nouveau nom d'utilisateur n'existe pas
			if !database.UserExist(formValues.username) {
				if !utils.CheckPseudo(formValues.username) {
					http.Redirect(w, r, "modifuser?err=pseudo", http.StatusSeeOther)
					return
				} else {
					user.Pseudo = formValues.username
					pseudochanged = true
				}
			} else {
				http.Redirect(w, r, "modifuser?err=name", http.StatusSeeOther)
				return
			}
		}
		if user.Email != formValues.email {
			log.Printf("L'email de l'utilisateur à changé: %#v, au lieu de %v\n", formValues.email, user.Email)
			if !database.EmailExist(formValues.email) {
				emailchanged = true
				user.Email = formValues.email
			} else {
				http.Redirect(w, r, "modifuser?err=email", http.StatusSeeOther)
				return
			}
		}
		if formValues.password1 != "" || formValues.password2 != "" {
			if (formValues.password1 == formValues.password2) && (formValues.password1 != formValues.actualpassword) {
				log.Printf("Les nouveaux mots de passes sont identiques et différents de l'actuel\n")
				hash, salt := utils.NewPwd(formValues.password1)
				log.Printf("valeurs utilisateur %#v\n", user)
				// Le nouveau mot de passe est différent de l'ancien
				user.HashedPwd = hash
				user.Salt = salt
				passchanged = true
			} else {
				http.Redirect(w, r, "modifuser?err=pass1", http.StatusSeeOther)
				return
			}
		}
		if emailchanged {
			// il faut valiser le nouvel email
			// il faut créer un utilisateur temporaire : newTempModifUser
			log.Printf("L'email de l'utilisateur à changé: %#v, au lieu de %v\n", formValues.email, user.Email)
			user.Email = formValues.email
			newTempModifUser.ConfirmID = utils.GenerateConfirmationID()
			newTempModifUser.CreationTime = time.Now()
			newTempModifUser.User = user

			utils.SendMail(&newTempModifUser, "mail2.gohtml")
			utils.TempModifUsers = append(utils.TempModifUsers, newTempModifUser)
			log.Printf("ModifUserHandlerPost: newTempModifUser.ConfirmID= %#v\n", newTempModifUser.ConfirmID)
			log.Printf("ModifUserHandlerPost: newTempModifUser.User.Pseudo= %#v\n", newTempModifUser.User.Pseudo)
			http.Redirect(w, r, "/logout", http.StatusSeeOther) // provisoire
			return
		}
		if (pseudochanged || passchanged) && !emailchanged {
			log.Printf("Le nom de l'utilisateur à changé: %#v, au lieu de %v\n", formValues.username, user.Pseudo) //testing
			// Si l'email n'est pas modifié, il n'est pas nécessaire de le revalider
			// Par contre il faut accuser réception de la modification des données utilisateur par courriel
			// sans demander de réponse
			// il faut mettre à jour l'utilisateur et éventuellement son mot de passe
			// mais pas son email
			// l'Id de l'utilisateur reste le même
			database.UpdateUser(user)
			utils.SendMailUpdate(&user)
			// il faut déconnecter l'utilisateur
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
		}

		log.Printf("valeurs utilisateur %#v\n", user)
	} else {
		log.Printf("le mot de passe entré n'est pas le bon\n")
		http.Redirect(w, r, "modifuser?err=pass2", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
