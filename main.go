package main

import "fmt"
import "reflect"

type Company struct {
	Name       string  `json:"name"`
	AgeInYears float64 `json:"age_in_years"`
	Origin     string  `json:"origin"`
	HeadOffice string  `json:"head_office"`
	Address    []struct {
		Street   string `json:"street"`
		Landmark string `json:"landmark"`
		City     string `json:"city"`
		Pincode  int    `json:"pincode"`
		State    string `json:"state"`
	} `json:"address"`
	Sponsers struct {
		Name string `json:"name"`
	} `json:"sponsers"`
	Revenue      string   `json:"revenue"`
	NoOfEmployee int      `json:"no_of_employee"`
	StrText      []string `json:"str_text"`
	IntText      []int    `json:"int_text"`
}


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


//Populates the given struct with data from the source
func PopulateStruct(source map[string]interface{}, sink interface{}) error {

	//Extract the reflected object from the sink
	//Indirect extracts the value if sink is a pointer
	sinkReflected := reflect.Indirect(reflect.ValueOf(sink))

	//extract the type of the reflected object
	sinkReflectedType := sinkReflected.Type()

	//iterate over the fields of Struct sink
	for i := 0 ; i < sinkReflectedType.NumField(); i++ {
		
		//extract the field type
		field := sinkReflectedType.Field(i)

		//extract the value of field
		fieldReflected := sinkReflected.Field(i)

		//extract the json tag string from the filed
		tag := field.Tag.Get("json")

		//check if the field value is settable
		if !fieldReflected.CanSet(){
			return fmt.Errorf("The %v field is not Settable",field.Name)
		}

		//if it is settable search the MAP for the field json tag as the key
		//and store the reference to the MAP where key is located
		foundMap,err := SearchKey(tag,source); 

		if err != nil {
			return err
		}

		//Case to match the kind of field value
		switch fieldReflected.Kind() {
		
		//if the field is a nested struct handle it recursively with the new source as searched MAP
		case reflect.Struct:
			err := PopulateStruct(foundMap[tag].(map[string]interface{}),fieldReflected.Addr().Interface())
			if err != nil {
				return err
			}

		//Check the type of slice elements if it is struct handle recursively
		case reflect.Slice:

			//store the type of slice elements
			fieldType := fieldReflected.Type().Elem()

			//Make a slice of the same type as the field type
			//and length same as the size of searched MAP reference
			slice := reflect.MakeSlice(reflect.SliceOf(fieldType), len(foundMap[tag].([]interface{})), len(foundMap[tag].([]interface{})))
			
			//Iterate over the slice recursively calling the nested Structs
			//and store the values in a slice reference
			if fieldType.Kind() == reflect.Struct{
				for j,val := range foundMap[tag].([]interface{}) {

					err := PopulateStruct(val.(map[string]interface{}),slice.Index(j).Addr().Interface())
					if err != nil {
						return err
					}
						
				}
			//Since the element type is no special data structure simply append the values to the slice 
			} else  {

					for j, v := range foundMap[tag].([]interface{}) {
						slice.Index(j).Set(reflect.ValueOf(v))
					}
			}

			//finally set the field to the created slice	
			fieldReflected.Set(slice)

		default:
			//If the field value is no special data structure and the value found from the MAP is assignable to the field
			//simply set the field with the found value 
			if reflect.TypeOf(foundMap[tag]).AssignableTo(fieldReflected.Type()) {
				fieldReflected.Set(reflect.ValueOf(foundMap[tag]))
			} 
		}

	}

	//if no problem is encountered return nil
    return nil
}


//Prints the struct
func printStruct(data interface{}) {

	//extract the value from data
	//Indirect fetches the value if the data is a pointer
    val := reflect.Indirect(reflect.ValueOf(data))

    fmt.Println(val.Type().Name())

	//Iterate over the struct fields
    for i := 0; i < val.NumField(); i++ {
		//Extract the field from reflected data
        field := val.Field(i)

		//Print the field name
        fmt.Printf("  - %s: ", val.Type().Field(i).Name)

		//Cases to match the kind of field
        switch field.Kind() {
		
		//Handle the structs recursively
        case reflect.Struct:
            fmt.Print("{")
            printStruct(field.Interface())
			fmt.Println("}")
		
        case reflect.Slice:

			//If the element type of slice is Struct
			//Iterate over each element recursively calling it
			if field.Type().Elem().Kind() == reflect.Struct{
			
			fmt.Print("[")
            for j := 0; j < field.Len(); j++ {
                fmt.Println("\n{")
                printStruct(field.Index(j).Interface())
				fmt.Println("},")
            }

			fmt.Println("]")
			} else {  
				//if the slice has no special data structures simply print it
				fmt.Println(field.Interface())
			}
        default:
			//if the field has no special data structure simply print it
            fmt.Println(field.Interface())
        }
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

	//variable of struct Company to Unmarshal the source
	var company Company

	//Call to populate the company reference with data from the source
	err := PopulateStruct(source,&company)
	if err != nil {
		fmt.Println(err)
	}

	//Call to print the Company Struct
	printStruct(company)


}