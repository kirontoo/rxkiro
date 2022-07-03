package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

	var fact db.Fact
	for rows.Next() {
		f := db.Fact{}
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

/**
can take 0 to 2 parameters
if no parameter given, randNum will return a random number from 0 to 100000
if 1 parameter, n, given, randNum will return a random number from 0 to n
if 2 parameters, min and max, given, randNum will return a random number from min to max

If more than 2 parametrs is given, they will be ignored.
*/
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

func isPositiveInt(n int) bool {
	return !(n < 0)
}

func convertToInt(n string) (int, string) {
	value, err := strconv.Atoi(n)
	if err != nil {
		return -1, fmt.Sprintf("Invalid param: %s. Must be a number", n)
	}

	return value, ""
}

func randomCmd(s string, rawCmd string) string {
	raw := strings.ToLower(strings.Trim(rawCmd, "${}"))
	cmdVars := strings.Split(raw, " ")

	value := randNum()
	if len(cmdVars) == 3 {
		max, err := convertToInt(cmdVars[2])
		if len(err) > 0 {
			return err
		}

		min, err := convertToInt(cmdVars[1])
		if len(err) > 0 {
			return err
		}

		if !isPositiveInt(max) || !isPositiveInt(min) {
			return "Invalid range value, should be greater or equal to 0"
		}

		if max < min {
			return fmt.Sprintf("Invalid max value, %d, should be less than %d", max, min)
		}

		value = randNum(min, max)
	}

	if len(cmdVars) == 2 {
		max, err := convertToInt(cmdVars[1])
		if len(err) > 0 {
			return err
		}

		if !isPositiveInt(max) {
			return "Invalid range value, should be greater or equal to 0"
		}

		value = randNum(max)
	}

	return strings.ReplaceAll(s, rawCmd, fmt.Sprintf("%d", value))
}
