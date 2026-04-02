package main // This denote that this is a standalone progarm and is not a package

// In the program go.mod explicitly define where is the entry point similar to pyproject.taml
import (
	"fmt"
	"log"
	"net/http"
)
var prefix="/api/v1"

func main() { // This define global entry point for program
	router := http.NewServeMux()
	router.HandleFunc(fmt.Sprintf("GET %s/api", prefix), func (w http.ResponseWriter, r *http.Request) {
		
		w.Write([]byte("Hello World"))
	})
	log.Println("Initiating Server")
	server := http.Server {
		Addr : ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}	

