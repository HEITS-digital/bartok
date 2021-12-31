package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getRadomWelcomeMessage(user string) string {
	rand.Seed(time.Now().Unix())
	messages := []string{
		"Hello <@%v>, and welcome to HEITS :wave:",
		"HOORAY!\nWelcome to the team <@%v> :dog_hooray:",
		"Welcome aboard <@%v>\nWoot Woot :woot_woot:",
		"Welcome to the tribe <@%v>. Bravo! :blob_clap:",
		"<@%v> just joined our family!\nWelcome to the club :mafia:",
	}
	n := rand.Int() % len(messages)
	return fmt.Sprintf(messages[n], user)
}

func getNewMemberDM() string {
	return fmt.Sprintf(teamJoinWelcomeMessageFormat, "CEC0Z16QL", "CSKGXKXS5", "C02054LCV6E", "CEC2Y6QD9", "C01S8NR19TR", "C01NY7FN34Y")
}

var randomEnReplies = []string{
	"Hey <@%v>. Where is the beef?",
	"Sorry <@%v> but I can't deal with you now.\nThis week is so very busy and my skin is broken",
	"Yes <@%v>\nI have superpowers because I was born at a very young age",
	"Stand back <@%v>, your hair makes me nervous",
	"Hey <@%v>.\nWould you like to kiss my flamingo? :flamingo:",
	"<@%v> on a scale of 1 to 5, how anxious are you when using public bathrooms?",
	"Stop asking for my number <@%v>!!!",
	"Are you afraid of raccoons <@%v>?",
	"Pickled cabbage -> that's my secret\nWhat's yours <@%v>?",
	"<@%v> -> :talktothehand:",
	"<@%v> -> :lalala:",
}

var randomRoReplies = []string{
	"Hai sa lasam prostiile pt alta data <@%v>",
	"<@%v>, pt binele tuturor hai sa pretindem ca ti-ai vazut de treaba ta",
	"Hey hey, pare ca <@%v> si-a terminat mai devreme treaba azi. Asta da productivitate!",
	"Hopa si <@%v>. Orice numai sa treaca ziua cat mai repede",
	"N-am timp acum <@%v>. Hai sa ne auzim mai tarziu... mult mai tarziu",
	"Uite cum a mai trecut un an si <@%v> tot nu se invata minte",
	"Ia zi <@%v>, cat ai dat pentru bagatu-n seama?",
	"Iar incepe <@%v> cu cicaleala\nNici nu se putea altfel!",
	"<@%v>, stii bancul cu iarna?\nIar n-ai de lucru?",
	"<@%v>, stii bancul cu Nicu?\nNi cum trece ziua si n-ai lucrat nimic",
}

var randomEnAnswers = []string{
	"Is that a serious question <@%v>?",
	"Hey <@%v> you know what?\nI’ll answer you in a bit. I’m now waiting for motivation to build up",
	"Sorry <@%v> but I just found something more important to deal with at this moment",
	"I have no idea how to answer this <@%v> :thinking_face:",
	"I don't think I'm qualified to answer this now <@%v>",
	"Only questions and questions... You know, I have questions too <@%v>.\nBut no one's curious",
}

var randomRoAnswers = []string{
	"Mda… Alta intrebare <@%v>! :alta-intrebare:",
	"Habar n-am ce sa-ti raspund la asta <@%v>. Lasa-ma sa ma mai gandesc",
	"Haha <@%v>. Ce te face sa crezi ca am timp pt intrebari acum?",
	"Revin imediat cu un raspuns <@%v>. Momentan mi-am luat o pauza pt gustare :leafy_green:",
	"Scuze <@%v>, dar momentan n-am destula motivatie sa-ti raspund",
	"Iti raspund mai tarziu <@%v>\nDeocamdata nu simt nevoia",
	"Nu stiu... dar zi-ne tu <@%v>. Pare ca deja le stii pe toate",
	"Raspunde-ti singur <@%v>. Oricum, nu cred ca-i prima data cand aleg sa te ignor",
	"Sa stii ca exista si intrebari castigatoare <@%v>\nInsa la cum te stiu, slabe sanse sa vina de la tine",
}

const teamJoinWelcomeMessageFormat string = `Welcome to HEITS.digital :wave: ! We are super excited that you joined us, and wish you the best of luck on this new adventure. 
I’m Bartók the goat, and I am here to share some useful information with you:
*1. Internal meetings*
- Each Monday at 11am we have the Internal & Informal meeting, where we disc	uss important company updates.
- Once a month we meet and share knowledge, during the HEITS talks initiative. Come and find out cool stuff, both technical and non-technical.
*2. Slack channels*
- If you ever need help from our workspace’s administrators, please reach out in <#%s>
- Engineering -> <#%s>
- Administrative & Financial stuff -> <#%s>
- Games, Hobbies & Fun -> <#%s>, <#%s> & <#%s>
There are quite a few other channels, depending on your interests or location. Just click on the :heavy_plus_sign: next to the channel list in the sidebar, and click Browse Channels to search for anything that interests you.
*3. PTO*
- This is our vacation planner https://heims.heits.digital/. You can authenticate using your Google account and add your vacation days here. Your Google calendar will later reflect the PTO days.
- For any other information regarding our benefits, or other administrative aspects, you can always reach Lidia Rusu from HR or Florina Condulet from Finance & Administration.
*4. Stay connected*
- Here’s our website https://heits.digital/ - check it out
- Facebook page: https://www.facebook.com/heits.digital - Like & Share
- Linkedin page: https://www.linkedin.com/company/heits-digital/ - Follow & Share
Hope I could be of help and I am working on adding new useful functions. If you have any suggestions, please drop a message to the engineering team.
Sit back, relax, enjoy our community and have fun! :happygoat:`
