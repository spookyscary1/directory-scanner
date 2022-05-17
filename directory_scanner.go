package main
 
import (
    "flag"
    "fmt"
	 "net/http"
	"os"
	"log"
	"bufio"
	"strings"
)
 
func main() {
    // variables declaration  
    var host string    
    var wordlist string      
 
    // flags declaration using flag package
    flag.StringVar(&host, "h","", "Specify a host to scan.")
    flag.StringVar(&wordlist, "w","", "Specicify a wordlist to use")

	// flag.Var(&host, "h", "Specify a host to scan.")
    // flag.Var(&wordlist, "w", "Specicify a wordlist to use")


    flag.Parse()  // after declaring flags we need to call it
	if(!(strings.HasSuffix(host,"/"))){
		host=host+"/"
	}
	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){

		resp, err := http.Get(host +scanner.Text())
		if err != nil{
			log.Fatal(err)
		}
		if resp.StatusCode==200 ||resp.StatusCode ==301||resp.StatusCode==302{
			fmt.Println(host +scanner.Text())
		}
		
		// fmt.Println(host +scanner.Text())
	}

	if err := scanner.Err(); err !=nil {
		log.Fatal(err)
	}
 

}