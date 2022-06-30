package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/db"
)

var pingedCount = 0

var botCommands = map[string]interface{}{
	"ping":   ping,
	"mefact": randomMeFact,
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

func randomAnimalFact(b *RxKiro) {
	query := fmt.Sprintf(`select * from "AnimalFacts" order by random() limit 1`)
	rows, err := b.db.SqlDb.Query(query)

	if err != nil {
		b.Log.Fatal().Err(err).Str("cmd", "animalFact").Send()
	}

	defer rows.Close()

	var fact db.FunFact
	for rows.Next() {
		f := db.FunFact{}
		if err := rows.Scan(&f.Id, &f.CreatedAt, &f.Value); err != nil {
			b.Log.Fatal().Err(err).Str("cmd", "animalFact").Msg("Random me fact")
		}

		fact = f
	}

	b.Send(fact.Value)
}

func randomMeFact(b *RxKiro) {
	query := fmt.Sprintf(`select * from "FunFacts" order by random() limit 1`)
	rows, err := b.db.SqlDb.Query(query)

	if err != nil {
		b.Log.Fatal().Err(err).Str("cmd", "meFact").Send()
	}

	defer rows.Close()

	var fact db.FunFact
	for rows.Next() {
		f := db.FunFact{}
		if err := rows.Scan(&f.Id, &f.CreatedAt, &f.Value); err != nil {
			b.Log.Fatal().Err(err).Str("cmd", "meFact").Msg("Random me fact")
		}

		fact = f
	}

	b.Send(fact.Value)
}

func (b *RxKiro) runCounter(cmd db.Command) string {
	if cmd.Counter.Valid {
		count := cmd.Counter.Int64 + 1
		// update in db
		_, err := b.db.SqlDb.Exec(fmt.Sprintf(`UPDATE "Commands" SET Counter = %d WHERE id = %d`, count, cmd.Id))
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
