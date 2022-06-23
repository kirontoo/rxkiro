package bot

import (
	"regexp"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

var pingedCount = 0

var botCommands = map[string]interface{}{
	"ping": ping,
}

func ping(b *RxKiro) {
	b.Log.Log().Msg("pinged")
	pingedCount += 1

	switch pingedCount {
	case 3:
		b.Send("...pong...")
		break
	case 4:
		b.Send("pong!")
		break
	case 5:
		b.Send("PONG!")
		b.Send("Where is this going?")
		pingedCount = 0
	default:
		b.Send("pong")
	}
}

// func randomMeFact(b *RxKiro) {
// 	res := b.db.Rpc("rand_fun_fact", "exact", map[string]string{})
// 	if res != "" {
// 		fact := trimQuotes(res)
// 		b.Log.Info().Str("fact", fact).Msg("Random animal fact")
// 		b.Send(fact)
// 	}
// }

// func randomAnimalFact(b *RxKiro) {
// 	res := b.db.Rpc("rand_animal_fact", "exact", map[string]string{})
// 	if res != "" {
// 		fact := trimQuotes(res)
// 		b.Log.Info().Str("fact", fact).Msg("Random animal fact")
// 		b.Send(fact)
// 	}
// }

// func (b *RxKiro) runCounterCmd(data Command) string {
// 	data.Counter = incrementCounter(data.Counter)
// 	b.Log.Debug().Int64("counter", data.Counter).Msg("Should increment")
// 	updatedRes, _, err := b.db.From(CommandTable).Update(data, "representation", "exact").Eq("name", data.Name).Execute()

// 	var jsonRes Command
// 	json.Unmarshal(updatedRes, &jsonRes)
// 	if err != nil {
// 		b.Log.Error().Interface("Command", jsonRes).Msg("Error updating")
// 		b.Log.Error().Err(err)
// 	}

// 	return data.Name + ": " + strconv.Itoa(int(data.Counter))
// }

func incrementCounter(count int64) int64 {
	return (count + 1)
}

func (b *RxKiro) replaceCmdVariables(rawCmd string, s string, msg twitch.PrivateMessage) string {
	cmdVar := strings.ToLower(strings.Trim(rawCmd, "${}"))
	switch cmdVar {
	case "user":
		username := "@" + msg.User.DisplayName
		updatedMsg := strings.ReplaceAll(s, rawCmd, username)
		return updatedMsg
	case "mention":
		var updatedMsg string
		m := strings.Split(msg.Message, " ")
		if len(m) < 2 {
			// TODO: what if there was no username mention?
			// should be a better response for this
			// respond with a !ERR: no mention found
			updatedMsg = strings.ReplaceAll(s, rawCmd, "")
		} else {
			mentionedUser := m[1]
			updatedMsg = strings.ReplaceAll(s, rawCmd, mentionedUser)
		}
		return updatedMsg
	default:
		return ""
	}
}

func (b *RxKiro) findCmdVars(s string) []string {
	var matches = []string{}

	r, err := regexp.Compile(`\${(.*?)}`)
	if err != nil {
		b.Log.Error().Err(err).Send()
	}

	matched := r.Match([]byte(s))
	if matched {
		matches = r.FindAllString(s, -1)
	}

	return matches
}
