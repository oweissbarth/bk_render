
func request(){

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
