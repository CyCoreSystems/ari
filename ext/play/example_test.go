package play

import (
	"context"
	"fmt"
	"strings"

	"github.com/CyCoreSystems/ari"
)

func ExamplePlay() {
	h := ari.NewChannelHandle(key, "exampleChannel")

	res, err := Play(context.TODO(), h,
		URI("sound:tt-monkeys", "sound:vm-goodbye"),
	).Result()
	if err != nil {
		fmt.Println("Failed to play audio", err)
		return
	}

	if len(res.DTMF) > 0 {
		fmt.Println("Got a DTMF during playback:", res.DTMF)
	}
	return
}

func ExamplePlay_async() {
	h := ari.NewChannelHandle(key, "exampleChannel")

	bridgeSub := h.Subscribe(ari.EventTypes.ChannelEnteredBridge)
	defer bridgeSub.Cancel()

	sess := Play(context.TODO(), h,
		URI("characters:ded", "sound:tt-monkeys",
			"number:192846", "digits:43"),
	)

	select {
	case <-bridgeSub.Events():
		fmt.Println("Channel entered bridge during playback")
	case <-sess.Done():
		if sess.Err() != nil {
			fmt.Println("Prompt failed", err)
		} else {
			fmt.Println("Prompt complete")
		}
	}
	return
}
func ExamplePrompt() {
	h := ari.NewChannelHandle(key, "exampleChannel")

	res, err := Prompt(context.TODO(), h,
		URI("tone:1004/250", "sound:vm-enter-num-to-call",
			"sound:astcc-followed-by-pound"),
		MatchHash(), // match any digits until hash
		Replays(3),  // repeat prompt up to three times, if no match
	).Result()
	if err != nil {
		fmt.Println("Failed to play", err)
		return
	}

	if res.MatchResult == Complete {
		fmt.Println("Got valid, terminated DTMF entry", res.DTMF)
		// hash is automatically trimmed from res.DTMF
	}
	return
}

func ExamplePrompt_customMatchFunc() {
	h := ari.NewChannelHandle(key, "exampleChannel")

	res, err := Prompt(context.TODO(), h,
		URI("sound:agent-user"),
		MatchFunc(func(in string) (string, MatchResult) {
			// This is a custom match function which will
			// be run each time a DTMF digit is received
			pat := strings.TrimSuffix(in, "#")

			user := db.Lookup(pat)
			if user == nil {
				if pat != in {
					// pattern was hash-terminated but no match
					// was found, so there is no match possible
					return pat, Invalid
				}
				return in, Incomplete
			}
			return pat, Complete
		}),
	).Result()
	if err != nil {
		fmt.Println("Failed to play prompt", err)
		return
	}

	if res.MatchResult == Complete {
		fmt.Println("Got valid user", res.DTMF)
	}
	return
}
