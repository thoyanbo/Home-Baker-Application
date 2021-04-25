package feature

import (
	"net/http"
	"regexp"
	"strings"

	ds "HomeBakerAppGIA2/datastruct"

	log "github.com/sirupsen/logrus"
)

//Http response to show search menu and search result if search performed matches a an existing sales item. Search is case-sensitive and requires exact match.
func SearchMenu(res http.ResponseWriter, req *http.Request) {
	var item string
	var itemFound bool = false
	var i int
	var foundItem ds.ItemInfo

	if req.Method == http.MethodPost {
		item = strings.TrimSpace(req.FormValue("name"))
		itemRegExp := regexp.MustCompile(`^[\w'\-,.][^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{2,30}$`) //name regexp to check for name pattern match
		if !itemRegExp.MatchString(item) {
			http.Error(res, "You have entered an invalid name field.", http.StatusBadRequest)
			log.Warning("Invalid user input for name field")
			return
		}

		item = pol.Sanitize(item)

		for i = 0; i < len(ds.Items); i++ {
			if item == ds.Items[i].Name {
				itemFound = true
				foundItem = ds.Items[i]
				break
			}
		}
	}

	err := tpl.ExecuteTemplate(res, "search.gohtml", foundItem) //displays search result if item found, else found item is empty struct and will not display
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at search template.", err)
		return
	}

	if itemFound == false {
		err := tpl.ExecuteTemplate(res, "itemNotFound.gohtml", item) //displays error message if item cannot be found
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Error("Error at itemNotFound template.", err)
			return
		}
	}
}

//Http response to show results of chosen form category.
func ViewCategory(res http.ResponseWriter, req *http.Request) {
	var data []ds.ItemInfo
	if req.Method == http.MethodPost {
		cat := req.FormValue("category")
		for i := 0; i < len(ds.Items); i++ {
			if cat == ds.Items[i].Category {
				data = append(data, ds.Items[i])
				//fmt.Println(data)
			}
		}
	}
	err := tpl.ExecuteTemplate(res, "viewCategory.gohtml", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at viewCategory template.", err)
		return
	}
}

//Http response to show recommended items
func ViewRecommended(res http.ResponseWriter, req *http.Request) {
	var data []ds.ItemInfo
	for i := 0; i < len(ds.Items); i++ {
		if ds.Items[i].Recommended == true {
			data = append(data, ds.Items[i])
		}
	}

	err := tpl.ExecuteTemplate(res, "viewRecommended.gohtml", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at viewRecommended template.", err)
		return
	}
}

//http response to show all sales items sorted by Price
func ViewPrice(res http.ResponseWriter, req *http.Request) {
	sortPrice(ds.Items)
	err := tpl.ExecuteTemplate(res, "viewPrice.gohtml", ds.Items)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error("Error at viewPrice template.", err)
		return
	}
}

// This function does a selection sort of the array by price.
func sortPrice(ii []ds.ItemInfo) {
	var n int
	n = len(ds.Items)
	mutex.Lock() // mutex lock to allow sorting to be done on global variable item array at any single time.
	{
		for last := n - 1; last >= 1; last-- {
			// select most expensive item in array
			largest := indexOfLargest(ii, last+1)

			//swap largest item array[largest] with array[last]
			swap(&ii[largest], &ii[last])
		}
	}
	mutex.Unlock()
}

//Function returns index of largest item through numberical comparison
func indexOfLargest(ii []ds.ItemInfo, n int) int {
	largestIndex := 0
	for i := 1; i < n; i++ {
		if ii[i].UnitPrice > ii[largestIndex].UnitPrice {
			largestIndex = i
		}
	}
	return largestIndex
}

func swap(x *ds.ItemInfo, y *ds.ItemInfo) {
	temp := *x
	*x = *y
	*y = temp
}
