package main
 
import (
    "flag"
    "fmt"
	"net/http"
	"os"
	"log"
	"bufio"
	"strings"
	"errors"
	"strconv"
	"crypto/tls"
)
 
func main() {
    // variables declaration  
    var host string    
    var wordlist string
	var Statuscodes string
	var codeBehavior string

 
    // flags declaration using flag package
    flag.StringVar(&host, "h","", "Specify a host to scan.")
    flag.StringVar(&wordlist, "w","", "Specicify a wordlist to use")
	flag.StringVar(&Statuscodes,"c","200,302,301","specify a status codes to display")
	flag.StringVar(&codeBehavior, "b", "allow", "specify to allow or deny the listed status codes")


    flag.Parse()  // after declaring flags we need to call it


	// sanity check code behavior variable 
	if !(strings.ToLower(codeBehavior)=="allow" || strings.ToLower(codeBehavior)=="deny" ){
		fmt.Println("Error: b must be set to allow or deny")
		os.Exit(4)
	}
	// parse status codes
	mapCodeList, err :=parseCodeList(Statuscodes)
	if err != nil{
		fmt.Println("Error: list of status codes not valid")
		os.Exit(4)
	}
	
	// append a slash to the host if when does not already exist
	if(!(strings.HasSuffix(host,"/"))){
		host=host+"/"
	}
	// opens wordlist 
	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
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
// opens lines in the wordlist 
	for scanner.Scan(){

		resp, err := client.Get(host +scanner.Text())
		if err != nil{
			log.Fatal(err)
		}

		_, exists := mapCodeList[resp.StatusCode]

		if codeBehavior=="allow" && exists== true{
				printResult(resp.StatusCode, host +scanner.Text())
		} else if (codeBehavior=="deny" && exists==false) {
			printResult(resp.StatusCode, host +scanner.Text())
		} 
		
	}

	if err := scanner.Err(); err !=nil {
		log.Fatal(err)
	}
 

}
func parseCodeList (list string) (map[int]bool, error) {
	array := strings.Split(list,",")
	var codes= make(map[int]bool)
	for i := 0; i < len(array); i++ {
		x, err := strconv.ParseInt(array[i], 10, 64)
		if err != nil {
			return codes, errors.New("invalid comma seperated status code list")
		}
		codes[int(x)]=true
		
	}
	return codes, nil
}
func printResult (aStatusCode int,anUrl string ) {
	fmt.Print(anUrl+" ")
	fmt.Println(aStatusCode)
}