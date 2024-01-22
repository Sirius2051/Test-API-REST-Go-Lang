package main

import (
    "encoding/json"
    "log"
    "net/http"
    "fmt"
)

type Joke struct {
    ID string `json:"id"`
	URL string `json:"url"`
}

const url = "https://api.chucknorris.io/jokes/random"

func getJoke(jokes map[string]bool) (Joke, error) {
    var joke Joke
    for {
        
		response, err := http.Get(url)
        
		if err != nil {
            
			return joke, err
        }
        
		defer response.Body.Close()

        err = json.NewDecoder(response.Body).Decode(&joke)
        
		if err != nil {
            
			return joke, err
        }

        if _, exists := jokes[joke.ID]; !exists {
            
			jokes[joke.ID] = true
            break
        }
    }

    return joke, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		
			return
		}

        jokes := make(map[string]bool)
        
		jokesList := make([]Joke, 0, 25)

		for len(jokes) < 25 {
		
			joke, err := getJoke(jokes)
			
			if err != nil {
			
				log.Fatal(err)
			}
	
			jokesList = append(jokesList, joke)
		}

        w.Header().Set("Content-Type", "application/json")
        
		json.NewEncoder(w).Encode(jokesList)
	} ) 
	
	svr := http.Server{
		Addr: ":8080",
	}

	err := svr.ListenAndServe()

	if err != nil {
		
		panic(err)
	}
}