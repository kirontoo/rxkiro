package bot

import (
	"encoding/json"
	"strconv"
)

var pingedCount = 0

var botCommands = map[string]interface{}{
	"ping":       ping,
	"animalfact": randomAnimalFact,
	"mefact": randomMeFact,
}

func ping(b *RxKiro) {
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

func randomMeFact(b *RxKiro) {
	res := b.db.Rpc("rand_fun_fact", "exact", map[string]string{})
	if res != "" {
		fact := trimQuotes(res)
		b.Log.Info().Str("fact", fact).Msg("Random animal fact")
		b.Send(fact)
	}
}

func randomAnimalFact(b *RxKiro) {
	res := b.db.Rpc("rand_animal_fact", "exact", map[string]string{})
	if res != "" {
		fact := trimQuotes(res)
		b.Log.Info().Str("fact", fact).Msg("Random animal fact")
		b.Send(fact)
	}
}

func (b *RxKiro) runCounterCmd(data Command) string {
	data.Counter = incrementCounter(data.Counter)
	b.Log.Debug().Int64("counter", data.Counter).Msg("Should increment")
	updatedRes, _, err := b.db.From(CommandTable).Update(data, "representation", "exact").Eq("name", data.Name).Execute()

	var jsonRes Command
	json.Unmarshal(updatedRes, &jsonRes)
	if err != nil {
		b.Log.Error().Interface("Command", jsonRes).Msg("Error updating")
		b.Log.Error().Err(err)
	}

	return data.Name + ": " + strconv.Itoa(int(data.Counter))
}

func incrementCounter(count int64) int64 {
	return (count + 1)
}
