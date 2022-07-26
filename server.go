package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Student struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func main() {
	fmt.Println("Starting server at port 8080...")
	http.HandleFunc("/", AnonymousDataHandler)
	http.HandleFunc("/add-student-data", AddStudentDataHandler)
	http.HandleFunc("/fetch-student-data", FetchStudentDataHandler)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Println("Error listening to the server!")
		return
	}
}

func AnonymousDataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "Invalid! Please hit a valid request!")
}

func AddStudentDataHandler(w http.ResponseWriter, r *http.Request) {
	stud := Student{}
	if err := json.NewDecoder(r.Body).Decode(&stud); err != nil {
		fmt.Println("Error while decoding the request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	studentDetails, _ := json.Marshal(stud)

	// Store student details to a file
	err := os.WriteFile("student_data.txt", studentDetails, 0777)
	if err != nil {
		fmt.Println("Could not write student data to the file! os.WriteFile() failed with an error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Println("Student data is successfully stored in the file:", string(studentDetails))
	}
}

func FetchStudentDataHandler(w http.ResponseWriter, r *http.Request) {
	studentDetails, err := os.ReadFile("student_data.txt")
	if err != nil {
		fmt.Println("Could not read student data from the file! os.ReadFile() failed with an error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err = json.NewEncoder(w).Encode(studentDetails); err != nil {
		fmt.Println("Error while encoding the student details:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "Application/json")
	fmt.Println("Student data is successfully fetched from the file:", string(studentDetails))
}
