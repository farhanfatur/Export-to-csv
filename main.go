package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type M map[string]interface{}

var csvHeader = []string{"id", "employee_name", "employee_age", "employee_salary"}

func take(id int, employes <-chan int, result chan<- int) {
	for x := range employes {
		fmt.Println("Employes:", id, "started ", x)
		time.Sleep(time.Second)
		fmt.Println("Employes:", id, "finished ", x)
		result <- x * 2
	}
}
func takeJson(url string) []M {
	var client = &http.Client{Timeout: 10 * time.Second}
	get, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer get.Body.Close()
	read, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Fatal(err)
	}
	var jsonDecode []M
	str := strings.Replace(string(read), "\ufeff", "", -1)
	errr := json.Unmarshal([]byte(str), &jsonDecode)
	if errr != nil {
		log.Fatal(errr)
	}
	return jsonDecode
}

func saveToCsv(rows []M) {
	file, err := os.Create(fmt.Sprintf("%v", time.Now()) + "-employees.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	Writer := csv.NewWriter(file)
	defer Writer.Flush()
	for i, data := range rows {
		if i == 0 {
			Writer.Write(csvHeader)
		} else {
			csvContent := make([]string, 0)
			for _, key := range csvHeader {
				if val, ok := data[key]; ok {
					csvContent = append(csvContent, fmt.Sprintf("%v", val))
				} else {
					csvContent = append(csvContent, "")
				}
			}
			Writer.Write(csvContent)
		}
	}
}

func main() {
	const api = "http://dummy.restapiexample.com/api/v1/employees"
	employes := make(chan int, 100)
	var data = takeJson(api)
	// fmt.Println(len(data))
	// var makeString = make([]string, 0)
	var eachConvert = make([]M, 0)
	for _, each := range data {
		eachConvert = append(eachConvert, each)
	}
	saveToCsv(eachConvert)
	// for i := 1; i <= 3; i++ {
	// 	go take(i, employes, results)
	// }

	// for j := 1; j <= 5; j++ {
	// 	employes <- j
	// }

	// close(employes)

	// for a := 1; a <= 5; a++ {
	// 	<-results
	// }
}
