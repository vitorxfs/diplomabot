package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lpernett/godotenv"
)

const studentName string = "Vitor Sanches";
const url string = "https://igce.rc.unesp.br/index.php"

type XmlData struct {
	Cmd []Cmd `xml:"cmd"`
}

type Cmd struct {
	Key string `xml:"t,attr"`
	Value string `xml:",chardata"`
}

func findBodyHTML(xmlBinary []byte) (string, error) {
	var xmlData XmlData
	xml.Unmarshal(xmlBinary, &xmlData)

	cmdBodyIndex := slices.IndexFunc(xmlData.Cmd, func (cmd Cmd) bool {
		return cmd.Key == "idCorpo"
	})

	if (cmdBodyIndex == -1) {
		return "", errors.New("could not find idCorpo in xml")
	}

	return xmlData.Cmd[cmdBodyIndex].Value, nil
}

func findPDFLink(htmlString string) (string, error) {
	reader := strings.NewReader(htmlString)
	q, error := goquery.NewDocumentFromReader(reader)
	if (error != nil) { return "", error }

	anchor := q.Find("#wrapper-gridder-block-315 > div > a");
	href, exists := anchor.Attr("href");
	if (!exists) { return "", errors.New("attribute href not found") }

	return href, nil
}

type PdfParseResponse struct {
	Text string `json:"text"`
}

func readPdfFromUrl(pdfUrl string) (string, error) {
	token := os.Getenv("UTILS_TOKEN")
	baseUrl := os.Getenv("UTILS_BASE_URL")
	url := strings.Join([]string{baseUrl, "api/pdf?", "url=", pdfUrl, "&token=", token}, "")

	req, error := http.NewRequest("GET", url, nil);
	if (error != nil) { return "", error }

	res, error := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if (error != nil) { return "", error }

	body, _ := io.ReadAll(res.Body);

	var pdfParseResponse PdfParseResponse;
	error = json.Unmarshal(body, &pdfParseResponse)
	if (error != nil) { return "", error }

	return pdfParseResponse.Text, nil
}

func containsName(text string, name string) bool {
	text = strings.ToLower(text);
	name = strings.ToLower(name);
	return strings.Contains(text, name);
}

func sendTelegram(msg string) error {
	botId := os.Getenv("TELEGRAM_BOT_ID")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botId, chatId, msg);

	req, error := http.NewRequest("GET", url, nil);
	if (error != nil) { return error }

	res, error := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if (error != nil) { return error }
	if (res.StatusCode != 200) { return errors.New(fmt.Sprintf("error %d when sending message to telegram"))}

	return nil;
}

func getSTGBodyXML() ([]byte, error) {
	payload := strings.NewReader("xajax=exibeCorpo&xajaxr=1726772075210&xajaxargs%5B%5D=942")

	req, error := http.NewRequest("POST", url, payload)
	if (error != nil) { return nil, error }

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, error := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if (error != nil) { return nil, error }

	return io.ReadAll(res.Body)
}

func main() {
	godotenv.Load()

	fmt.Println("Starting...")

	body, error := getSTGBodyXML()
	if (error != nil) { panic(error) }

	html, error := findBodyHTML(body)
	if error != nil { panic(error) }

	url, error := findPDFLink(html)
	if error != nil { panic(error) }

	pdfText, error := readPdfFromUrl(url);
	if error != nil { panic(error) }

	if (containsName(pdfText, studentName)) {
		error := sendTelegram("🚨 Diploma está pronto para ser retirado!! 🎉🎉🎉")
		if error != nil { panic(error) }
		return;
	}

	error = sendTelegram("Diploma ainda não está pronto para ser retirado")
	if error != nil { panic(error) }

	return;
}
