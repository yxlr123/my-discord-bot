import discord
import requests
import re
import time
# from disnake.ext import commands
import datetime


class MyClient(discord.Client):

    async def on_ready(self):
        print('Logged on as', self.user)

    async def on_message(self, message):
        # 防止机器人自己回自己（
        if message.author == self.user:
            return
        today = datetime.datetime.now().weekday() + 1
        now_hour = datetime.datetime.now().hour
        r18 = "0"
        print(str(message.author) + " 说: " + str(message.content))
        print()
        zhenxun = 0
        setuu = 0
        if zhenxun == 1 and message.content != "":
            await message.channel.send("快去看真寻")
        if message.content != "已经在康了" and today == 5:
            zhenxun = 1
        if message.content[0:4:] == "让我康康":
            r18 = "1"
            setuu = 1
            if len(message.content) > 4:
                keyword = message.content[4::]
        if message.content[0:2:] == '我要' and message.content[-2::] == "色图":
            setuu = 1
            if len(message.content) > 4:
                keyword = message.content[2:-3:1]
        if setuu == 1 and message.content != "":
            await message.channel.send("年轻人，要节制")
            print("有人点了一份色图")
            if len(message.content) == 4:
                url = f"https://api.lolicon.app/setu/v2?r18={r18}&proxy=yxlr-cdn.ml"
            else:
                url = f"https://api.lolicon.app/setu/v2?r18={r18}&proxy=yxlr-cdn.ml&keyword={keyword}"
            setu = requests.get(url).json()
            if setu["data"] == []:
                await message.channel.send("好像没有这个标签的色图")
            else:
                tags = str(setu["data"][0]["tags"])[1:-1:1]
                setu_url = setu['data'][0]['urls']["original"]
                # setu_url = "https://pixiv.re/" + str(setu["data"][0]["pid"]) + ".png"
                title = setu["data"][0]["title"]
                author = setu["data"][0]["author"]
                size = str(setu["data"][0]["width"]) + "x" + str(
                    setu["data"][0]["height"])
                embed = discord.Embed(title=title, url=setu_url)
                embed.set_author(name=author)
                embed.set_image(url=setu_url)
                embed.add_field(name='R18',
                                value=setu['data'][0]["r18"],
                                inline=True)
                embed.add_field(name='Size', value=size, inline=True)
                embed.add_field(name='Tags', value=tags, inline=False)
                embed.set_footer(text="Powerd by api.lolicon.app")
                await message.channel.send(embed=embed)
        if message.content == "不可以瑟瑟":
            await message.channel.send("一分钟后继续瑟瑟")
            print("有坏蛋暂停了瑟瑟")
            time.sleep(60)
            await message.channel.send("我又来了，嘿嘿嘿")

        if message.content == "我要bing图":
            bing_img = requests.get("https://api.muvip.cn//api/bing/index.php?rand=false&day=0&size=1920x1080&info=true").json()
            bing_title = str(bing_img["title"])
            bing_url = bing_img["url"]
            bing_link = bing_img["link"]
            embed = discord.Embed(title=bing_title, url=bing_link)
            embed.set_image(url=bing_url)
            await message.channel.send(embed=embed)

        if message.content == "我要一言":
            hitokoto_all = requests.get("https://v1.hitokoto.cn/?c=a").json()
            hitokoto = str(hitokoto_all["hitokoto"])
            hitokoto_from = "《" + hitokoto_all["from"] + "》"
            hitokoto_who = hitokoto_all["from_who"]
            embed = discord.Embed(title=hitokoto)
            if hitokoto_who:
                kongkong = "ㅤ"*len(hitokoto) + "——" + hitokoto_who + hitokoto_from
            else:
                kongkong = "ㅤ"*len(hitokoto) + "——" + hitokoto_from
            embed.set_footer(text=kongkong)
            await message.channel.send(embed=embed)


intents = discord.Intents.default()
intents.message_content = True
client = MyClient(intents=intents)
client.run("token")
