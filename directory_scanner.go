package main
 
import (
    "flag"
    "fmt"
	 "net/http"
	"os"
	"log"
	"bufio"
	"strings"
	"crypto/tls"
)
 
func main() {
    // variables declaration  
    var host string    
    var wordlist string      
 
    // flags declaration using flag package
    flag.StringVar(&host, "h","", "Specify a host to scan.")
    flag.StringVar(&wordlist, "w","", "Specicify a wordlist to use")


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
// creating a client that skips verifying tls certs and does not follow redirects
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}	
	client := &http.Client{
		Transport: tr, 
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
}

	for scanner.Scan(){

		resp, err := client.Get(host +scanner.Text())
		if err != nil{
			log.Fatal(err)
		}
		if resp.StatusCode==200 ||resp.StatusCode ==301||resp.StatusCode==302{
			fmt.Print(host +scanner.Text()+" ")
			fmt.Println(resp.StatusCode)
		}
		
	}

	if err := scanner.Err(); err !=nil {
		log.Fatal(err)
	}
 

}