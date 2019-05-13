package controllers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"os"
	"strings"
	"io"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
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

// Get serves as a GET request
func (controller *HabitsController) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
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

	habitID, err := primitive.ObjectIDFromHex(id[0])
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
	r.HandleFunc("/", controller.Get).Methods(http.MethodGet)
	r.HandleFunc("/CreateHabit", controller.CreateHabit).Methods(http.MethodPost)
}

// SetHabitsController sets the controller for the sets up the habits controllet
func SetHabitsController(r *mux.Router, db *mongo.Database) {
	habit := database.HabitsDB{Habits: db.Collection("habits")}
	
	habitsController := HabitsController{habitsDB: &habit}
	habitsController.initializeController(r)
}

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
	/** if len(r.MultipartForm.Value["userEmail"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing owner's email.")
		return
	} */

	// get userEmail from header
	userEmail := r.Header.Get("userEmail")

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
	defer im.Close()
	if err != nil {
		println(err.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error opening the image file.")
		return
	}

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

	imageURL := "/" + parseUserEmail(userEmail)

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
		UserEmail:	userEmail,
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
	habit.Image = ""
	habit.UserEmail = ""
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