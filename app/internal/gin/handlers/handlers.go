package handlers

import (
	"fmt"
	chatgptgo "github.com/AidenHadisi/chat-gpt-go"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"itmo-profile/config"
	"itmo-profile/internal/database"
	"itmo-profile/model/events"
	"itmo-profile/model/user"
	"itmo-profile/model/vote"
	"itmo-profile/pkg/tools"
	"log"
	"net/http"
	"strconv"
)

func Login(context *gin.Context) {
	idS := context.PostForm("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbUser, err := user.User{}.Get(id, db)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lowestName, lowestValue := dbUser.GetLowestPoint()

	userEventsSlice, err := dbUser.GetUserEventsByTag(lowestName, db)
	if err != nil {
		log.Println(err)
	}

	eventIdSlice := []int{0}
	for _, value := range userEventsSlice {
		eventIdSlice = append(eventIdSlice, value.Id)
	}
	eventsSlice, err := events.Event{}.GetEventsByTagNoUser(lowestName, db, eventIdSlice)
	if err != nil {
		log.Println(err)
	}
	var message string
	if len(eventsSlice) <= 0 {
		api := chatgptgo.NewApi(config.GetEnv("GPT_TOKEN", ""))
		request := &chatgptgo.Request{
			MaxTokens:   100,
			Temperature: 0.5,
			N:           1,
			Model:       "gpt-3.5-turbo-16k-0613",
			Messages: []*chatgptgo.Message{
				{
					Role:    "user",
					Content: "Сгенерируйте короткую мотивирующую(подбадривающую) цитату(пример: ИТМО рождает гениев и ты один из них, ты молодец, ты совсем справишься и тд) на русском языке, учитывая это слово, оно должно быть в тексте и должно быть подставлено по смыслу:" + lowestName + ". Цитата должна поддержать человека и быть персонализирована(он связан с ниу итмо). Также цитата должна заканчиваться по смыслу и общее количество слов в цитате обязательно должно быть не больше 15 слов",
				},
			},
		}

		response, err := api.Chat(request)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message = response.Choices[0].Message.Content

		_, err = db.Exec("INSERT INTO message_log (user_id, be_itmo_value, be_itmo, gpt_text) values ($1,$2,$3,$4)", dbUser.Id, lowestValue, lowestName, message)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	quoted := pq.QuoteIdentifier(lowestName)
	if (lowestValue + 5) >= 100 {
		_, err := db.Exec(fmt.Sprintf("UPDATE users SET %s=100 WHERE id=$1", quoted), dbUser.Id)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.Exec(fmt.Sprintf("UPDATE users SET %s=%s+5 WHERE id=$1", quoted, quoted), dbUser.Id)
		if err != nil {
			log.Println(err)
		}
	}

	_, CatImageId := dbUser.GetCatImageSrc(lowestValue, db)

	context.JSON(http.StatusOK, gin.H{"user": dbUser, "events": eventsSlice, "message": message, "catImage": CatImageId})
	return
}

func BeItmo(context *gin.Context) {
	idS := context.PostForm("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbUser, err := user.User{}.Get(id, db)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var Response map[string]map[string]interface{}
	Response = map[string]map[string]interface{}{}
	tagMap := dbUser.GetBeMap()
	for index := range tagMap {
		Response[index] = map[string]interface{}{}
		voteDb := vote.GetVoteByTag(index, db)
		voteDb.Status = voteDb.CheckUserVote(dbUser.Id, db)
		Response[index]["vote"] = voteDb

		eventSlice, err := events.Event{}.GetEventsByTag(index, db)
		if err != nil {
			log.Println(err)
		}

		Response[index]["events"] = eventSlice
	}

	context.JSON(http.StatusOK, gin.H{"data": Response})
}

func Vote(context *gin.Context) {
	voteId := context.PostForm("voteId")
	userId := context.PostForm("userId")
	beItm := context.PostForm("beItmo")
	vote1Value1 := tools.ConvertStrToInt(context.PostForm("vote1Value1"))
	vote1Value2 := tools.ConvertStrToInt(context.PostForm("vote1Value2"))
	vote2Value1 := tools.ConvertStrToInt(context.PostForm("vote2Value1"))
	vote2Value2 := tools.ConvertStrToInt(context.PostForm("vote2Value2"))
	vote2Value3 := tools.ConvertStrToInt(context.PostForm("vote2Value3"))

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if vote1Value1 == 1 {
		err := vote.UpdateBeItmo(beItm, tools.ConvertStrToInt(userId), 5, db)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if vote1Value2 == 1 {
		err := vote.UpdateBeItmo(beItm, tools.ConvertStrToInt(userId), -5, db)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if vote2Value1 == 1 {
		err := vote.UpdateBeItmo(beItm, tools.ConvertStrToInt(userId), 5, db)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if vote2Value2 == 1 {
		err := vote.UpdateBeItmo(beItm, tools.ConvertStrToInt(userId), 5, db)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if vote2Value3 == 1 {
		err := vote.UpdateBeItmo(beItm, tools.ConvertStrToInt(userId), -5, db)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	_, err = db.Exec("INSERT INTO vote_log (user_id, vote_id, vote_date, vote_1_value_1, vote_1_value_2, vote_2_value_1, vote_2_value_2, vote_2_value_3) VALUES ($1,$2,current_date,$3,$4,$5,$6,$7)",
		userId, voteId, vote1Value1, vote1Value2, vote2Value1, vote2Value2, vote2Value3)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := user.User{}.Get(tools.ConvertStrToInt(userId), db)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": dbUser})
}
