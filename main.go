package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Message represents a message sent by a user
type Message struct {
	Sender   string
	Receiver string
	Content  string
}

// User represents a user in the system
type User struct {
	Username string
	Online   bool
	Messages []Message
}

// Data to be rendered in the template
type TemplateData struct {
	Users []User
}

var (
	users []User
)

func main() {
	// Define routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/send", sendMessageHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	fmt.Println("Server started at localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Render template with user list
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Users: users,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Get username from request
	username := r.FormValue("username")

	// Create a new user
	newUser := User{
		Username: username,
		Online:   true,
		Messages: []Message{},
	}

	// Add user to the list
	users = append(users, newUser)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Get sender, receiver, and message from request
	sender := r.FormValue("sender")
	receiver := r.FormValue("receiver")
	content := r.FormValue("content")

	// Find sender and receiver in the user list
	var senderUser, receiverUser *User
	for i := range users {
		if users[i].Username == sender {
			senderUser = &users[i]
		}
		if users[i].Username == receiver {
			receiverUser = &users[i]
		}
	}

	// If sender or receiver not found, return error
	if senderUser == nil || receiverUser == nil {
		http.Error(w, "Sender or receiver not found", http.StatusBadRequest)
		return
	}

	// Create message
	message := Message{
		Sender:   sender,
		Receiver: receiver,
		Content:  content,
	}

	// Add message to receiver's messages
	receiverUser.Messages = append(receiverUser.Messages, message)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
