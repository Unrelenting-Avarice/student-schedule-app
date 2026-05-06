package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Global variables to store the logged-in user's information
var username string
var adminPrivileges bool

type CalendarData struct {
	Days          []string
	Hours         []string
	HourCount     int
	SlicesPerHour int
	Status        string
	Selected      map[string]bool
	ReadOnly      bool
	admin         bool
	Username      string
}

type ScheduleData struct {
	Files []string
}

var funcMap = template.FuncMap{
	"until": until,
}

var calendarTmpl = template.Must(
	template.New("calendar.html").
		Funcs(funcMap).
		ParseFiles("templates/calendar.html"),
)

var scheduleTmpl = template.Must(
	template.New("schedule.html").
		Funcs(funcMap).
		ParseFiles("templates/schedule.html"),
)

func until(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i
	}
	return a
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/calendar", calendarHandler)
	http.HandleFunc("/save-selection", saveSelectionHandler)
	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/view-calendar", viewCalendarHandler)
	http.HandleFunc("/modify-status", modifyStatus)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
func readStatusFromFile(file string) (string, error) {
	filePath := "saveFiles/" + file // just use the filename exactly
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("file is empty")
}

func menuHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "templates/menu.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	username = r.FormValue("username")
	password := r.FormValue("password")

	if password == "student12" { // Student credentials
		adminPrivileges = false
		http.Redirect(w, r, "/menu?user="+username, http.StatusSeeOther)
		return
	} else if username == "manager" && password == "admin12" { // Admin credentials
		adminPrivileges = true
		http.Redirect(w, r, "/menu?user="+username, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func loadSelections(filename string) map[string]bool {
	file, err := os.Open(filename)
	if err != nil {
		return map[string]bool{} // empty if no file
	}
	defer file.Close()

	selected := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		selected[line] = true
	}

	return selected
}

func calendarHandler(w http.ResponseWriter, r *http.Request) {

	status, err := readStatusFromFile(username)
	if err != nil {
		status = "NA"
	}

	hours := []string{}
	for h := 8; h <= 18; h++ {

		period := "AM"
		hour := h

		if h == 12 {
			period = "PM"
		} else if h > 12 {
			hour = h - 12
			period = "PM"
		}

		hours = append(hours, fmt.Sprintf("%d %s", hour, period))
	}

	selected := loadSelections("saveFiles/" + username + ".txt")

	data := CalendarData{
		Days:          []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
		Hours:         hours,
		HourCount:     len(hours),
		SlicesPerHour: 6,
		Status:        status,
		Selected:      selected,
		ReadOnly:      false,
		admin:         adminPrivileges,
		Username:      username,
	}

	calendarTmpl.Execute(w, data)
}

func saveSelectionHandler(w http.ResponseWriter, r *http.Request) {

	var selected []struct {
		Day   string `json:"day"`
		Hour  string `json:"hour"`
		Slice int    `json:"slice"`
	}

	json.NewDecoder(r.Body).Decode(&selected)

	var dataList []string                  // nil slice to hold the formatted data
	dataList = append(dataList, "Pending") // Add the approval status as the first line <<< CRUCIAL

	for _, s := range selected {
		data := fmt.Sprintf("%s-%s-%d", s.Day, s.Hour, s.Slice) // Format the data as "Day-Hour-Slice"
		dataList = append(dataList, data)                       // Append the formatted string to the slice
		// fmt.Printf("%s\n", data)                                // Three letter code -- hr Index -- slice index

	}

	content := strings.Join(dataList, "\n")                                  // Join the slice of strings into a single string with newlines
	err := os.WriteFile("saveFiles/"+username+".txt", []byte(content), 0644) // Write the content to data.json
	if err != nil {                                                          // Handle any potential errors
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Submit!"))
}
func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	fileNames := []string{}

	files, err := os.ReadDir("saveFiles")
	if err != nil {
		fmt.Println("non nil")
		panic(err)
	}

	for _, file := range files {
		// fmt.Println(file.Name())
		fileNames = append(fileNames, file.Name())
	}

	data := ScheduleData{
		Files: fileNames,
	}

	scheduleTmpl.Execute(w, data)
}

func viewCalendarHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}

	status, err := readStatusFromFile(fileName) // <-- FIXED
	if err != nil {
		status = "NA"
	}

	hours := []string{}
	for h := 8; h <= 18; h++ {
		period := "AM"
		hour := h
		if h == 12 {
			period = "PM"
		} else if h > 12 {
			hour = h - 12
			period = "PM"
		}
		hours = append(hours, fmt.Sprintf("%d %s", hour, period))
	}

	selected := loadSelections("saveFiles/" + fileName)

	data := CalendarData{
		Days:          []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
		Hours:         hours,
		HourCount:     len(hours),
		SlicesPerHour: 6,
		Status:        status,
		Selected:      selected,
		ReadOnly:      true,
		admin:         adminPrivileges,
	}

	calendarTmpl.Execute(w, data)
}

func modifyStatus(w http.ResponseWriter, r *http.Request) {
	// Check admin via query para
	if adminPrivileges == false {
		http.Error(w, "Forbidden: admin required", http.StatusForbidden)
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")
	if action != "approve" && action != "deny" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	var status string
	switch action {
	case "approve":
		status = "Approved"
	case "deny":
		status = "Denied"
	}

	// Read existing selections to preserve them
	existingSelections := loadSelections("saveFiles/" + file)
	content := status + "\n"
	for sel := range existingSelections {
		content += sel + "\n"
	}

	err := os.WriteFile("saveFiles/"+file, []byte(content), 0644)
	if err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Status updated successfully"))
}
