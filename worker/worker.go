package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"log"
	"bytes"
	"io/ioutil"
	"io"
	"time"
	"strconv"
	"fmt"
	"strings"
	"github.com/BurntSushi/toml"
)

type Worker struct{
	Id int		`json:"id",omitempty`
	Name string `json:"name"`
	Ip int		`json:"ip"`
}

type Job struct{
	Id  int			`json:"id",omitempty`
	StartFrame int	`json:"startFrame"`
	EndFrame int	`json:"endFrame"`
	TaskId int		`json:"taskId"`
	JobFile string	`json:"jobFile"`
	Done bool		`json:"done"`
}


type Config struct {
	Bcr_server string
	Output_dir string
	Blender_command string
}



var worker Worker
var config Config

func main(){

	readConfig()

	registerWorker()

	for true {
		got, job := checkForJob()
		if got == true {
			do(job)
		}
		time.Sleep(5*time.Second)
	}

}

func readConfig(){
	configfile := "config.ini"

	_, err := os.Stat(configfile)

	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

}

func do(job Job){
	println("doing job...")
	downloadJobFile(job)
	//blender -noaudio --background -o /home/oliver/render/frame_### -s 0 -e 4 -a -E CYCLES -F PNG test.blend
	cmd := exec.Command(config.Blender_command, 	"-noaudio",
									"--background",
									"-o", config.Output_dir+strconv.Itoa(job.Id)+"/frame_####",
									"-s", strconv.Itoa(job.StartFrame),
	 							  	"-e", strconv.Itoa(job.EndFrame),
								  	"-a",
									"-E", "CYCLES",
									"-F", "PNG",
							   	  	"job.blend")
	println("running: "+strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", string(output))
	if err != nil{
		log.Fatal(err)
	}
	job.Done = true
	reportDone(job)
}

func reportDone(job Job){
	println("reporting job done")

	payload, err := json.Marshal(job)

	if err != nil {
		println(err)
	}

	req, err := http.NewRequest("PUT",
		config.Bcr_server+"/worker/"+strconv.Itoa(worker.Id)+"/job/"+strconv.Itoa(job.Id), bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		println("Error reported")
	}

	println("reported")

}

func downloadJobFile(job Job){
	out, err := os.Create("job.blend")
	defer out.Close()

	if err != nil{
		println("error could not create file")
	}

	resp, err := http.Get(job.JobFile)

	if err != nil{
		println("could not download file")
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
}

func checkForJob() (bool, Job){
	println("Checking for job")
	req, err := http.NewRequest("GET",
		config.Bcr_server+"/worker/"+strconv.Itoa(worker.Id)+"/job", nil)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		println("no job.")
		return false, Job{}
	}

	println("Got job!")

	body, _ := ioutil.ReadAll(resp.Body)


	job := Job{}

	json.Unmarshal(body, &job)

	return true, job

}

func registerWorker(){

	hostname,  _ := os.Hostname()
	worker = Worker{Name: hostname, Ip: 4544534} //TODO real ip

	payload, err := json.Marshal(worker)

	if err != nil {
		println(err)
	}


	req, err := http.NewRequest("POST", config.Bcr_server+"/worker", bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		log.Fatal("Server ERROR")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &worker)

}
