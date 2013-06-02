package main
//package controllers

import (
	"fmt"
	)

type Q struct {
     question_id int
     answer string
}

type R struct {
     user_id string
     questions []Q
}

var best_girl_id = "girl"
var best_girl_phone = "123412341234"


func FindAnswer(test_suit_id int) []R{

     q1 := Q{101, "answer 1"}
     q2 := Q{102, "answer 2"}
     r := []R{{"hadwinzhy",[]Q{q1,q2}}, {"girl", []Q{q2,q1}}}
     return r
}

func GetBestFitUser(current_user_id string, test_suit_id int) bool{
     all := FindAnswer(test_suit_id)  
     // find girl's answer
     
     var girl = R{};
     for _, g:= range all{
     	 if g.user_id == best_girl_id{
	    girl = g
	    break
	 }
     }
     fmt.Println( girl)
     fmt.Println("----------------")
     for _, b:= range all{
     	 if b.user_id == best_girl_id{
	    continue	    
	 }
	 for _, boyQuestion:= range b.questions{
	     girl_answer
	     if boyQuestion.answer  != girl.questions[boyQuestion.question_id]{
	     	     fmt.Println(boyQuestion.answer)
             }
	 }

     }

     fmt.Println("----------------")
     return false
}


func GetResult(user_id string, test_suit_id int) string{
     has_found := GetBestFitUser(user_id, test_suit_id)    

     if has_found {
        return "Congratulation, u got her information:\n " + best_girl_phone
     }
     return "sorry, u can't get her name"
}


func main (){
   //new data 
   var m map[int]string

   m = make(map[int]string)
   m[101] = "first answer"
   m[102] = "second answer"
   
   fmt.Println(GetResult("asdf", 1))
}
