package setu

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

var (
	SetuCmd = discordgo.ApplicationCommand{
		Name:        "getsetu",
		Description: "请求一张色图",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "r18",
				Description: "(默认为\"false\")",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "true",
						Value: "1",
					},
					{
						Name:  "false",
						Value: "0",
					},
				},
			},
			{
				Name:        "keyword",
				Description: "图片的关键词",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}

	CommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"getsetu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			setuOptions := i.ApplicationCommandData().Options
			setuR18 := "0"
			keyword := ""
			for _, r := range setuOptions {
				if r.Name == "r18" {
					setuR18 = r.StringValue()
				}
				if r.Name == "keyword" {
					keyword = r.StringValue()
				}
			}

			embed, err := GetSetu(setuR18, keyword)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
			}
		},
	}
)

func GetSetu(r18 string, key string) (*discordgo.MessageEmbed, error) {
	client := resty.New() // 创建一个restry客户端
	resp, err := client.R().
	    SetBody(map[string]string{
			"r18":r18,
			"keyword":key,
			"proxy":"yxlr-cdn.ml",
		}).
	    Post(fmt.Sprintf("https://api.lolicon.app/setu/v2?r18=%v&keyword=%v&proxy=yxlr-cdn.ml",r18,key))
	if err != nil {
		return nil , err
	}
	json := string(resp.Body())
	// 判断api报错或者内容为空的情况（keyword没找到）
	if a := gjson.Get(json, "error").String(); a != "" {
		return nil, errors.New("api出错：" + a)
	}
	setu := gjson.Get(json, "data.0")
	if setu.Raw == "" {
		log.Println(json)
		err = errors.New("api返回了空值")
		return nil, err
	}

	setuR18 := &discordgo.MessageEmbedField{
		Name:   "R18",
		Value:  setu.Get("r18").String(),
		Inline: false,
	}
	setuSize := &discordgo.MessageEmbedField{
		Name:   "Size",
		Value:  setu.Get("width").String() + "x" + setu.Get("height").String(),
		Inline: false,
	}
	var tags string
	for _, tag := range setu.Get("tags").Array() {
		tags = tag.Str + "," + tags
	}
	setuTags := &discordgo.MessageEmbedField{
		Name:  "Tags",
		Value: tags,
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://www.pixiv.net/users/" + setu.Get("uid").String(),
			Name:    setu.Get("author").String(),
			IconURL: "https://i.imgur.com/pECIFHB.png",
		},
		URL:    setu.Get("urls.original").String(),
		Title:  setu.Get("title").String(),
		Fields: []*discordgo.MessageEmbedField{setuR18, setuSize, setuTags},
		Image: &discordgo.MessageEmbedImage{
			URL:    setu.Get("urls.original").String(),
			Width:  int(setu.Get("width").Int()),
			Height: int(setu.Get("height").Int()),
		},
	}
	return embed, nil
}
