package main

import "leaderboard/app/leaderboard"

// @title Leaderboard Service
// @version 0.0.4
// @description Simple & fast leaderboard service

// @contact.name Berke Emrecan Arslan
// @contact.url https://beremaran.com
// @contact.email berke.emrecan.arslan@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host leaderboard-v2-lb-ecs-tg-584908050.eu-central-1.elb.amazonaws.com
func main() {
	leaderboard.Run()
}
