package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

/**
 ************ Structures definitions
 */
// Structure for DryMeasures
type DryMeasure struct {
	Id          string `json:"id"`
	Cups        int    `json:"cups"`
	Tablespoons int    `json:"tablespoons"`
	Teaspoons   int    `json:"teaspoons"`
	Grams       int    `json:"grams"`
}

// Structure for LiquidMeasures
type LiquidMeasure struct {
	Gallons int `json:"gallons"`
	Quarts  int `json:"quarts"`
	Pints   int `json:"pints"`
	Cups    int `json:"cups"`
	FluidOz int `json:"fluidOz"`
	//todo: add also ml
}

// -----------------------------------------------------------------

type CookingConvertor struct {
	quantity float32
	fromUnit Unit
	toUnit   Unit
}

// metric from https://www.thecalculatorsite.com/cooking/cooking-calculator.php
type Unit string

const (
	Cups        Unit = "cups"
	Gallons          = "gal"
	Grams            = "g"
	Kilograms        = "kg"
	Liters           = "l"
	Milliliters      = "ml"
	Deciliters       = "dl"
	Ounces           = "oz"
	Pounds           = "lb"
	Pints            = "pt"
	Tablespoons      = "tbsp"
	Teaspoons        = "tsp"
)

//todo try to send the cookingConvertor as pointer *CookingConvertor, otherwise it's a copy
func (converter *CookingConvertor) Convert() (*CookingConvertor, error) {
	fmt.Printf("Convert the value of %f \n", converter.quantity)

	if err := converter.toUnit.IsValid(); err != nil {
		return nil, err
	}
	if err := converter.fromUnit.IsValid(); err != nil {
		return nil, err
	}

	return converter, nil
}

func (unit Unit) IsValid() error {
	switch unit {
	case Cups, Gallons, Grams, Kilograms, Liters, Tablespoons, Milliliters, Teaspoons:
		return nil
	}
	return errors.New("invalid conversion unit type")
}

func (converter *CookingConvertor) From(fromVal Unit) *CookingConvertor {
	fmt.Printf("from %s \n", fromVal)
	converter.fromUnit = fromVal
	return converter
}

func (converter *CookingConvertor) To(toVal Unit) *CookingConvertor {
	fmt.Printf("to %s \n", toVal)
	converter.toUnit = toVal
	return converter
}

//todo in the end should look like convert(5).from(Cups).to(Kg)

/**
************
 */

// declare a global DryMeasures array
// that we can then populate in our main function
// to simulate a database
var DryMeasures []DryMeasure

func getAllDryMeasures(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllDryMeasures")
	json.NewEncoder(w).Encode(DryMeasures)
}

func getAllDryMeasuresById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllDryMeasuresById")
	vars := mux.Vars(r)
	pathParamId := vars["id"]

	// Loop over all of our DryMeasures
	// if the dryMeasure.Id equals the key we pass in
	// return the article encoded as JSON
	for _, dryMeasure := range DryMeasures {
		if dryMeasure.Id == pathParamId {
			json.NewEncoder(w).Encode(dryMeasure)
		}
	}
}

func createNewDryMeasure(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewDryMeasure")
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)

	var dryMeasure DryMeasure
	json.Unmarshal(reqBody, &dryMeasure)

	// update our global Articles array to include
	// our new Article
	DryMeasures = append(DryMeasures, dryMeasure)

	json.NewEncoder(w).Encode(dryMeasure)
}

func deleteDryMeasureById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteDryMeasureById")

	vars := mux.Vars(r)
	id := vars["id"]

	// we then need to loop through all our DryMeasures
	for index, dryMeasure := range DryMeasures {
		// if our id path parameter matches one of our
		// dryMeasure
		if dryMeasure.Id == id {
			// updates our DryMeasures array to remove the
			// measure
			DryMeasures = append(DryMeasures[:index], DryMeasures[index+1:]...)
		}
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	var router = mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/dries", getAllDryMeasures)
	router.HandleFunc("/dry", createNewDryMeasure).Methods("POST")
	router.HandleFunc("/dry/{id}", deleteDryMeasureById).Methods("DELETE") //this has to be defined before the other /dry/{id}
	router.HandleFunc("/dry/{id}", getAllDryMeasuresById)                  // just for the sake of example for now

	log.Fatal(http.ListenAndServe(":10000", router))
}

// Start the web server
func main() {
	fmt.Println("***** REST API - Cooking CookingConvertor v1.0 *****")
	DryMeasures = []DryMeasure{
		DryMeasure{Id: "0", Cups: 1, Tablespoons: 1, Teaspoons: 1, Grams: 5},
		DryMeasure{Id: "1", Cups: 1, Tablespoons: 2, Teaspoons: 2, Grams: 10},
	}

	//this could represent the whole lib
	var converter = CookingConvertor{quantity: 1, fromUnit: Cups, toUnit: Tablespoons}
	_, _ = converter.Convert()

	// converter.Convert(858.9).From("g").To("ml")
	//converter.From("tbsp").To("c").Convert(1)

	handleRequests()

}
