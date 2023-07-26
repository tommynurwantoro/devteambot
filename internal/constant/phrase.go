package constant

type Phrase struct {
	Kill    []string
	Idle    []string
	Suicide []string
}

func NewPhrase() Phrase {
	phrase := Phrase{
		Kill: []string{
			"%s :crossed_swords: >> **%s** slit **%s** throat with a sharp rock.\n",
			"%s :crossed_swords: >> **%s** stabbed **%s** with a cucumber. What an unfortunate series of events!\n",
			"%s :crossed_swords: >> **%s** stabbed **%s** with a pickle after he/she tried to stab them with a cucumber.\n",
			"%s :crossed_swords: >> **%s** slapped **%s** with an eggplant and died of a head concussion.\n",
			"%s :crossed_swords: >> **%s** accidentally step over **%s** eyes and crush the brain.\n",
			"%s :crossed_swords: >> **%s** decided to try bungee jumping, but the operator, **%s**, didn't attach the line.\n",
			"%s :crossed_swords: >> **%s** come up behind in a car, **%s** had his AirPods in and...Bye!\n",
			"%s :crossed_swords: >> **%s** drew his weapon faster than **%s** in a duel.\n",
			"%s :crossed_swords: >> **%s** slipped on a eggplant and accidentally elbowed **%s** to death.\n",
			"%s :crossed_swords: >> **%s** one hundred hand slapped **%s** and knocked their soul out their back.\n",
			"%s :crossed_swords: >> **%s** killed **%s** by shoving a cucumber down their throat. Bon appetit!\n",
			"%s :crossed_swords: >> **%s** didn't like **%s** took the first slice of pizza\n",
			"%s :crossed_swords: >> **%s** ran over **%s** with their steamboat, then reversed back over the peices, just to be sure the job was done.\n",
		},
		Idle: []string{
			"%s :shushing_face: >> **%s** was meditating under the waterfalls.\n",
			"%s :shushing_face: >> **%s** camped at a cave because it was dark outside.\n",
			"%s :shushing_face: >> **%s** missed an arrow.\n",
			"%s :shushing_face: >> **%s** listened to his favorite music.\n",
			"%s :shushing_face: >> **%s** pet a cat.\n",
			"%s :shushing_face: >> **%s** found a coin on the floor.\n",
			"%s :shushing_face: >> **%s** found a broken arrow on the floor.\n",
			"%s :shushing_face: >> **%s** learnt water breathing style from tanjiro.\n",
			"%s :shushing_face: >> **%s** foraged plants for dinner.\n",
			"%s :shushing_face: >> **%s** helped a turtle cross the road.\n",
		},
		Suicide: []string{
			"%s :skull_crossbones: >> **%s** tried to swim in lava.\n",
			"%s :skull_crossbones: >> **%s** was practicing knife throwing when the knife bounced back and stabbed him in the heart.\n",
			"%s :skull_crossbones: >> **%s** hired a fake hitman off the internet and got into a shootout with the police.\n",
			"%s :skull_crossbones: >> **%s** was killed by Cyborg Khuga.\n",
			"%s :skull_crossbones: >> **%s** stepped on a Lego.\n",
			"%s :skull_crossbones: >> **%s** used the wrong side of the sword. What an idiot.\n",
			"%s :skull_crossbones: >> **%s** got trampled by a bull crypto market.\n",
			"%s :skull_crossbones: >> **%s** lost a deadly game of Rock, Paper, Scissors.\n",
			"%s :skull_crossbones: >> **%s** discovered that the floor is lava.\n",
			"%s :skull_crossbones: >> **%s** had a good vegetarian meal, didn't know the main vegetable was Venusaur.\n",
			"%s :skull_crossbones: >> **%s** fought the law. The law won.\n",
			"%s :skull_crossbones: >> **%s** stepped on a bee.\n",
			"%s :skull_crossbones: >> **%s** slipped on a banana peel and split their head open.\n",
			"%s :skull_crossbones: >> **%s** discovered what happens when Superman sneezes in your face.\n",
			"%s :skull_crossbones: >> **%s** entered a wormhole.\n",
			"%s :skull_crossbones: >> **%s** was trapped in an infinite loop of tik tok videos. Found their corpse a year later.\n",
			"%s :skull_crossbones: >> **%s** died by eating pineapple on pizza, yuck!\n",
			"%s :skull_crossbones: >> **%s** died from eating to much ice cream. He got a brain freeze and died a cold lonely sad death.\n",
			"%s :skull_crossbones: >> **%s** lost a stare contest to a 3 eyed frog.\n",
			"%s :skull_crossbones: >> **%s** spent too much time researching what an NFT is... and died.\n",
			"%s :skull_crossbones: >> **%s** listed their NFT below floor price and was brutally murdered by the community.\n",
			"%s :skull_crossbones: >> **%s** died because of choked by it's saliva at sleep.\n",
			"%s :skull_crossbones: >> **%s** fell off a ladder\n",
			"%s :skull_crossbones: >> **%s** was struct by lightning\n",
		},
	}

	return phrase
}
