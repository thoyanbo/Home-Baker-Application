package main

import (
	"net/http"
)

func searchMenu(res http.ResponseWriter, req *http.Request) { //http response to show search menu and search result if search is performed and item found
	var item string
	var itemFound bool = false
	var i int
	var foundItem itemInfo

	if req.Method == http.MethodPost {
		item = req.FormValue("name")
		for i = 0; i < len(items); i++ {
			if item == items[i].Name {
				itemFound = true
				foundItem = items[i]
				break
			}
		}
	}

	err := tpl.ExecuteTemplate(res, "search.gohtml", foundItem) //displays search result if item found, else found item is empty struct and will not display
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if itemFound == false {
		err := tpl.ExecuteTemplate(res, "itemNotFound.gohtml", item) //displays error message if item cannot be found
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func viewCategory(res http.ResponseWriter, req *http.Request) { //http response to show result of chosen category
	var data []itemInfo
	if req.Method == http.MethodPost {
		cat := req.FormValue("category")
		for i := 0; i < len(items); i++ {
			if cat == items[i].Category {
				data = append(data, items[i])
				//fmt.Println(data)
			}
		}
	}
	err := tpl.ExecuteTemplate(res, "viewCategory.gohtml", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func viewRecommended(res http.ResponseWriter, req *http.Request) { //http response to show recommended items
	var data []itemInfo
	for i := 0; i < len(items); i++ {
		if items[i].Recommended == true {
			data = append(data, items[i])
		}
	}

	err := tpl.ExecuteTemplate(res, "viewRecommended.gohtml", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func viewPrice(res http.ResponseWriter, req *http.Request) { //http response to show items sorted by price
	sortPrice(items)
	err := tpl.ExecuteTemplate(res, "viewPrice.gohtml", items)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func sortPrice(ii []itemInfo) { // this function does a selection sort of the array by price.
	var n int
	n = len(items)
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

func indexOfLargest(ii []itemInfo, n int) int {
	largestIndex := 0
	for i := 1; i < n; i++ {
		if ii[i].UnitPrice > ii[largestIndex].UnitPrice {
			largestIndex = i
		}
	}
	return largestIndex
}

func swap(x *itemInfo, y *itemInfo) {
	temp := *x
	*x = *y
	*y = temp
}
