package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-zoox/fetch"
	"github.com/spf13/viper"
)

func Req(url, data string) []byte {
	var response *fetch.Response
	var err error
	if len(data) > 0 {
		req := fetch.Query{}
		req.Set("data", data)
		response, err = fetch.Post(url, &fetch.Config{Query: req})
		if err != nil {
			fmt.Printf("Can not reach the server: %v", err)
			os.Exit(1)
		}

	} else {
		response, err = fetch.Get(url)
		if err != nil {
			fmt.Printf("Can not reach the server: %v", err)
			os.Exit(1)
		}
	}

	return response.Body
}

func main() {
	var blockchain []*Block
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/gui-blockchain")
	viper.AddConfigPath("$HOME/.gui-blockchain")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			path, _ := os.LookupEnv("XDG_CONFIG_HOME")
			path += "/gui-blockchain/"
			os.MkdirAll(path, 0755)
			_, err1 := os.Create(path + "config.toml")
			if err1 != nil {
				fmt.Println(err1)
			}
		} else {
			fmt.Printf("Error: %v", err)
			return
		}
	}
	a := app.New()

	content := container.NewStack()
	w := a.NewWindow("GUI for blockchain")

	data := string(Req(viper.GetString("url"), ""))
	blocks := strings.Split(data, "|/")
	for _, v := range blocks {
		t := Deserialize([]byte(v))
		blockchain = append(blockchain, t)
	}
	show_blockchain := widget.NewList(
		func() int {
			return len(blockchain)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("Text Editor"))
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			text := co.(*fyne.Container).Objects[0].(*widget.Label)
			text.SetText(string(blockchain[lii].Name))
		},
	)
	show_blockchain.OnSelected = func(id widget.ListItemID) {
		t := container.NewCenter(
			container.NewVBox(
				widget.NewLabel(string(blockchain[id].Name)),
				widget.NewLabel("Hash: "+fmt.Sprintf("%x", blockchain[id].Hash)),
				widget.NewLabel("Previous hash: "+fmt.Sprintf("%x", blockchain[id].PrevHash)),
				widget.NewLabel("Timestamp: "+fmt.Sprintf("%v", blockchain[id].Timestamp)),
			),
		)
		content.Objects = []fyne.CanvasObject{t}
	}
	uploadbtn := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {})
	upload := container.NewHBox(layout.NewSpacer(), uploadbtn)

	mainsplit := container.NewHSplit(show_blockchain, content)
	mainsplit.Offset = 0.2
	split := container.NewVSplit(upload, mainsplit)
	split.Offset = 0.025
	w.SetContent(split)
	w.ShowAndRun()
}
