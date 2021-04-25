//Package main implements the Home Baker Web Application.
//		Home Baker Web App is a simple web application that allows user clients to:
//		1. View Sales Items
//		2. Search for sales items
//		3. Create order bookings
//		4. View an edit existing bookings
//		5. Delete existing orders (admin feature)
//		6. Manage orders overview (admin feature)
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	ds "HomeBakerAppGIA2/datastruct"
	feat "HomeBakerAppGIA2/feature"
	ses "HomeBakerAppGIA2/session"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	tpl   *template.Template
	mutex sync.Mutex
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	preLoadUser()                           //preload users
	preLoadOrders()                         // this is done to load some existing orders into the program
	feat.UpdateWeeklySchedule(ds.OrderList) //update weekly schedule list
	feat.UpdateWeeklyOrder(ds.OrderList)    //update weekly item list

	//below codes are for initializing third party logrus
	var filename string = "log/logfile.log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	logChecksum, err := ioutil.ReadFile("log/checksum")
	if err != nil {
		fmt.Println(err)
	}

	str := string(logChecksum)

	if b, err := ComputeSHA256("log/logfile.log"); err != nil {
		fmt.Printf("Err: %v", err)
	} else {
		hash := hex.EncodeToString(b)
		if str == hash {
			fmt.Println("Log integrity OK.")
		} else {
			fmt.Println("File Tampering detected.")
		}
	}

	Formatter := new(log.TextFormatter)
	log.SetOutput(io.MultiWriter(file, os.Stdout)) //default logger will be writing to file and os.Stdout
	log.SetLevel(log.WarnLevel)                    //only log the warning severity level or higher
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true

}

//This function computes the hash (SHA256 value) and returns the hash sum based on contents of the file
func ComputeSHA256(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func main() {
	r := mux.NewRouter()
	r.Handle("/favicon.ico", http.NotFoundHandler()).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/", ses.Index).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/signup", ses.Signup).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/login", ses.Login).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/logout", ses.Logout).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/menu", menu).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/display", feat.Display).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/search", feat.SearchMenu).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/viewPrice", feat.ViewPrice).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/viewCategory", feat.ViewCategory).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/recommended", feat.ViewRecommended).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/createOrder", feat.CreateNewOrder).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/deleteOrder", feat.DelOrder).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/viewOrEdit", feat.ViewOrEditOrder).Methods("POST", "GET").Schemes("https")
	r.HandleFunc("/overview", feat.ViewAllOrders).Methods("POST", "GET").Schemes("https")
	err := http.ListenAndServeTLS(":5221", "./cert.pem", "./key.pem", r)
	if err != nil {
		log.Fatal("HTTPS Server error: ", err)
	}
	//http.ListenAndServe(":5221", r)
}

func menu(res http.ResponseWriter, req *http.Request) {
	myUser := ses.GetUser(res, req)
	tpl.ExecuteTemplate(res, "menu.gohtml", myUser)
}
