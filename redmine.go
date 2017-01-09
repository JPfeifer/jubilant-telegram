package main

import (
	"os/user"
	"flag"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {	
	setnewentry()
}	

func setnewentry(){
	 url,user,password,project_id,user_id,activity_id,hours,comment := readconfig()

	client := &http.Client{}

	body := []byte("<time_entry>\n <project_id>"+project_id+"</project_id>\n <user_id>"+user_id+"</user_id>\n <activity_id>"+activity_id+"</activity_id>\n <hours>"+hours+"</hours>\n <comments>"+comment+"</comments>\n </time_entry>")
	url="http://"+url+"/time_entries.xml"
	req, _ := http.NewRequest("POST",url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/xml")
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)	
	if !strings.HasPrefix(resp.Status, "201"){
		resp_body, _ := ioutil.ReadAll(resp.Body)	
		fmt.Println(string(resp_body))
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
	flag.StringVar(&url,"url", "0.0.0.0" , "Server IP")
	flag.StringVar(&user, "user", "" , "login user")
	flag.StringVar(&password,"password", "" , "login password")
	flag.StringVar(&project_id,"project_id", "" , "redmine project_id")
	flag.StringVar(&user_id, "user_id", "" , "redmine user_id")
	flag.StringVar(&activity_id,"activity_id", "" , "redmine activity_id")
	flag.StringVar(&hours,"hours", "" , "hours spend on tasks")
	flag.StringVar(&comment,"comment", "" , "comment for redmine")
	flag.Parse()
    	if err != nil {
		fmt.Println("Errored when finding User home")
    	}
	var uri string = usr.HomeDir+"/.redmine"
	content, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Println("Errored when reading config")
			}
	lines := strings.Split(string(content), "\n")
	for  _,v := range lines{
		if (strings.HasPrefix(v,"url")) && (strings.EqualFold(url,"0.0.0.0")) {
				url=v[4:len(v)]
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
