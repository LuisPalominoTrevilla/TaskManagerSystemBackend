package controllers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"os"
	"strings"
	"io"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
	database "github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gorilla/mux"

)
// HabitsController holds important data for the habits controller
type HabitsController struct {
	habitsDB *database.HabitsDB
}
// Get serves as a GET request for either all or one specific habit
func (controller *HabitsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) < 1 {
		habits := controller.getAll()
		if habits == nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "Error retrieving habits from database.")
			return
		}
		response := struct {
			Habits []*models.Habit
		} {
			habits,
		}
		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(response)
        return
	}
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}
	filter := bson.D{{"_id", habitID}}
	var habit models.Habit
	err2 := controller.habitsDB.GetByID(filter, &habit)
	if(err2 != nil) {
		println(err2.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving habit from database.")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(habit)
}
func (controller *HabitsController) getAll() []*models.Habit {
	var results []*models.Habit
	filter := bson.D{}
	results, err := controller.habitsDB.Get(filter)
	if(err != nil) {
		results = nil
	}
	return results
}
func (controller *HabitsController) initializeController(r *mux.Router) {
	r.HandleFunc("/{id}", controller.Get).Methods(http.MethodGet)
	r.HandleFunc("/", controller.Get).Methods(http.MethodGet)
	r.HandleFunc("/", controller.CreateHabit).Methods(http.MethodPost)
	r.HandleFunc("/{id}", controller.EditHabit).Methods(http.MethodPut)
	r.HandleFunc("/{id}", controller.DeleteHabit).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", controller.SumPoints).Methods(http.MethodGet)
	r.HandleFunc("/{id}/{completionStatus}", controller.SumPoints).Methods(http.MethodGet)
}

// SetHabitsController sets the controller for the sets up the habits controllet
func SetHabitsController(r *mux.Router, db *mongo.Database) {
	habit := database.HabitsDB{Habits: db.Collection("habits")}
	
	habitsController := HabitsController{habitsDB: &habit}
	habitsController.initializeController(r)
}
// CreateHabit serves to receive a POST request to create a new habit
func (controller *HabitsController) CreateHabit(w http.ResponseWriter, r *http.Request) {
	var maxBytes int64 = 64 * 1024 * 1024
	validImageFormats := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
	}
	// Parse multipart form data
	err := r.ParseMultipartForm(maxBytes)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the data provided.")
		return
	}
	// Ensure that both the image and age are contained in multipartform
	if len(r.MultipartForm.File["image"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing image.")
		return
	}
	if len(r.MultipartForm.Value["title"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing title.")
		return
	}
	if len(r.MultipartForm.Value["type"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing type.")
		return
	}
	if len(r.MultipartForm.Value["difficulty"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing difficulty.")
		return
	}
	// get userEmail from header
	userId := r.Header.Get("userId")
	// get age from image
	hType, err := strconv.Atoi(r.MultipartForm.Value["type"][0])
	if (err != nil || hType < -1 || hType > 1){
		println(err.Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Either type was not a number or is not within the permitted range.")
		return
	}
	difficulty, err := strconv.Atoi(r.MultipartForm.Value["difficulty"][0])
	if (err != nil){
		println(err.Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Difficulty was not a number.")
		return
	}
	// get image file header
	imFileHeader := r.MultipartForm.File["image"][0]
	im, err := imFileHeader.Open()
	
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error opening the image file.")
		return
	}
	defer im.Close()
	if _, exists := validImageFormats[imFileHeader.Header["Content-Type"][0]]; !exists {
		w.WriteHeader(400)
		fmt.Fprint(w, "Image file format not valid.")
		return
	}
	if imFileHeader.Size/1000000 > 5 {
		w.WriteHeader(413)
		fmt.Fprint(w, "La imagen pesa más de 5 MB.")
		return
	}
	imageURL := "/" + parseUserEmail(userId)
	// ensure dir exists and create final file
	err2 := os.MkdirAll("static"+imageURL, os.ModePerm)
	if(err2 != nil){
		println(err2.Error())
	}
	imageURL += "/"
	filename := imFileHeader.Filename
	filename = strings.Replace(filename, " ", "", -1)
	additionalNum := ""
	for fileExists("static" + imageURL + additionalNum + filename) {
		if additionalNum == "" {
			additionalNum = "1"
		} else {
			num, _ := strconv.Atoi(additionalNum)
			num++
			additionalNum = strconv.Itoa(num)
		}
	}
	imageURL += additionalNum + filename
	file, err := os.Create("static" + imageURL,)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error creating file.")
		return
	}
	defer file.Close()
	// write image file to dir
	_, err = io.Copy(file, im)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error copying the image file.")
		return
	}
	habit := models.Habit{
		Title:     	r.MultipartForm.Value["title"][0],
		Type: 		hType,
		Difficulty: difficulty,
		UserId:		userId,
		Image:		"/images" + imageURL,
		Score:      0,
	}
	result, err := controller.habitsDB.Insert(habit)
	if err != nil {
		println(err.Error())
		_ = os.Remove(("/static" + imageURL))
		w.WriteHeader(500)
		fmt.Fprint(w, "Error uploading image to database.")
		return
	}
	habit.ID = result.InsertedID.(primitive.ObjectID)
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(habit)
}
func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}
func parseUserEmail(email string) string {
	emailParsed := ""
	for i:=0; i < len(email); i++ {
		if(email[i] == '@') {
			break
		} 
		emailParsed += string(email[i])
	}
	return emailParsed
}
// EditHabit serves to receive a PUT request to modify an existing habit
func (controller *HabitsController) EditHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}
	// Parse url-encoded data
	err = r.ParseForm()
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the data provided.")
		return
	}
	if len(r.Form["title"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing title.")
		return
	}
	if len(r.Form["type"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing type.")
		return
	}
	if len(r.Form["difficulty"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing difficulty.")
		return
	}
	// get age from image
	hType, err := strconv.Atoi(r.Form["type"][0])
	if (err != nil || hType < -1 || hType > 1){
		println(err.Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Either type was not a number or is not within the permitted range.")
		return
	}
	difficulty, err := strconv.Atoi(r.Form["difficulty"][0])
	if (err != nil){
		println(err.Error())
		w.WriteHeader(400)
		fmt.Fprint(w, "Difficulty was not a number.")
		return
	}
	filter := bson.D{{"_id", habitID}}
	updateDoc := bson.D{{"$set", bson.D{{"title", r.Form["title"][0]},
		{"type", hType},
		{"difficulty", difficulty}}}}
	_, err = controller.habitsDB.UpdateOne(filter, updateDoc)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error updating data in the database.")
		return
	}
	var result models.Habit
	err = controller.habitsDB.GetByID(filter, &result)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving new values.")
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func (controller *HabitsController) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}
	filter := bson.D{{"_id", habitID}}
	deleteResult, err := controller.habitsDB.DeleteOne(filter)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while deleting habit.")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(deleteResult)
}

func (controller * HabitsController) SumPoints(w http.ResponseWriter, r * http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}

	status := vars["completionStatus"]
	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No completion status was provided.")
		return
	}

	completionStatus, err := strconv.Atoi(status)

	if err != nil || completionStatus > 1 || completionStatus < 0{
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the completion status: either not a number or not within range [0,1].")
		return
	}

	var habit models.Habit

	err = controller.habitsDB.GetByID(bson.D{{"_id", habitID}}, &habit)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving habit from the database.")
		return
	}

	if (habit.Type == 1 && completionStatus == 0) || (habit.Type == -1 && completionStatus == 1) {
		w.WriteHeader(400)
		fmt.Fprint(w, "Incompatible habit type and completion status.")
		return
	}

	//get score as int
	score := habit.Score

	//sumar difficulty a score
	if(score < 0){ // es rojo
		if(habit.Type == 1){
			score += float32(habit.Difficulty)
		}
		if(habit.Type == -1){
			score -= float32(habit.Difficulty)*2
		}
	} else if(score >= 0 && score < 10){ // es naranja
		if(habit.Type == 1){
			score += float32(habit.Difficulty)
		}
		if(habit.Type == -1){
			score -= float32(habit.Difficulty)*1.5
		}
	} else if(score >= 10 && score < 40){ // es amarillo 
		if(habit.Type == 1){
			score += float32(habit.Difficulty)
		}
		if(habit.Type == -1){
			score -= float32(habit.Difficulty)
		}
	} else if(score >= 40 && score < 50){
		if(habit.Type == 1){
			score += float32(habit.Difficulty)*1.5
		}
		if(habit.Type == -1){
			score -= float32(habit.Difficulty)
		}
	} else{
		if(habit.Type == 1){
			score += float32(habit.Difficulty)*2
		}
		if(habit.Type == -1){
			score -= float32(habit.Difficulty)
		}
	}

}