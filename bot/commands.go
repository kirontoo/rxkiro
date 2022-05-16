package bot

import (
	"encoding/json"
	"strconv"
)

var botCommands = map[string]interface{}{
	"ping": func(b *RxKiro) {
		// TODO after pinged 5 times in a row, send back a special msg might need to save a the time and a db record for this
		b.Send("pong")
		b.Send("where is this going?")
	},
	"animalFact": randomAnimalFact,
}

func randomAnimalFact(b *RxKiro) {
	const table = "AnimalFacts"
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
