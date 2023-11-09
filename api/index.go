package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mr-destructive/geneapi/geneapi"
	"github.com/mr-destructive/geneapi/geneapi/types"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//create a tempalte
		templateStr := `<!DOCTYPE html>
        <html>
            <head>
                <title>GeneAI</title>
                <script src="https://unpkg.com/htmx.org@1.9.6"></script>
            </head>
            <body>
                <h1>GeneAI</h1>
                <form hx-post="/" hx-swap="innerHTML" hx-target="#result">
                    <input type="text" name="prompt" placeholder="Prompt">
                    <select name="model">
                        <option value="palm2" selected>PaLM2</option>
                        <option value="cohereai">Cohere</option>
                        <option value="openai">OpenAI</option>
                    </select>
                    <input type="number" name="max-tokens" placeholder="Max Tokens" min="20" max="800">
                    <input type="number" name="temperature" placeholder="Temperature" min="0" max="1">
                    <input type="submit" value="Submit">
                </form>
                <p id="result"></p>
            </body>
        </html>`

		temp := template.New("index")
		t, err := temp.Parse(templateStr)
		t.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else {
		prompt := r.FormValue("prompt")
		//convert to int
		maxTokens, err := strconv.Atoi(r.FormValue("max-tokens"))
		if err != nil {
			panic(err)
		}
		temperature, err := strconv.ParseFloat(r.FormValue("temperature"), 64)
		model := r.FormValue("model")
		apiKeys := map[string]string{
			"openai":   os.Getenv("openai"),
			"palm2":    os.Getenv("palm2"),
			"cohereai": os.Getenv("cohereai"),
		}
		req := &types.Request{
			Prompt:      prompt,
			MaxTokens:   maxTokens,
			Temperature: temperature,
		}
		fmt.Println(req)
		resp, err := geneapi.GeneAI(*req, model, apiKeys)
        fmt.Println(resp)
		if err != nil {
			panic(err)
		}
		tmpl := template.New("result")
		t, err := tmpl.Parse(resp)
		t.Execute(w, nil)
	}
}

func readEnv() (map[string]string, error) {
	//open the .env file
	env, err := os.ReadFile("api/.env")
	if err != nil {
		return nil, err
	}
	// read from the env file as = sepearted values
	apiKeys := make(map[string]string)
	err = fmt.Errorf("failed to parse env vars")
	for _, v := range strings.Split(string(env), "\n") {
		if strings.Contains(v, "=") {
			apiKeys[strings.Split(v, "=")[0]] = strings.Split(v, "=")[1]
		}
	}
	return apiKeys, nil
}
