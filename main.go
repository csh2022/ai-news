package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type News struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Image       string `json:"image"`
	ArticleLink string `json:"articleLink"`
}

var newsData = []News{
	{
		ID:          1,
		Category:    "AI Simulation",
		Title:       "AI, Simulation, And The Generative Adversarial Network",
		Summary:     "AI enhances simulations by analyzing rich data, using tools like GANs, VAEs, and digital twins to significantly improve processes in manufacturing, healthcare, and synthetic data generation.",
		Image:       "https://images.unsplash.com/photo-1620712943543-bcc极光4688e7485?ixlib=rb-4.0.3&auto=format&fit=crop&w=1350&q=80",
		ArticleLink: "https://www.baidu.com",
	},
	{
		ID:          2,
		Category:    "AI Talent",
		Title:       "Behind the AI talent war: Why tech giants are paying millions to top hires",
		Summary:     "The AI arms race is intensifying, with tech giants like Meta offering massive signing bonuses to secure top AI talent from competitors.",
		Image:       "https://images.unsplash.com/photo-1576091160550-2173dba999ef?ixlib=rb-4.0.3&auto=format&fit=crop&w=1350&q=80",
		ArticleLink: "#",
	},
	{
		ID:          3,
		Category:    "Agentic AI",
		Title:       "The Rise of Agentic AI: How Businesses Are Boosting Productivity",
		Summary:     "Businesses are moving beyond basic automation and generative AI, adopting 'Agentic AI'—autonomous systems that can plan, execute, and complete multi-step tasks.",
		Image:       "https极光://images.unsplash.com/photo-1677442136606-a10723f9极光e1c3?ixlib=rb-4.0.3&auto=format&fit=crop&w=1350&q=80",
		ArticleLink: "#",
	},
	{
		ID:          4,
		Category:    "Optical Computing",
		Title:       "Microsoft doing light work with Analog Optical Computer prototype",
		Summary:     "Microsoft researchers have unveiled a new prototype Analog Optical Computer (AOC) that uses light instead of electricity.",
		Image:       "https://images.unsplash.com/photo-1576675192299-0094e1b83cb5?ixlib=rb-4.0.3&auto=format&fit=crop&w=1350&q=80",
		ArticleLink: "#",
	},
}

func getNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(newsData)
}

func getNewsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}
	
	for _, news := range newsData {
		if news.ID == id {
			json.NewEncoder(w).Encode(news)
			return
		}
	}
	
	http.Error(w, "News not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/api/news", getNews).Methods("GET")
	r.HandleFunc("/api/news/{id}", getNewsByID).Methods("GET")
	
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}