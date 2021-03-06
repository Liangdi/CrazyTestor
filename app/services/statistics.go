/**
 * Created with IntelliJ IDEA.
 * User: liangdi
 * Date: 3/30/13
 * Time: 12:34 PM
 * To change this template use File | Settings | File Templates.
 */
package services

import (
	"labix.org/v2/mgo"
	"strconv"
	"CrazyTestor/app/models"
	"labix.org/v2/mgo/bson"
)

var (
	_statisticsServiceInstance *StatisticsService
)

const (
	STATISTICS_COLLECTION = "statistics"
)

type StatisticsService struct {
	MongodbService
}

func (us *StatisticsService) Get(id int64) string {
	return "test:" + strconv.Itoa(int(id))
}

func (us *StatisticsService) Count() int {
	ret := 0
	ret, _ = us.c.Count()
	return ret
}

func (us *StatisticsService) Delete(id int64){
	us.c.Remove(bson.M{"id":id})
}

func (us *StatisticsService) DeleteByQuestion(questionId int64) {
	us.c.Remove(bson.M{"questionid":questionId})
}

func (us *StatisticsService) GetOne() *models.Statistics {
	result := &models.Statistics{}
	us.c.Find(nil).One(result)
	return result
}

func (us *StatisticsService) Add(s *models.Statistics) {
	us.mutex.Lock()
	s.Id = GetIdsService().GetNext(us.coll)
	us.c.Insert(s)
	us.mutex.Unlock()
}

func (us *StatisticsService) Find(questionId int64) []models.Statistics{
	result :=[]models.Statistics{}
	//us.c.Find(bson.M{"questionid": questionId}).All(&result)
	return result
}

func (us *StatisticsService) FindAnswer(userId string, suitId int64) map[int64]string{
	ret := make(map[int64]string)
	us.c.Find(bson.M{"userId":userId, "suitId":suitId}).Select(bson.M{"questionId":1, "content":2}).All(&ret)	
				
	return ret
}

func InitStatisticsService(session *mgo.Session, db *mgo.Database) {
	if _statisticsServiceInstance == nil {
		_statisticsServiceInstance = &StatisticsService{}
		_statisticsServiceInstance.Init(session, db, STATISTICS_COLLECTION)
	}
}

func GetStatisticsService() *StatisticsService {
	return _statisticsServiceInstance
}



