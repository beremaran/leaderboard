package services_test

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"leaderboard/app/leaderboard/tasks"
	"testing"
)
import . "github.com/onsi/ginkgo"
import . "github.com/onsi/gomega"

const KeyPrefix = "LB_"

var mRedis *miniredis.Miniredis

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "services")
}

var _ = BeforeSuite(func() {
	var err error
	mRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
	mRedis.Close()
})

var _ = Describe("the redis service", func() {
	var (
		redisService api.RedisService
		userService  *services.UserService
	)

	JustBeforeEach(func() {
		userService, redisService = buildDependencies(mRedis.Addr())
	})

	JustAfterEach(func() {
		mRedis.FlushAll()
	})

	Context("RedisService.GetSortedSetSize()", func() {
		When("given non-existent sorted set", func() {
			It("returns error", func() {
				_, err := redisService.GetSortedSetSize(uuid.New().String())
				Expect(err).NotTo(BeNil())
			})
		})

		When("given a valid sorted set", func() {
			It("returns correct size", func() {
				generateUsers(userService, redisService, 33)

				size, err := redisService.GetSortedSetSize("GLOBAL")
				Expect(err).To(BeNil())
				Expect(size).To(BeEquivalentTo(33))
			})
		})
	})

	Context("RedisService.GetScore()", func() {
		When("given non-existent sorted set", func() {
			It("returns error", func() {
				_, err := redisService.GetScore(uuid.New().String(), uuid.New().String())
				Expect(err).NotTo(BeNil())
			})
		})

		When("given a valid sorted set", func() {
			It("returns correct size", func() {
				guid, err := userService.Create(&api.UserProfile{
					UserId:      "a-guid",
					DisplayName: "hi",
					Country:     "XX",
					Points:      133,
				})
				Expect(err).To(BeNil())
				Expect(guid).To(BeEquivalentTo("a-guid"))

				score, err := redisService.GetScore("GLOBAL", "a-guid")
				Expect(err).To(BeNil())
				Expect(score).To(BeEquivalentTo(133))
			})
		})
	})
})

var _ = Describe("the user service", func() {
	var (
		userService  *services.UserService
		redisService api.RedisService
	)

	JustBeforeEach(func() {
		userService, redisService = buildDependencies(mRedis.Addr())
		generateUsers(userService, redisService, 20)
	})

	JustAfterEach(func() {
		mRedis.FlushAll()
	})

	Context("UserService.Create()", func() {

		When("profile is given with a non-empty UserId", func() {
			It("does not modify given UserId", func() {
				guid, err := userService.Create(&api.UserProfile{
					UserId:      "a-guid",
					DisplayName: "hi",
					Country:     "XX",
				})
				Expect(err).To(BeNil())
				Expect(guid).To(BeEquivalentTo("a-guid"))

				profile, err := userService.GetByID("a-guid")
				Expect(err).To(BeNil())
				Expect(profile.UserId).To(BeEquivalentTo("a-guid"))
			})
		})

		When("profile is given with a empty UserId", func() {
			It("generates a new GUID for the profile", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())
				Expect(guid).NotTo(BeEquivalentTo(""))
			})
		})
	})

	Context("UserService.GetByID()", func() {
		When("profile with given guid does not exist", func() {
			It("returns profile as nil", func() {
				profile, _ := userService.GetByID(uuid.New().String())
				Expect(profile).To(BeNil())
			})

			It("returns error", func() {
				_, err := userService.GetByID(uuid.New().String())
				Expect(err).NotTo(BeNil())
			})
		})
		When("profile with given GUID exist", func() {
			It("does not return error", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				_, err = userService.GetByID(guid)
				Expect(err).To(BeNil())
			})

			It("returns the profile with correct guid", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				returnedProfile, err := userService.GetByID(guid)
				Expect(err).To(BeNil())
				Expect(returnedProfile.UserId).To(BeEquivalentTo(guid))
			})

			It("returns the profile with correct country", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				returnedProfile, err := userService.GetByID(guid)
				Expect(err).To(BeNil())
				Expect(returnedProfile.Country).To(BeEquivalentTo(profile.Country))
			})

			It("returns the profile with correct display name", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				returnedProfile, err := userService.GetByID(guid)
				Expect(err).To(BeNil())
				Expect(returnedProfile.DisplayName).To(BeEquivalentTo(profile.DisplayName))
			})
		})
	})

	Context("UserService.GetByIDWithRank()", func() {
		When("profile with given guid does not exist", func() {
			It("returns profile as nil", func() {
				profile, _ := userService.GetByIDWithRank(uuid.New().String(), "GLOBAL")
				Expect(profile).To(BeNil())
			})

			It("returns error", func() {
				_, err := userService.GetByIDWithRank(uuid.New().String(), "GLOBAL")
				Expect(err).NotTo(BeNil())
			})
		})
		When("profile with given GUID exist", func() {
			It("returns the correct profile", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				returnedProfile, err := userService.GetByIDWithRank(guid, "GLOBAL")
				Expect(err).To(BeNil())
				Expect(returnedProfile.UserId).To(BeEquivalentTo(guid))
				Expect(returnedProfile.Rank).NotTo(BeEquivalentTo(0))
				Expect(returnedProfile.Country).To(BeEquivalentTo(profile.Country))
				Expect(returnedProfile.DisplayName).To(BeEquivalentTo(profile.DisplayName))
			})
		})
	})

	Context("UserService.SetRank()", func() {
		When("a profile with valid guid is given", func() {
			It("sets the rank", func() {
				var profile = api.UserProfile{
					DisplayName: "hi",
					Country:     "XX",
				}

				guid, err := userService.Create(&profile)
				Expect(err).To(BeNil())

				profile.UserId = guid
				err = userService.SetRank(&profile, "GLOBAL")
				Expect(err).To(BeNil())
				Expect(profile.Rank).NotTo(BeEquivalentTo(0))
			})
		})

		When("a profile with invalid guid is given", func() {
			It("returns error", func() {
				var profile = api.UserProfile{
					UserId:      uuid.New().String(),
					DisplayName: "hi",
					Country:     "XX",
				}

				err := userService.SetRank(&profile, "GLOBAL")
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("UserService.GetAllByID()", func() {
		When("a valid guid list is given", func() {
			It("returns profiles", func() {
				const nUsers = 5
				var uuids []string
				for i := 0; i < nUsers; i++ {
					guid := uuid.New().String()
					uuids = append(uuids, guid)
					profile := &api.UserProfile{
						UserId:      guid,
						DisplayName: "hi",
						Country:     "XX",
					}

					_, err := userService.Create(profile)
					Expect(err).To(BeNil())
				}

				profiles, err := userService.GetAllByID(uuids...)
				Expect(err).To(BeNil())

				for i := 0; i < nUsers; i++ {
					Expect(profiles[i].UserId).To(BeEquivalentTo(uuids[i]))
				}
			})
		})
	})
})

var _ = Describe("Leaderboard Service", func() {
	var (
		leaderboardService api.LeaderboardService
	)

	BeforeEach(func() {
		leaderboardService = getRedisMockedLeaderboardService(mRedis.Addr(), 20)
	})

	AfterEach(func() {
		mRedis.FlushAll()
	})

	Context("LeaderboardService.GetPage()", func() {

		When("page=1, pageSize=10", func() {

			It("returns 10 items", func() {
				page, err := leaderboardService.GetPage("GLOBAL", 1, 10)
				Expect(err).To(BeNil())
				Expect(len(page)).To(BeEquivalentTo(10))
			})

			It("returns items in ascending rank order", func() {
				page, err := leaderboardService.GetPage("GLOBAL", 1, 10)
				Expect(err).To(BeNil())

				for i := 0; i < len(page)-1; i++ {
					Expect(page[i].Rank < page[i+1].Rank).To(BeTrue())
				}
			})

			It("returns items in ascending score order", func() {
				page, err := leaderboardService.GetPage("GLOBAL", 1, 10)
				Expect(err).To(BeNil())

				for i := 0; i < len(page)-1; i++ {
					Expect(page[i].Points > page[i+1].Points).To(BeTrue())
				}
			})

			It("returns first page where it returns first row as top player", func() {
				page, err := leaderboardService.GetPage("GLOBAL", 1, 10)
				Expect(err).To(BeNil())
				Expect(page[0].Rank).To(BeEquivalentTo(1))
			})

		})
	})
})

func buildDependencies(redisAddr string) (*services.UserService, *services.RedisService) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		PoolSize: 64,
	})

	redisService := services.NewRedisService(redisClient, KeyPrefix)
	userService := services.NewUserService(redisService, KeyPrefix)

	return userService, redisService
}

func getRedisMockedLeaderboardService(redisAddr string, nPrefillUsers int) *services.LeaderboardService {
	userService, redisService := buildDependencies(redisAddr)
	generateUsers(userService, redisService, nPrefillUsers)

	return services.NewLeaderboardService(userService, redisService, KeyPrefix)
}

func generateUsers(userService *services.UserService, redisService api.RedisService, nUsers int) {
	task := tasks.NewGenerateUsersSingletonTask(userService, redisService)
	_, _ = task.Start(uint64(nUsers), 1)

	status, _ := task.Status()
	for status.RemainingUsers > 0 {
		status, _ = task.Status()
	}
}
