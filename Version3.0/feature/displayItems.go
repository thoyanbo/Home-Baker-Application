package feature

import (
	"errors"
	"fmt"
	"net/http"

	ds "HomeBakerAppGIA2/datastruct"

	log "github.com/sirupsen/logrus"
)

//Function Displays all sales item to user as HTML response
func Display(res http.ResponseWriter, req *http.Request) {

	sortItems(ds.Items)
	err := tpl.ExecuteTemplate(res, "display.gohtml", ds.Items)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}

}

//Function sorts the item array by category.
func sortItems(ii []ds.ItemInfo) error {
	defer func() { //to handle potential panic situation
		if err := recover(); err != nil {
			fmt.Println("Oops, panic occurred:", err)
			log.Panic("Panic occured at sort items:", err)
		}
	}()
	var n int
	n = len(ds.Items)
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

// Function returns index of last item in alphabetical order through string comparison
func indexOfLast(ii []ds.ItemInfo, n int) (int, error) {
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
