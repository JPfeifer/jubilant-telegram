package main

import (
	"os/user"
	"flag"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
	"io"
)

func main() {
	setnewentry()
}

func setnewentry(){
	url,user,password,project_id,user_id,activity_id,hours,comment:= readconfig()
	var url_time_entries string=url+"/time_entries.xml"
	client := &http.Client{}
	body := []byte("<time_entry>\n <project_id>"+project_id+"</project_id>\n <user_id>"+user_id+"</user_id>\n  <activity_id>"+activity_id+"</activity_id>\n <hours>"+hours+"</hours>\n <comments>"+comment+"</comments>\n </time_entry>")
	req, _ := http.NewRequest("POST",url_time_entries, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/xml")
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
		fmt.Println(string(resp.Status))
	if !strings.HasPrefix(resp.Status, "201"){
		resp_body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(resp_body))
		if strings.HasPrefix(resp.Status, "422"){
			if strings.Contains(string(resp_body), "Projekt muss ausgefüllt werden"){
				fmt.Println("Missing issue or project id")
				getissues(url,user,password)
				getprojects(url,user,password)
			}
			if strings.Contains(string(resp_body), "Aktivität muss ausgefüllt werden"){
				fmt.Println("Missing issue or project id")
				gettracker(url,user,password)
			}
		}
						}

}

func readconfig()(string,string,string,string,string,string,string,string){
	usr, err := user.Current()
	var url string
	var user string
	var password string
	var project_id string
	var user_id string
	var activity_id string
	var hours string
	var comment string
	flag.StringVar(&url,"url", "http://exmaple.com/redmine" , "Server URL")
	flag.StringVar(&user, "user", "" , "login user")
	flag.StringVar(&password,"password", "" , "login password")
	flag.StringVar(&project_id,"project_id", "" , "redmine project_id or issue_id")
	flag.StringVar(&user_id, "user_id", "" , "redmine user_id")
	flag.StringVar(&activity_id,"activity_id", "" , "redmine activity_id")
	flag.StringVar(&hours,"hours", "" , "hours spend on tasks")
	flag.StringVar(&comment,"comment", "" , "comment for redmine")
	flag.Parse()
	var uri string = usr.HomeDir+"/.redmine"
	content, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Println("Errored when reading config")
			}
	lines := strings.Split(string(content), "\n")
	for  _,v := range lines{
		if (strings.HasPrefix(v,"url")) && (strings.EqualFold(url,"http://exmaple.com/redmine")) {
				url=v[4:len(v)]
				fmt.Println(url)
		}else if strings.HasPrefix(v, "user=") && (strings.EqualFold(user,"")) {
				user=v[5:len(v)]
		}else if (strings.HasPrefix(v, "password=")) && (strings.EqualFold(password,"")) {
				password=v[9:len(v)]
		}else if (strings.HasPrefix(v, "project_id=")) && (strings.EqualFold(project_id,"")) {
				project_id=v[11:len(v)]
		}else if (strings.HasPrefix(v,"user_id=")) && (strings.EqualFold(user_id,"")) {
				user_id=v[8:len(v)]
		}else if (strings.HasPrefix(v, "activity_id=")) && (strings.EqualFold(activity_id,"")) {
				activity_id=v[12:len(v)]
				}
		}
					return url,user,password,project_id,user_id,activity_id,hours,comment
}

func getissues(url,user,password string){
	url=url+"/issues.json"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)
	type Message struct {
	Issues []struct {
		Author struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"author"`
		CreatedOn   string `json:"created_on"`
		Description string `json:"description"`
		DoneRatio   int    `json:"done_ratio"`
		ID          int    `json:"id"`
		Priority    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"priority"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
			StartDate string `json:"start_date"`
		Status    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		Subject string `json:"subject"`
		Tracker struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tracker"`
		UpdatedOn string `json:"updated_on"`
		} `json:"issues"`
		Limit      int `json:"limit"`
		Offset     int `json:"offset"`
		TotalCount int `json:"total_count"`
			}
	var m Message
	dec := json.NewDecoder(strings.NewReader(string(resp_body)))
	for{
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("DECODE ERROR %v",err)
		}
		}
fmt.Println("available issue_id are:")
		for j := 0; j < m.TotalCount; j++ {
      			fmt.Printf("ID:%v ",m.Issues[j].ID)
        		fmt.Printf("Name:%v \n", m.Issues[j].Subject)
    				}
}

func getprojects(url,user,password string){
	url=url+"/projects.json"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)
	type Message struct {
		Limit    int `json:"limit"`
		Offset   int `json:"offset"`
		Projects []struct {
			CreatedOn   string `json:"created_on"`
			Description string `json:"description"`
			ID          int    `json:"id"`
			Identifier  string `json:"identifier"`
			IsPublic    bool   `json:"is_public"`
			Name        string `json:"name"`
			Status      int    `json:"status"`
			UpdatedOn   string `json:"updated_on"`
		} `json:"projects"`
		TotalCount int `json:"total_count"`
	}
	var m Message
	dec := json.NewDecoder(strings.NewReader(string(resp_body)))
	for{
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("DECODE ERROR %v",err)
		}
		}

		fmt.Println("available Projects are:")
		for j := 0; j < m.TotalCount; j++ {
						fmt.Printf("ID:%v ",m.Projects[j].ID)
						fmt.Printf("Name:%v \n", m.Projects[j].Name)
    				}
	}


	func gettracker(url,user,password string){
		url=url+"/trackers.json"
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(user, password)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}
		defer resp.Body.Close()

		resp_body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(resp_body))
		type Message struct {
			Trackers []struct {
		DefaultStatus struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"default_status"`
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"trackers"`
		}
		var m Message
		dec := json.NewDecoder(strings.NewReader(string(resp_body)))
		for{
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("DECODE ERROR %v",err)
			}
			}

			fmt.Println("available activity id are:")
			for j := 0; j < len(m.Trackers); j++ {
							fmt.Printf("ID:%v ",m.Trackers[j].ID)
							fmt.Printf("Name:%v \n", m.Trackers[j].Name)
	    				}
							}
