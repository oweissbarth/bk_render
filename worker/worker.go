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
}

const BCR_SERVER = "http://localhost:5000"
const OUTPUT_DIR = "/home/oliver/render/"
const BLENDER_COMMAND = "blender"

var worker Worker

func main(){

	registerWorker()

	for true {
		got, job := checkForJob()
		if got == true {
			do(job)
		}
		time.Sleep(5*time.Second)
	}

}

func do(job Job){
	println("doing job...")
	downloadJobFile(job)
	//blender -noaudio --background -o /home/oliver/render/frame_### -s 0 -e 4 -a -E CYCLES -F PNG test.blend
	cmd := exec.Command("blender", 	"-noaudio",
									"--background",
									"-o", OUTPUT_DIR+strconv.Itoa(job.Id)+"/frame_####",
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
	reportDone(job)
}

func reportDone(job Job){

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
		BCR_SERVER+"/worker/"+strconv.Itoa(worker.Id)+"/job", nil)

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


	req, err := http.NewRequest("POST", BCR_SERVER+"/worker", bytes.NewBuffer(payload))

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
