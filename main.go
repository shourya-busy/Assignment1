package main

import "fmt"
import "reflect"


//searches through the source for possible occurrence of the key
func SearchKey(key string, source map[string]interface{}) (map[string]interface{},error) {

	//If the key exists in the first level return it
	if _, ok := source[key]; ok {
		return source,nil
	}
	
	//if the key is not found on above level iterate through the MAP values
	for _,val := range source {
		valReflected := reflect.ValueOf(val)
		
		//Case to match the kind of reflected value inside MAP
		//Only checking for two data structures that could possibly hold a nested MAP
	    switch valReflected.Kind() {
        
		//recursively search through the nested MAP
		case reflect.Map :
			nestedMap := valReflected.Interface().(map[string]interface{})
			
			//if the key is located return the entire MAP where it was found
			if found, err := SearchKey(key,nestedMap); err == nil {
				return found,nil
			}
        //Iterate through the slice to check if any element is a MAP
		case reflect.Slice :
			nestedSlice := valReflected.Interface().([]interface{})

			for _, value := range nestedSlice {
				
				valueReflected := reflect.ValueOf(value)
				

				if valueReflected.Kind() == reflect.Map {
					nestedMapInSlice := valueReflected.Interface().(map[string]interface{})
					
					//if the key is located return the entire MAP where it was found
					if found, err := SearchKey(key,nestedMapInSlice); err == nil {
						return found,nil
					}
				}

			}

		}

	}
	//if no value is found return an error
	return nil, fmt.Errorf("Key Not Found")

}

//Sets the Value for a Key in source
func SetKeyValue(key string,value interface{},source map[string]interface{}){
	//uses the search function to locate the MAP where key exists
	if foundMap,err := SearchKey(key,source); err != nil {
		fmt.Println(err)
	} else {
		//Since maps are sent through reference the value can be updated
		foundMap[key] = value
		fmt.Println("Key Updated")
		fmt.Printf("%v : %v\n\n",key,foundMap[key])
	}
}

//Removes the key from the source
func RemoveKey(key string, source map[string]interface{}) {
	//uses the search function to locate the MAP where key exists
	if foundMap,err := SearchKey(key,source); err != nil {
		fmt.Println(err)
	} else {
		//delete function deleted the key from the MAP 
		delete(foundMap,key)
		fmt.Printf("Key Deleted : %v\n\n",key)
	}
}

func main() {
	//Input data structure
	var source  = map[string]interface{}{
		"name" : "Tolexo Online Pvt. Ltd",
		"age_in_years" : 8.5,
		"origin" : "Noida",
		"head_office" : "Noida, Uttar Pradesh",
		"address" : []interface{}{
			map[string]interface{}{
				"street" : "91 Springboard",
				"landmark" : "Axis Bank",
				"city" : "Noida",
				"pincode" : 201301,
				"state" : "Uttar Pradesh",
				},
			map[string]interface{}{
				"street" : "91 Springboard",
				"landmark" : "Axis Bank",
				"city" : "Noida",
				"pincode" : 201301,
				"state" : "Uttar Pradesh",
				},
		},
		"sponsers" : map[string]interface{}{
			"name" : "One",
		},
		"revenue" : "19.8 million$",
		"no_of_employee" : 630,
		"str_text" : []interface{}{"one","two"},
		"int_text" : []interface{}{1,3,4},
	}


	//sample value to update a particular key
	newMap := map[string]interface{}{
		"status" : "complete",
		"array" : []interface{}{1.0,3.4},
	}

	//Function to set a value to a particular key
	SetKeyValue("pincode",newMap,source)

	//Function to remove a particular key
	RemoveKey("status",source)

}