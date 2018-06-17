package main

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/base64"
	"context"
	"google.golang.org/api/texttospeech/v1beta1"
	"golang.org/x/oauth2/google"
)

const (
	TEXT = "こんにちは！こちらは認証キーを使ったパターンです"
	GENDER = "FEMALE"
	LANGUAGE = "ja-JP"
	VOICENAME = "ja-JP-Standard-A"

	AUDIO_ENCODE = "LINEAR16"
	SPEAKING_RATE = 1

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

	b, err := ioutil.ReadFile("TextSpeech.json")
	if err != nil{
		panic(err)
	}
	config, err := google.JWTConfigFromJSON(b, texttospeech.CloudPlatformScope)
	client := config.Client(ctx)
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


