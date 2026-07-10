package items

import (
	"log"
	"sync"
)

func ProcessItems(items []Item) {
	var wg sync.WaitGroup
	for _, v := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("%s %s", v.ID, v.Name)
		}()
	}
	wg.Wait()
}
