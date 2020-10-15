package main

import "leaderboard/app/leaderboard"

// @title Leaderboard Service
// @version 0.0.4
// @description Simple & fast leaderboard service

// @contact.name Berke Emrecan Arslan
// @contact.url https://beremaran.com
// @contact.email berke.emrecan.arslan@gmail.com

// @license.name The MIT License (MIT)
// @license.url https://mit-license.org/

// @host leaderboard-v2-lb-ecs-tg-584908050.eu-central-1.elb.amazonaws.com
func main() {
	leaderboard.Run()
}
