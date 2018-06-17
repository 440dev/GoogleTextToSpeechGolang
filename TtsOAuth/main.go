package main

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"fmt"
	"bufio"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/base64"
	"context"
	"google.golang.org/api/texttospeech/v1beta1"
)

const (
	TEXT = "こんにちは！こちらはOAuthのパターンです"
	GENDER = "FEMALE"
	LANGUAGE = "ja-JP"
	VOICENAME = "ja-JP-Standard-A"

	AUDIO_ENCODE = "LINEAR16"
	SPEAKING_RATE = 1

	CLIENT_ID = " Your Client ID "
	SECRET_KEY = " Your Secret KEY "

	OUTPUT = "output.mp3"
)

func main(){
	synthesizeSpeechRequest := &texttospeech.SynthesizeSpeechRequest{
		AudioConfig: &texttospeech.AudioConfig{
			AudioEncoding: AUDIO_ENCODE,
			SpeakingRate: SPEAKING_RATE},
		Input: &texttospeech.SynthesisInput{
			Text: TEXT},
		Voice: &texttospeech.VoiceSelectionParams{
			Name:         VOICENAME,
			LanguageCode: LANGUAGE,
			SsmlGender: GENDER}}

	payload, _ := json.Marshal(synthesizeSpeechRequest)

	ctx := context.Background()
	config := oauth2.Config{
		ClientID:CLIENT_ID,
		ClientSecret:SECRET_KEY,
		Endpoint:oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL:"urn:ietf:wg:oauth:2.0:oob",
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	url := config.AuthCodeURL("")
	fmt.Println("Get authorize code:", url)
	fmt.Print("Input code:")

	var s string
	var sc = bufio.NewScanner(os.Stdin)
	if sc.Scan(){
		s = sc.Text()
	}

	token, err := config.Exchange(ctx, s)
	if err != nil {
		panic(err)
		return
	}

	client := config.Client(ctx, token)

	req, err := http.NewRequest(
		"POST",
		"https://texttospeech.googleapis.com/v1beta1/text:synthesize",
		bytes.NewReader(payload),
	)

	if err != nil{
		panic(err)
		return
	}

	response, err := client.Do(req)
	if err != nil{
		panic(err)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK{
		panic(err)
		return
	}

	var audio texttospeech.SynthesizeSpeechResponse
	if err := json.Unmarshal(body, &audio); err != nil{
		panic(err)
	}

	data, _ := base64.StdEncoding.DecodeString(audio.AudioContent)
	file, _:=os.Create(OUTPUT)
	defer file.Close()
	file.Write(data)

	fmt.Println("Generate Speech File:",OUTPUT)

	return
}