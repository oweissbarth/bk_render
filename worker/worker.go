package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	log "github.com/Sirupsen/logrus"
	"bytes"
	"io/ioutil"
	"io"
	"time"
	"strconv"
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

type Chunk struct{
	Id  int				`json:"id",omitempty`
	StartFrame int		`json:"startFrame"`
	EndFrame int		`json:"endFrame"`
	JobId int			`json:"jobId"`
	JobFile string		`json:"jobFile"`
	JobFileType string 	`json:"jobFileType"`
	Done bool			`json:"done"`
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

	log.SetLevel(log.InfoLevel)


	readConfig()

	registerWorker()

	for true {
		got, job := checkForChunk()
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
		log.Error("Config file is missing: ", configfile)
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}

func do(job Chunk){
	println("doing job...")
	cleanWorkingDir()
	downloadChunkFile(job)
	//blender -noaudio --background -o /home/oliver/render/frame_### -s 0 -e 4 -a -E CYCLES -F PNG test.blend

	// The order of options seems to be crucial here
	cmd := exec.Command(config.Blender_command, 	"-noaudio",
									"--background",
									config.Working_dir+"job.blend",
									"-E", "CYCLES",
									"-F", "PNG",
									"-o", config.Output_dir+strconv.Itoa(job.JobId)+"/frame_####",
									"-s", strconv.Itoa(job.StartFrame),
	 								"-e", strconv.Itoa(job.EndFrame),
									"-a")
	log.Info("running: "+strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()

	log.Debug(output)

	if err != nil{
		log.Error(err)
		os.Exit(1)
	}
	job.Done = true
	reportDone(job)
}

func cleanWorkingDir(){
	log.Info("cleaning working dir...")
	os.RemoveAll(config.Working_dir)
	log.Info("cleaned working dir.")
}

func reportDone(job Chunk){
	log.Info("reporting job done")

	payload, err := json.Marshal(job)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	req, err := http.NewRequest("PUT",
		config.Bcr_server+"/worker/"+strconv.Itoa(worker.Id)+"/job/"+strconv.Itoa(job.Id), bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error("Marking job as done failed. Got", resp.StatusCode)
	}

	log.Info("Reported done.")
}

func downloadChunkFile(job Chunk){
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
		log.Error("Could not download job file: ",err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)

	if job.JobFileType == ".zip"{
		log.Info("unzipping...")
		//Unzip(config.Working_dir+filename, config.Working_dir)
		unzip.Extract(config.Working_dir+filename, config.Working_dir)
		log.Info("unzipped.")

		files, _ := ioutil.ReadDir(config.Working_dir)

		found := false

		for _, f := range files {
			log.Debug("checking: "+f.Name())
			if filepath.Ext(config.Working_dir+f.Name()) == ".blend" {
				log.Info("found .blend")
				os.Rename(config.Working_dir+f.Name(), config.Working_dir+"job.blend")
				found = true
				break
			}
		}

		if !found {
			log.Error("No .blend found in archive")
			os.Exit(1)
		}


	}

}

func checkForChunk() (bool, Chunk){
	log.Info("Checking for job")
	req, err := http.NewRequest("GET",
		config.Bcr_server+"/worker/"+strconv.Itoa(worker.Id)+"/job", nil)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		log.Info("no job.")
		return false, Chunk{}
	}

	log.Info("Got job!")

	body, _ := ioutil.ReadAll(resp.Body)


	job := Chunk{}

	json.Unmarshal(body, &job)

	return true, job

}

func registerWorker(){

	hostname,  _ := os.Hostname()
	worker = Worker{Name: hostname, Ip: 4544534} //TODO real ip

	payload, err := json.Marshal(worker)

	if err != nil {
		log.Error("Could not register worker:", err)
		os.Exit(1)	}


	req, err := http.NewRequest("POST", config.Bcr_server+"/worker", bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		log.Error("Could not register worker. Got status code: ", resp.StatusCode)
		os.Exit(1)
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
