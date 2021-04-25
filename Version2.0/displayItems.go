package main

import (
	"errors"
	"fmt"
	"net/http"
)

func display(res http.ResponseWriter, req *http.Request) { //function displays all sales item to user

	sortItems(items)
	err := tpl.ExecuteTemplate(res, "display.gohtml", items)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

}

func sortItems(ii []itemInfo) error { //Note that display items by default will sort the items by category and display by category.
	defer func() { //to handle possible panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
		}
	}()
	var n int
	n = len(items)
	mutex.Lock() // mutex lock to allow sorting to be done on global variable item array at any single time.
	{
		for last := n - 1; last >= 1; last-- {
			// select the last alphabetical item in array
			largest, _ := indexOfLast(ii, last+1)

			//swap largest item array[largest] with array[last]
			swap(&ii[largest], &ii[last])
		}
	}
	mutex.Unlock()
	return nil
}

func indexOfLast(ii []itemInfo, n int) (int, error) { // this function returns index of last item in alphabetical order through string comparison
	if len(ii) == 0 {
		return 0, errors.New("the list is empty")
	}
	lastIndex := 0
	for i := 1; i < n; i++ {
		if ii[i].Name > ii[lastIndex].Name {
			lastIndex = i
		}
	}
	return lastIndex, nil
}
