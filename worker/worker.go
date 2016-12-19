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
	"path/filepath"
	"archive/zip"
	"github.com/BurntSushi/toml"
	"github.com/xyproto/unzip"
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
	JobFileType string `json:"jobFileType"`
	Done bool		`json:"done"`
}


type Config struct {
	Bcr_server string
	Output_dir string
	Blender_command string
	Working_dir string
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
	cleanWorkingDir()
	downloadJobFile(job)
	//blender -noaudio --background -o /home/oliver/render/frame_### -s 0 -e 4 -a -E CYCLES -F PNG test.blend

	// The order of options seems to be crucial here
	cmd := exec.Command(config.Blender_command, 	"-noaudio",
									"--background",
									config.Working_dir+"job.blend",
									"-E", "CYCLES",
									"-E", "CYCLES",
									"-o", config.Output_dir+strconv.Itoa(job.TaskId)+"/frame_####",
									"-s", strconv.Itoa(job.StartFrame),
	 								"-e", strconv.Itoa(job.EndFrame),
									"-a")
	println("running: "+strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", string(output))
	if err != nil{
		log.Fatal(err)
	}
	job.Done = true
	reportDone(job)
}

func cleanWorkingDir(){
	println("cleaning working dir...")
	os.RemoveAll(config.Working_dir)
	println("cleaned working dir.")
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
	var filename string
	if job.JobFileType == ".zip"{
		filename = "job.zip"
	}else{
		filename = "job.blend"
	}
	os.Mkdir(config.Working_dir, 0755)
	out, err := os.Create(config.Working_dir+filename)
	defer out.Close()

	if err != nil{
		log.Fatal(err)

		println("error could not create file")
	}

	resp, err := http.Get(job.JobFile)

	if err != nil{
		log.Fatal(err)
		println("could not download file")
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)

	if job.JobFileType == ".zip"{
		println("unzipping...")
		//Unzip(config.Working_dir+filename, config.Working_dir)
		unzip.Extract(config.Working_dir+filename, config.Working_dir)
		println("unzipped.")

		files, _ := ioutil.ReadDir(config.Working_dir)

		println(files)

		found := false

		for _, f := range files {
			println("checking: "+f.Name())
			if filepath.Ext(config.Working_dir+f.Name()) == ".blend" {
				println("found .blend")
				os.Rename(config.Working_dir+f.Name(), config.Working_dir+"job.blend")
				found = true
				break
			}
		}

		if !found {
			log.Fatal("no jobfile found.")
		}


	}

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


func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
