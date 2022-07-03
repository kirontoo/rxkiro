package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/db"
)

var pingedCount = 0

var botCommands = map[string]interface{}{
	"ping":       ping,
	"mefact":     randomMeFact,
	"animalfact": randomAnimalFact,
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

func randomFact(b *RxKiro, table string) string {
	query := fmt.Sprintf(`select * from "%s" order by random() limit 1`, table)
	rows, err := b.db.Store.Query(query)

	if err != nil {
		b.Log.Fatal().Err(err).Str("cmd", "animalFact").Send()
	}

	defer rows.Close()

	var fact db.FunFact
	for rows.Next() {
		f := db.FunFact{}
		if err := rows.Scan(&f.Id, &f.CreatedAt, &f.Value); err != nil {
			b.Log.Fatal().Err(err).Msg("finding random fact")
		}

		fact = f
	}

	return fact.Value
}

func randomAnimalFact(b *RxKiro) {
	msg := randomFact(b, "AnimalFacts")
	b.Send(msg)
}

func randomMeFact(b *RxKiro) {
	msg := randomFact(b, "FunFacts")
	b.Send(msg)
}

func (b *RxKiro) updateCounter(cmd db.Command) string {
	if cmd.Counter.Valid {
		count := cmd.Counter.Int64 + 1
		// update in db
		_, err := b.db.Store.Exec(fmt.Sprintf(`UPDATE "Commands" SET Counter = %d WHERE id = %d`, count, cmd.Id))
		if err != nil {
			b.Log.Error().Msg("Could not update counter.")
			return fmt.Sprintf("Could not update counter: %s", cmd.Name)
		}

		b.Log.Info().Int64("new_count", count).Int64("old_count", cmd.Counter.Int64).Str("cmd", cmd.Name).Msg("Update counter")

		return fmt.Sprintf("%s: %d", cmd.Name, count)
	}

	return fmt.Sprintf("Invalid counter: %s", cmd.Name)
}

func (b *RxKiro) replaceCmdVariables(rawCmd string, s string, msg twitch.PrivateMessage) string {
	raw := strings.ToLower(strings.Trim(rawCmd, "${}"))
	cmdVars := strings.Split(raw, " ")
	cmdVar := cmdVars[0]

	b.Log.Debug().Str("cmdvar", cmdVar).Send()
	b.Log.Print(cmdVars)

	switch cmdVar {
	case "user":
		username := "@" + msg.User.DisplayName
		updatedMsg := strings.ReplaceAll(s, rawCmd, username)
		return updatedMsg
	case "mention":
		var updatedMsg string
		m := strings.Split(msg.Message, " ")
		if len(m) < 2 {
			updatedMsg = strings.ReplaceAll(s, rawCmd, "")
		} else {
			mentionedUser := m[1]
			updatedMsg = strings.ReplaceAll(s, rawCmd, mentionedUser)
		}
		return updatedMsg
	case "random":
		value := randNum()
		if len(cmdVars) == 3 {
			// Has to be 3 for format: ${random minNum maxNum}
			max, err := strconv.Atoi(cmdVars[2])
			if err != nil {
				b.Log.Error().Str("cmd", cmdVar).Msg("Invalid params")
				return fmt.Sprintf("Invalid param: %s. Must be a number", cmdVars[2])
			}

			min, err := strconv.Atoi(cmdVars[1])
			if err != nil {
				b.Log.Error().Str("cmd", cmdVar).Msg("Invalid params")
				return fmt.Sprintf("Invalid param: %s. Must be a number", cmdVars[1])
			}

			if !isValidInt(max) || !isValidInt(min) {
				return "Invalid range value, should be greater or equal to 0"
			}

			if max < min {
				return fmt.Sprintf("Invalid max value, %d, should be less than %d", max, min)
			}

			value = randNum(min, max)
		}

		if len(cmdVars) == 2 {
			max, err := strconv.Atoi(cmdVars[1])
			if err != nil {
				b.Log.Error().Str("cmd", cmdVar).Msg("Invalid params")
				return fmt.Sprintf("Invalid param: %s. Must be a number", cmdVars[1])
			}

			if !isValidInt(max) {
				return "Invalid range value, should be greater or equal to 0"
			}

			value = randNum(max)
		}

		return strings.ReplaceAll(s, rawCmd, fmt.Sprintf("%d", value))
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

func randNum(randRange ...int) int {
	rand.Seed(time.Now().UnixNano())
	size := len(randRange)
	if size >= 2 {
		var min, max = randRange[0], randRange[1]
		randomNum := rand.Intn(max-min) + min
		return randomNum
	}

	if size == 1 {
		max := randRange[0]
		return rand.Intn(max)
	}

	return rand.Intn(100000)
}

func isValidInt(n int) bool {
	return !(n < 0)
}
