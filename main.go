package main
//import the libraries
import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"time"
)
//config structure
type Config struct {
	NoOfItems int64 `yaml:"no_of_items"`
}
//main function
func main() {
	//config object
	var cfg Config
	//calling readFile for reading configfile
	readFile(&cfg)
	//simple get function for our client
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ctx := context.TODO()
			// Parse the POST form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
				return
			}

			// number of items
			r.URL.Query() //.PostForm.Get("number_of_items")
			data := r.URL.Query()
			var numberOfRows string
			//if client pass the query parameter
			if len(data) > 0 {
				numberOfRows = cast.ToString(data.Get("number_of_rows"))
				if cast.ToInt64(numberOfRows) <= 0 {
					numberOfRows = cast.ToString(cfg.NoOfItems)
				}

			} else { //if not then read from the config file
				numberOfRows = cast.ToString(cfg.NoOfItems)
			}
			//calling the function
			lenOfProcessItems, err := processItems(ctx, cast.ToInt(numberOfRows))
			if err != nil {
				//failure return
				http.Error(w, fmt.Sprintf("err %v", err), http.StatusForbidden)
			} else {
				//success return
				w.Write([]byte(fmt.Sprintf("Maximum process %v", lenOfProcessItems)))
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusForbidden)
		}
	})
	// for starting port the actual client
	http.ListenAndServe(":8080", nil)
}

func processItems(context context.Context, numberOfRows int) (uint64, error) {
	var actualItemCountProcess uint64
	var data = make([]Item, numberOfRows)

	//if we need to populate item we can add item details this is comments due to service code is unavailable
	//for i := 0; i < numberOfRows; i++ {
	//	var item Item
	//	item.id = i+1 // if any
	//	item.Name = fmt.Sprintf("item %s", i+1) // if any
	//	data = append(data, item)
	//}

	var service Service
	//n certain number of item , timeInterval given time interval expected.
	noOfItems, timeInterval := service.GetLimits()
	// service is already busy
	if noOfItems == 0 {
		return 0, nil
	}
	//first time to process
	actualItemCountProcess = noOfItems
	//per process in seconds
	perProcess := cast.ToUint64(timeInterval.Seconds()) / cast.ToUint64(noOfItems)

	// to maintain batch number
	counter := uint64(0)
	//for making batch to process
	var b Batch
	//iterate data
	for _, d := range data {
		// when counter is equal to n(number of items process at a time) we will make batch and send it to process
		if noOfItems >= counter {
			b = append(b, d)
			counter = counter + 1
			continue
		} else {
			startTime := time.Now()
			err := service.Process(context, b)
			if err != nil {
				//if somehow Error get we need to minus n value from actualItemsCounter
				actualItemCountProcess = actualItemCountProcess -noOfItems
				return actualItemCountProcess,err
			}
			endTime := time.Now()
			// time difference for processing number of items in real time
			difference := endTime.Sub(startTime).Seconds() // 50 seconds
			// timeInterval (time interval) seconds minus actual seconds use by process function
			remainingTime := timeInterval.Seconds() - difference
			if remainingTime > 0 {
				//counter should be zero for making new batch
				counter = 0
				// new number of batch remaining time / perProcess
				noOfItems = cast.ToUint64(remainingTime) / cast.ToUint64(perProcess)
				//how many items process
				actualItemCountProcess = actualItemCountProcess + noOfItems
				//making new batch for new processing
				b = make(Batch, 0)
			} else {
				break
			}
			
		}

	}
	return actualItemCountProcess, nil
}
//process Error if we can't read config yaml file
func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

//reading the config file
func readFile(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()
	// decode the yaml file through library
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}
